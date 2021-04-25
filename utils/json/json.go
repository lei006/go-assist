package json

import "encoding/json"

func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
}
