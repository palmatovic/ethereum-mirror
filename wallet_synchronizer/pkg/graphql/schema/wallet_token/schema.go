package wallet_token

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
	wallet_token_get_service "wallet-synchronizer/pkg/service/wallet_token/get"
	wallet_token_list_service "wallet-synchronizer/pkg/service/wallet_token/list"
	token_url "wallet-synchronizer/pkg/url/token"
	wallet_url "wallet-synchronizer/pkg/url/wallet"
)

func Schema(database *gorm.DB) graphql.Schema {

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
