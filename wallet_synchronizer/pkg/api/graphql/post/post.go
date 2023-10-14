package post

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	graphql_util "wallet-synchronizer/pkg/graphql"
	json_util "wallet-synchronizer/pkg/util/json"
)

type Api struct {
	schema graphql.Schema
	query  string
	fields logrus.Fields
}

func NewApi(
	uuid string,
	url string,
	schema graphql.Schema,
	query string,
) *Api {
	return &Api{
		fields: logrus.Fields{"uuid": uuid, "url": url, "query": query},
		schema: schema,
		query:  query,
	}
}

func (a *Api) Post() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")

	result := graphql.Do(graphql.Params{
		Schema:        a.schema,
		RequestString: a.query,
	})

	if len(result.Errors) > 0 {
		gqlError := result.Errors[0]
		statusCode := graphql_util.MapGraphQLErrorToHTTPStatus(&gqlError)
		errB, _ := json.Marshal(result.Errors)
		logrus.WithFields(a.fields).WithError(errors.New(string(errB))).Errorf("terminated with failure")
		return statusCode, json_util.NewErrorResponse(statusCode, result.Errors)
	}

	logrus.WithFields(a.fields).Info("terminated with success")
	return fiber.StatusOK, result
}
