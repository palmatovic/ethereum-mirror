package get

type GetApi struct {
}

func NewGetApi() *GetApi {
	return &GetApi{}
}

func (Get *GetApi) Get() (httpStatus int, response interface{}) {
}
