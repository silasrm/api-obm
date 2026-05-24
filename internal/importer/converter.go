package importer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	reCreateTable    = regexp.MustCompile(`(?i)^CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS\s+`+"`"+`(\w+)`+"`"+`\s*\(`)
	reInsertIgnore   = regexp.MustCompile(`(?i)^INSERT\s+IGNORE\s+INTO\s+`+"`"+`(\w+)`+"`"+`\s*\(([^)]+)\)\s*VALUES\s*`)
	rePrimaryKey     = regexp.MustCompile(`(?i)^\s*PRIMARY\s+KEY\s*\(([^)]+)\)`)
	reBacktick       = regexp.MustCompile("`([^`]+)`")
	reBigIntAutoInc  = regexp.MustCompile("(?i)`(\\w+)`\\s+bigint\\s+unsigned\\s+NOT\\s+NULL\\s+AUTO_INCREMENT(?:\\s+COMMENT\\s+'[^']*')?")
	reBigIntUnsigned = regexp.MustCompile("(?i)bigint\\s+unsigned")
	reBigIntDefNull  = regexp.MustCompile("(?i)bigint\\s+DEFAULT\\s+NULL")
	reBigIntNotNull  = regexp.MustCompile("(?i)bigint\\s+NOT\\s+NULL")
	reTinyint1       = regexp.MustCompile("(?i)tinyint\\(1\\)")
	reTinyint        = regexp.MustCompile("(?i)tinyint")
	reSmallint       = regexp.MustCompile("(?i)smallint")
	reLongtextCS     = regexp.MustCompile("(?i)longtext\\s+CHARACTER\\s+SET\\s+\\S+\\s+COLLATE\\s+\\S+")
	reLongtextPlain  = regexp.MustCompile("(?i)longtext")
	reDouble         = regexp.MustCompile("(?i)\\bdouble\\b")
	reDecimal        = regexp.MustCompile("(?i)decimal\\((\\d+),(\\d+)\\)")
	reVarcharCS      = regexp.MustCompile("(?i)varchar\\((\\d+)\\)\\s+CHARACTER\\s+SET\\s+\\S+\\s+COLLATE\\s+\\S+")
	reVarcharCollate = regexp.MustCompile("(?i)varchar\\((\\d+)\\)\\s+COLLATE\\s+\\S+")
	reVarcharPlain   = regexp.MustCompile("(?i)varchar\\((\\d+)\\)")
	reTextCS         = regexp.MustCompile("(?i)\\btext\\s+CHARACTER\\s+SET\\s+\\S+\\s+COLLATE\\s+\\S+")
	reTimestampNull  = regexp.MustCompile("(?i)timestamp\\s+NULL\\s+DEFAULT\\s+NULL")
	reTimestamp      = regexp.MustCompile("(?i)\\btimestamp\\b")
	reIntUnsignedAuto = regexp.MustCompile("(?i)`(\\w+)`\\s+int\\s+unsigned\\s+NOT\\s+NULL\\s+AUTO_INCREMENT")
	reIntUnsigned    = regexp.MustCompile("(?i)int\\s+unsigned")
	reIntNotNull     = regexp.MustCompile("(?i)\\bint\\s+NOT\\s+NULL\\b")
	reIntPlain       = regexp.MustCompile("(?i)\\bint\\b")
	reEngine         = regexp.MustCompile("(?i)\\)\\s*ENGINE=.*")
	reUniqKey        = regexp.MustCompile("(?i)^\\s*UNIQUE\\s+KEY\\s+`?\\w+`?\\s+\\(([^)]+)\\)")
	reKey            = regexp.MustCompile("(?i)^\\s*KEY\\s+")
	reConstraint     = regexp.MustCompile("(?i)^\\s*CONSTRAINT\\s+")
	reMySQLSet       = regexp.MustCompile(`(?i)^/\*!.*\*/;?$`)
	reSetStatement   = regexp.MustCompile(`(?i)^SET\s+`)
	reBackslashQuote = regexp.MustCompile(`\\'`)
	reColumnComment  = regexp.MustCompile(`(?i)\s+COMMENT\s+'[^']*'`)
	reBigSerialCol   = regexp.MustCompile("(?i)^\\s*(\\w+)\\s+BIGSERIAL\\s*,?\\s*$")
	reSerialCol      = regexp.MustCompile("(?i)^\\s*(\\w+)\\s+SERIAL\\s*,?\\s*$")
)

var SkipTables = map[string]bool{
	"roles":      true,
	"role_users": true,
	"sessions":   true,
}

const BatchSize = 5000

type Options struct {
	FullDrop       bool
	AppendMetadata bool
}

type Converter struct {
	writer       *bufio.Writer
	state        string
	currentTable string
	autoIncCol   string
	pkCols       string
	insertHeader string
	insertBatch  []string
	colLines     []string
	extraLines   []string
}

func NewConverter(writer io.Writer, opts Options) *Converter {
	return &Converter{
		writer: bufio.NewWriter(writer),
		state:  "normal",
	}
}

func Convert(in io.Reader, out io.Writer, opts Options) error {
	c := NewConverter(out, opts)

	if opts.FullDrop {
		c.writeln("DROP SCHEMA public CASCADE;")
		c.writeln("CREATE SCHEMA public;")
	}

	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 16*1024*1024), 32*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\r")
		c.processLine(line)
	}

	if c.state == "insert" && len(c.insertBatch) > 0 {
		c.flushBatch()
	}

	if opts.AppendMetadata {
		c.writeln("")
		c.writeln("CREATE TABLE IF NOT EXISTS users (")
		c.writeln(" id SERIAL PRIMARY KEY,")
		c.writeln(" username VARCHAR(50) NOT NULL UNIQUE,")
		c.writeln(" password_hash VARCHAR(255) NOT NULL,")
		c.writeln(" active BOOLEAN NOT NULL DEFAULT TRUE,")
		c.writeln(" created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()")
		c.writeln(");")

		c.writeln("")
		c.writeln("CREATE TABLE IF NOT EXISTS obm_metadata (")
		c.writeln(" key VARCHAR(100) PRIMARY KEY,")
		c.writeln(" value TEXT NOT NULL,")
		c.writeln(" updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()")
		c.writeln(");")
	}

	c.writer.Flush()

	return scanner.Err()
}

func (c *Converter) processLine(line string) {
	switch c.state {
	case "normal":
		c.processNormal(line)
	case "create":
		c.processCreate(line)
	case "create_skip":
		c.processCreateSkip(line)
	case "insert":
		c.processInsert(line)
	case "insert_skip":
		c.processInsertSkip(line)
	}
}

func (c *Converter) processNormal(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}
	if strings.HasPrefix(line, "--") {
		return
	}
	if reMySQLSet.MatchString(line) {
		return
	}
	if reSetStatement.MatchString(line) {
		return
	}

	if m := reCreateTable.FindStringSubmatch(line); m != nil {
		c.currentTable = m[1]
		if SkipTables[strings.ToLower(c.currentTable)] {
			c.state = "create_skip"
			return
		}
		c.autoIncCol = ""
		c.pkCols = ""
		c.colLines = c.colLines[:0]
		c.extraLines = c.extraLines[:0]
		c.state = "create"
		return
	}

	if m := reInsertIgnore.FindStringSubmatch(line); m != nil {
		tableName := m[1]
		if SkipTables[strings.ToLower(tableName)] {
			c.state = "insert_skip"
			return
		}
		c.currentTable = tableName
		cols := reBacktick.ReplaceAllString(m[2], "$1")
		c.insertHeader = fmt.Sprintf("INSERT INTO %s (%s) VALUES", c.currentTable, cols)
		c.insertBatch = c.insertBatch[:0]

		rest := line[reInsertIgnore.FindStringIndex(line)[1]:]
		rest = strings.TrimSpace(rest)
		if rest != "" {
			c.state = "insert"
			c.processInsert(rest)
		} else {
			c.state = "insert"
		}
		return
	}
}

func (c *Converter) processCreate(line string) {
	if reEngine.MatchString(line) {
		c.emitTable()
		c.state = "normal"
		return
	}

	if pkMatch := rePrimaryKey.FindStringSubmatch(line); pkMatch != nil {
		pkCols := reBacktick.ReplaceAllString(pkMatch[1], "$1")
		c.pkCols = pkCols
		return
	}

	if m := reUniqKey.FindStringSubmatch(line); m != nil {
		cols := reBacktick.ReplaceAllString(m[1], "$1")
		c.extraLines = append(c.extraLines, " UNIQUE ("+cols+"),")
		return
	}

	if reKey.MatchString(line) {
		return
	}

	if reConstraint.MatchString(line) {
		return
	}

	if strings.HasSuffix(strings.TrimSpace(line), ");") {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimSuffix(trimmed, ");")
		trimmed = strings.TrimRight(trimmed, ",")
		if strings.TrimSpace(trimmed) != "" {
			converted := convertColumnDef(trimmed)
			c.colLines = append(c.colLines, " "+converted)
		}
		c.emitTable()
		c.state = "normal"
		return
	}

	converted := convertColumnDef(line)
	if c.autoIncCol == "" {
		if m := reBigIntAutoInc.FindStringSubmatch(line); m != nil {
			c.autoIncCol = m[1]
		}
		if m := reIntUnsignedAuto.FindStringSubmatch(line); m != nil {
			c.autoIncCol = m[1]
		}
	}
	c.colLines = append(c.colLines, " "+converted)
}

func (c *Converter) processCreateSkip(line string) {
	if reEngine.MatchString(line) {
		c.state = "normal"
		return
	}
	trimmed := strings.TrimSpace(line)
	if strings.HasSuffix(trimmed, ");") {
		c.state = "normal"
		return
	}
}

func (c *Converter) processInsert(line string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		if len(c.insertBatch) > 0 {
			c.flushBatch()
		}
		c.state = "normal"
		return
	}

	isEnd := strings.HasSuffix(trimmed, ";")
	if isEnd {
		trimmed = strings.TrimSuffix(trimmed, ";")
		trimmed = strings.TrimRight(trimmed, ",")
		trimmed = strings.TrimSpace(trimmed)
	}

	if trimmed == "" {
		if len(c.insertBatch) > 0 {
			c.flushBatch()
		}
		c.state = "normal"
		return
	}

	converted := reBackslashQuote.ReplaceAllString(trimmed, "''")
	converted = reBacktick.ReplaceAllString(converted, "$1")

	tuples := c.extractTuples(converted)
	for _, t := range tuples {
		row := "(" + t + ")"
		c.insertBatch = append(c.insertBatch, row)
		if len(c.insertBatch) >= BatchSize {
			c.flushBatch()
		}
	}

	if isEnd {
		if len(c.insertBatch) > 0 {
			c.flushBatch()
		}
		c.state = "normal"
	}
}

func (c *Converter) processInsertSkip(line string) {
	trimmed := strings.TrimSpace(line)
	if strings.HasSuffix(trimmed, ";") {
		c.state = "normal"
	}
}

func (c *Converter) emitTable() {
	c.writeln("CREATE TABLE IF NOT EXISTS " + c.currentTable + " (")
	isPKAutoInc := false
	if c.autoIncCol != "" && c.pkCols != "" {
		normalizedPK := strings.ReplaceAll(strings.ReplaceAll(c.pkCols, " ", ""), "`", "")
		normalizedAutoInc := strings.ReplaceAll(c.autoIncCol, "`", "")
		if strings.EqualFold(normalizedPK, normalizedAutoInc) {
			isPKAutoInc = true
		}
	}

	allLines := make([]string, 0, len(c.colLines)+len(c.extraLines)+1)
	for _, colLine := range c.colLines {
		if isPKAutoInc {
			if m := reBigSerialCol.FindStringSubmatch(colLine); m != nil {
				if strings.EqualFold(m[1], c.autoIncCol) {
					allLines = append(allLines, " "+c.autoIncCol+" BIGSERIAL PRIMARY KEY")
					continue
				}
			}
			if m := reSerialCol.FindStringSubmatch(colLine); m != nil {
				if strings.EqualFold(m[1], c.autoIncCol) {
					allLines = append(allLines, " "+c.autoIncCol+" SERIAL PRIMARY KEY")
					continue
				}
			}
		}
		allLines = append(allLines, colLine)
	}

	for _, extra := range c.extraLines {
		allLines = append(allLines, extra)
	}

	if !isPKAutoInc && c.pkCols != "" {
		allLines = append(allLines, " PRIMARY KEY ("+c.pkCols+")")
	}

	for i, l := range allLines {
		l = strings.TrimRight(strings.TrimSpace(l), ",")
		if i < len(allLines)-1 {
			l += ","
		}
		c.writeln(l)
	}

	c.writeln(");")
}

func (c *Converter) flushBatch() {
	if len(c.insertBatch) == 0 {
		return
	}
	c.writeln(c.insertHeader)
	for i, row := range c.insertBatch {
		if i < len(c.insertBatch)-1 {
			fmt.Fprintf(c.writer, "%s,\n", row)
		} else {
			fmt.Fprintf(c.writer, "%s\n", row)
		}
	}
	c.writeln("ON CONFLICT DO NOTHING;")
	c.insertBatch = c.insertBatch[:0]
}

func (c *Converter) extractTuples(data string) []string {
	var tuples []string
	var buf strings.Builder
	depth := 0
	inQuote := false
	escape := false

	for i := 0; i < len(data); i++ {
		ch := data[i]

		if escape {
			buf.WriteByte(ch)
			escape = false
			continue
		}

		if ch == '\\' && inQuote {
			buf.WriteByte(ch)
			escape = true
			continue
		}

		if ch == '\'' {
			inQuote = !inQuote
			buf.WriteByte(ch)
			continue
		}

		if inQuote {
			buf.WriteByte(ch)
			continue
		}

		if ch == '(' {
			if depth == 0 {
				buf.Reset()
			} else {
				buf.WriteByte(ch)
			}
			depth++
			continue
		}

		if ch == ')' {
			depth--
			if depth == 0 {
				tuples = append(tuples, buf.String())
				buf.Reset()
			} else {
				buf.WriteByte(ch)
			}
			continue
		}

		buf.WriteByte(ch)
	}

	return tuples
}

func convertColumnDef(line string) string {
	converted := line

	converted = reBigIntAutoInc.ReplaceAllString(converted, "$1 BIGSERIAL")
	converted = reBigIntUnsigned.ReplaceAllString(converted, "BIGINT")
	converted = reBigIntDefNull.ReplaceAllString(converted, "BIGINT")
	converted = reBigIntNotNull.ReplaceAllString(converted, "BIGINT NOT NULL")

	converted = reIntUnsignedAuto.ReplaceAllString(converted, "$1 SERIAL")
	converted = reIntUnsigned.ReplaceAllString(converted, "INTEGER")
	converted = reIntNotNull.ReplaceAllString(converted, "INTEGER NOT NULL")
	converted = reIntPlain.ReplaceAllString(converted, "INTEGER")

	converted = reVarcharCS.ReplaceAllString(converted, "VARCHAR($1)")
	converted = reVarcharCollate.ReplaceAllString(converted, "VARCHAR($1)")
	converted = reVarcharPlain.ReplaceAllString(converted, "VARCHAR($1)")

	converted = reLongtextCS.ReplaceAllString(converted, "TEXT")
	converted = reLongtextPlain.ReplaceAllString(converted, "TEXT")

	converted = reTextCS.ReplaceAllString(converted, "TEXT")

	converted = reTinyint1.ReplaceAllString(converted, "BOOLEAN")
	converted = reTinyint.ReplaceAllString(converted, "SMALLINT")

	converted = reSmallint.ReplaceAllString(converted, "SMALLINT")

	converted = reDecimal.ReplaceAllString(converted, "NUMERIC($1,$2)")
	converted = reDouble.ReplaceAllString(converted, "DOUBLE PRECISION")

	converted = reTimestampNull.ReplaceAllString(converted, "TIMESTAMPTZ")
	converted = reTimestamp.ReplaceAllString(converted, "TIMESTAMPTZ")

	converted = reBacktick.ReplaceAllString(converted, "$1")

	converted = reColumnComment.ReplaceAllString(converted, "")

	return converted
}

func (c *Converter) writeln(line string) {
	fmt.Fprintln(c.writer, line)
}
