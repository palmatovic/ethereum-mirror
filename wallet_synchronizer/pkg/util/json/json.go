package json

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
