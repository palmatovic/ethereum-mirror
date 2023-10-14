package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	token_db "wallet-synchronizer/pkg/database/token"
	token_get_service "wallet-synchronizer/pkg/service/token/get"
	token_list_service "wallet-synchronizer/pkg/service/token/list"
	graphql_util "wallet-synchronizer/pkg/util/graphql"
	json_util "wallet-synchronizer/pkg/util/json"
	token_url "wallet-synchronizer/pkg/util/url/token"
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
		Schema:        getSchema(a.database),
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

func getSchema(database *gorm.DB) graphql.Schema {

	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "TokenQuery",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: token_db.TokenGraphQL,
				Args: graphql.FieldConfigArgument{
					string(token_url.Id): &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tokenId, ok := p.Args[string(token_url.Id)].(string)
					if !ok || len(tokenId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(token_url.Id))
					}

					var token *token_db.Token
					var err error

					_, token, err = token_get_service.NewService(database, tokenId).Get()
					if err != nil {
						return nil, err
					}

					return token, nil
				},
			},
			"tokens": &graphql.Field{
				Type: graphql.NewList(token_db.TokenGraphQL),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var tokens *[]token_db.Token
					var err error
					_, tokens, err = token_list_service.NewService(database).List()
					if err != nil {
						return nil, err
					}
					return tokens, nil
				},
			},
		},
	})

	var schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		logrus.WithError(err).Panic("failed to create token schema")
	}

	return schema
}
