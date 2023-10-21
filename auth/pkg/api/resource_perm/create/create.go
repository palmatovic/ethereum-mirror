package create

type CreateApi struct {
}

func NewCreateApi() *CreateApi {
	return &CreateApi{}
}

func (create *CreateApi) Create() (httpStatus int, response interface{}) {
}
