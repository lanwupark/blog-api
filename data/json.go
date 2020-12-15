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

// FromJSONString 反序列化json数据
func FromJSONString(data string, i interface{}) error {
	return json.Unmarshal([]byte(data), i)
}

// ToJSONString 封装json数据 返回
func ToJSONString(i interface{}) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
