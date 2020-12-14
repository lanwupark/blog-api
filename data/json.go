package data

import (
	"encoding/json"
	"io"
)

// FromJSON deserialize
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

// ToJSON 封装json数据 将其写入到writer中
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}
