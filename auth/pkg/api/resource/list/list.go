package list

type ListApi struct {
}

func NewListApi() *ListApi {
	return &ListApi{}
}

func (List *ListApi) List() (httpStatus int, response interface{}) {
}
