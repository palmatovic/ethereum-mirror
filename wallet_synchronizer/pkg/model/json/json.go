package json

import "encoding/json"

type Response struct {
	Error *Error      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Error struct {
	Code   int         `json:"code"`
	Detail interface{} `json:"detail"`
}

func NewErrorResponse(httpStatus int, detail interface{}) Response {
	return Response{
		Error: &Error{
			Code:   httpStatus,
			Detail: detail,
		},
	}
}

func Stringify(data []byte) string {
	var outputMap map[string]interface{}
	var outputBytes []byte
	_ = json.Unmarshal(data, &outputMap)
	outputBytes, _ = json.Marshal(outputMap)
	return string(outputBytes)
}
