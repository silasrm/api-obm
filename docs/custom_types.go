package docs

type NullBool struct {
	Bool  bool `json:"bool"`
	Valid bool `json:"valid"`
}

type NullInt64 struct {
	Int64 int64 `json:"int64"`
	Valid bool  `json:"valid"`
}

type NullFloat64 struct {
	Float64 float64 `json:"float64"`
	Valid   bool    `json:"valid"`
}

type NullString struct {
	String string `json:"string"`
	Valid  bool   `json:"valid"`
}

type NullTime struct {
	Time  string `json:"time"`
	Valid bool   `json:"valid"`
}
