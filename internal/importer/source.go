package importer

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ReadCloserWrapper struct {
	io.Reader
	Cleanup func()
}

func (r *ReadCloserWrapper) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}

func (r *ReadCloserWrapper) Close() error {
	if r.Cleanup != nil {
		r.Cleanup()
		r.Cleanup = nil
	}
	return nil
}

func Resolve(source string) (io.ReadCloser, error) {
	if strings.HasPrefix(source, "mysql://") {
		return resolveMysql(source)
	}
	lower := strings.ToLower(source)
	if strings.HasSuffix(lower, ".zip") {
		return resolveZip(source)
	}
	if strings.HasSuffix(lower, ".sql.gz") {
		return resolveGzip(source)
	}
	return resolvePlain(source)
}

func resolveMysql(source string) (io.ReadCloser, error) {
	_, err := exec.LookPath("mysqldump")
	if err != nil {
		return nil, fmt.Errorf("mysqldump not found on PATH — required for mysql:// source")
	}

	u, err := url.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("invalid mysql URL: %w", err)
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "3306"
	}

	user := u.User.Username()
	password, _ := u.User.Password()
	database := strings.TrimPrefix(u.Path, "/")

	args := []string{
		"-h", host,
		"-P", port,
		"-u", user,
	}
	if password != "" {
		args = append(args, "-p"+password)
	}
	args = append(args, database, "--skip-comments", "--no-create-db", "--single-transaction", "--quick")

	cmd := exec.Command("mysqldump", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start mysqldump: %w", err)
	}

	return &ReadCloserWrapper{
		Reader:  stdout,
		Cleanup: func() { cmd.Wait() },
	}, nil
}

func resolveZip(source string) (io.ReadCloser, error) {
	zr, err := zip.OpenReader(source)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP: %w", err)
	}

	var target *zip.File
	for _, f := range zr.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".sql") {
			target = f
			break
		}
	}
	if target == nil {
		zr.Close()
		return nil, fmt.Errorf("no .sql file found inside ZIP")
	}

	rc, err := target.Open()
	if err != nil {
		zr.Close()
		return nil, fmt.Errorf("failed to open entry in ZIP: %w", err)
	}

	tmp, err := os.CreateTemp("", "obm-import-*.sql")
	if err != nil {
		rc.Close()
		zr.Close()
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(tmp, rc)
	rc.Close()
	zr.Close()
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("failed to extract from ZIP: %w", err)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("failed to seek temp file: %w", err)
	}

	name := tmp.Name()
	return &ReadCloserWrapper{
		Reader: tmp,
		Cleanup: func() {
			tmp.Close()
			os.Remove(name)
		},
	}, nil
}

func resolveGzip(source string) (io.ReadCloser, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	gz, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	return &ReadCloserWrapper{
		Reader: gz,
		Cleanup: func() {
			gz.Close()
			f.Close()
		},
	}, nil
}

func resolvePlain(source string) (io.ReadCloser, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &ReadCloserWrapper{
		Reader:  f,
		Cleanup: func() { f.Close() },
	}, nil
}

func ExtractVersion(source string) string {
	base := filepath.Base(source)
	parts := strings.Split(base, "-")
	for _, p := range parts {
		p = strings.TrimSuffix(p, filepath.Ext(p))
		if len(p) == 8 {
			allDigit := true
			for _, c := range p {
				if c < '0' || c > '9' {
					allDigit = false
					break
				}
			}
			if allDigit {
				return p
			}
		}
	}
	return ""
}
