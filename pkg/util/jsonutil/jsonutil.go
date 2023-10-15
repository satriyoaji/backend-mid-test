package jsonutil

import (
	"encoding/json"

	"github.com/labstack/gommon/log"
)

func Stringify(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Warn("Marshal data to json error: ", err)
		return ""
	}
	return string(b)
}

func Parse(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}
