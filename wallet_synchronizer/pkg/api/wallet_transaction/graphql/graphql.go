package graphql

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_transaction_graphql_service "wallet-synchronizer/pkg/service/wallet_transaction/graphql"
	graphql_util "wallet-synchronizer/pkg/util/graphql"
	json_util "wallet-synchronizer/pkg/util/json"
)

type Api struct {
	database *gorm.DB
	query    string
	fields   logrus.Fields
}

func NewApi(
	uuid string,
	url string,
	database *gorm.DB,
	query string,
) *Api {
	return &Api{
		fields:   logrus.Fields{"uuid": uuid, "url": url, "query": query},
		database: database,
		query:    query,
	}
}

func (a *Api) GraphQL() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")

	result := graphql.Do(graphql.Params{
		Schema:        wallet_transaction_graphql_service.NewService(a.database).Schema(),
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
