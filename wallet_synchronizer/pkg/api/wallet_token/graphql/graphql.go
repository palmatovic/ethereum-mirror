package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
	wallet_token_get_service "wallet-synchronizer/pkg/service/wallet_token/get"
	wallet_token_list_service "wallet-synchronizer/pkg/service/wallet_token/list"
	graphql_util "wallet-synchronizer/pkg/util/graphql"
	json_util "wallet-synchronizer/pkg/util/json"
	token_url "wallet-synchronizer/pkg/util/url/token"
	wallet_url "wallet-synchronizer/pkg/util/url/wallet"
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
		Name: "WalletTokenQuery",
		Fields: graphql.Fields{
			"wallet_token": &graphql.Field{
				Type: wallet_token_db.WalletTokenGraphQL,
				Args: graphql.FieldConfigArgument{
					string(wallet_url.Id): &graphql.ArgumentConfig{Type: graphql.String},
					string(token_url.Id):  &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					walletId, ok := p.Args[string(wallet_url.Id)].(string)
					if !ok || len(walletId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(wallet_url.Id))
					}

					tokenId, ok := p.Args[string(token_url.Id)].(string)
					if !ok || len(tokenId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(token_url.Id))
					}
					var walletToken *wallet_token_db.WalletToken
					var err error
					_, walletToken, err = wallet_token_get_service.NewService(database, walletId, tokenId).Get()
					if err != nil {
						return nil, err
					}
					return walletToken, nil
				},
			},
			"wallet_tokens": &graphql.Field{
				Type: graphql.NewList(wallet_token_db.WalletTokenGraphQL),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var walletTokens *[]wallet_token_db.WalletToken
					var err error
					_, walletTokens, err = wallet_token_list_service.NewService(database).List()
					if err != nil {
						return nil, err
					}
					return walletTokens, nil
				},
			},
		},
	})

	var schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		logrus.WithError(err).Panic("failed to create wallet_token schema")
	}

	return schema
}
