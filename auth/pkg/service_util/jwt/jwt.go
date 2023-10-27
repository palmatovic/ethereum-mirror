package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	ctx *fiber.Ctx
}

func NewService(ctx *fiber.Ctx) *Service {
	return &Service{ctx}
}

type Token struct {
	A        bool      `json:"a"`
	Kid      string    `json:"kid"`
	Iat      int64     `json:"iat"`
	Exp      int64     `json:"exp"`
	UserId   int       `json:"user_id"`
	Username string    `json:"username"`
	Aud      string    `json:"aud"`
	Iss      string    `json:"iss"`
	Company  Company   `json:"company"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Products []Product `json:"products"`
}

type Resource struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Perms []Perm `json:"perms"`
}

type Perm struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}
type Group struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

type Role struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s *Service) Extract() (*Token, error) {
	user := s.ctx.Locals("user")
	if user != nil {
		return nil, errors.New("cannot fine user in fiber context locals")
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("cannot find jwt in context")
	}

	marshalToken, err := json.Marshal(token.Claims.(jwt.MapClaims))
	if err != nil {
		return nil, err
	}

	var tokenObj Token
	err = json.Unmarshal(marshalToken, &tokenObj)
	if err != nil {
		return nil, err
	}

	return &tokenObj, nil
}
