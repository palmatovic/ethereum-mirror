package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
	wallet_create_service "wallet-synchronizer/pkg/service/wallet/create"
	wallet_get_service "wallet-synchronizer/pkg/service/wallet/get"
	wallet_list_service "wallet-synchronizer/pkg/service/wallet/list"
	wallet_url "wallet-synchronizer/pkg/url/wallet"
)

func Schema(database *gorm.DB) graphql.Schema {

	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "WalletQuery",
		Fields: graphql.Fields{
			"wallet": &graphql.Field{
				Type: wallet_db.WalletGraphQL,
				Args: graphql.FieldConfigArgument{
					string(wallet_url.Id): &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					walletId, ok := p.Args[string(wallet_url.Id)].(string)
					if !ok || len(walletId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(wallet_url.Id))
					}
					var wallet *wallet_db.Wallet
					var err error
					_, wallet, err = wallet_get_service.NewService(database, walletId).Get()
					if err != nil {
						return nil, err
					}
					return wallet, nil
				},
			},
			"wallets": &graphql.Field{
				Type: graphql.NewList(wallet_db.WalletGraphQL),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var wallets *[]wallet_db.Wallet
					var err error
					_, wallets, err = wallet_list_service.NewService(database).List()
					if err != nil {
						return nil, err
					}
					return wallets, nil
				},
			},
		},
	})

	var rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "WalletMutation",
		Fields: graphql.Fields{
			"create_wallet": &graphql.Field{
				Type: wallet_db.WalletGraphQL,
				Args: graphql.FieldConfigArgument{
					"wallet": &graphql.ArgumentConfig{Type: wallet_db.CreateWalletGraphQL},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					walletInput, ok := p.Args["wallet"].(map[string]interface{})
					if !ok || len(walletInput) == 0 {
						return nil, errors.New("wallet must be evaluated")
					}

					walletId, ok := walletInput[string(wallet_url.Id)].(string)
					if !ok || len(walletId) == 0 {
						return nil, fmt.Errorf("%s must be evaluated as a string", string(wallet_url.Id))
					}

					wByte, err := json.Marshal(walletInput)
					if err != nil {
						return nil, errors.New("cannot marshal wallet")
					}

					var wInput *wallet_db.Wallet
					err = json.Unmarshal(wByte, &wInput)
					if err != nil {
						return nil, errors.New("cannot unmarshal wallet")
					}
					var wallet *wallet_db.Wallet

					_, wallet, err = wallet_create_service.NewService(database, wInput).Create()
					if err != nil {
						return nil, err
					}

					return wallet, nil
				},
			},
		},
	})

	var schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		logrus.WithError(err).Panic("failed to create wallet schema")
	}

	return schema
}
