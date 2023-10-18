package token

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	token_db "wallet-synchronizer/pkg/database/token"
	token_get_service "wallet-synchronizer/pkg/service/token/get"
	token_list_service "wallet-synchronizer/pkg/service/token/list"
	token_url "wallet-synchronizer/pkg/url/token"
)

func Schema(database *gorm.DB) graphql.Schema {

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
