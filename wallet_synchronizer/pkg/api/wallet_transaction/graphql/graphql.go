package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_transaction_db "wallet-synchronizer/pkg/database/wallet_transaction"
	wallet_transaction_get_service "wallet-synchronizer/pkg/service/wallet_transaction/get"
	wallet_transaction_list_service "wallet-synchronizer/pkg/service/wallet_transaction/list"
	graphql_util "wallet-synchronizer/pkg/util/graphql"
	json_util "wallet-synchronizer/pkg/util/json"
	wallet_transaction_url "wallet-synchronizer/pkg/util/url/wallet_transaction"
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
		Name: "WalletTransactionQuery",
		Fields: graphql.Fields{
			"wallet_transaction": &graphql.Field{
				Type: wallet_transaction_db.WalletTransactionGraphQL,
				Args: graphql.FieldConfigArgument{
					string(wallet_transaction_url.Id): &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					walletTransactionId, ok := p.Args[string(wallet_transaction_url.Id)].(string)
					if !ok || len(walletTransactionId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(wallet_transaction_url.Id))
					}

					var walletTransaction *wallet_transaction_db.WalletTransaction
					var err error
					_, walletTransaction, err = wallet_transaction_get_service.NewService(database, walletTransactionId).Get()
					if err != nil {
						return nil, err
					}
					return walletTransaction, nil
				},
			},
			"wallet_transactions": &graphql.Field{
				Type: graphql.NewList(wallet_transaction_db.WalletTransactionGraphQL),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var walletTransactions *[]wallet_transaction_db.WalletTransaction
					var err error
					_, walletTransactions, err = wallet_transaction_list_service.NewService(database).List()
					if err != nil {
						return nil, err
					}
					return walletTransactions, nil
				},
			},
		},
	})

	var schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		logrus.WithError(err).Panic("failed to create wallet_transaction schema")
	}

	return schema
}
