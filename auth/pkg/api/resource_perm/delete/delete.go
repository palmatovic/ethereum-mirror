package delete

type DeleteApi struct {
}

func NewDeleteApi() *DeleteApi {
	return &DeleteApi{}
}

func (Delete *DeleteApi) Delete() (httpStatus int, response interface{}) {
}
