package graphql

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	token_db "wallet-synchronizer/pkg/database/token"
	token_get_service "wallet-synchronizer/pkg/service/token/get"
	token_list_service "wallet-synchronizer/pkg/service/token/list"
)

type Service struct {
	database *gorm.DB
}

func NewService(
	database *gorm.DB,
) *Service {
	return &Service{
		database: database,
	}
}

var tokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "token",
	Fields: graphql.Fields{
		"token_id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"symbol": &graphql.Field{
			Type: graphql.String,
		},
		"decimals": &graphql.Field{
			Type: graphql.Int,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"logo": &graphql.Field{
			Type: graphql.String,
		},
		"go_plus_response": &graphql.Field{
			Type: graphql.NewScalar(graphql.ScalarConfig{
				Name: "Json",
				Serialize: func(value interface{}) interface{} {
					var serialized map[string]interface{}
					_ = json.Unmarshal(value.([]byte), &serialized)
					return serialized
				},
			}),
		},
	},
})

func (s *Service) Schema() graphql.Schema {

	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "TokenQuery",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: tokenType,
				Args: graphql.FieldConfigArgument{
					"token_id": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tokenId, ok := p.Args["token_id"].(string)
					var token *token_db.Token
					var err error
					if ok {
						_, token, err = token_get_service.NewService(s.database, tokenId).Get()
						if err != nil {
							return nil, err
						}
					}
					return token, nil
				},
			},
			"tokens": &graphql.Field{
				Type: graphql.NewList(tokenType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var tokens *[]token_db.Token
					var err error
					_, tokens, err = token_list_service.NewService(s.database).List()
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
