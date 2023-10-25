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
	Kid          string       `json:"kid"`
	Iat          int64        `json:"iat"`
	Exp          int64        `json:"exp"`
	UserId       int          `json:"user_id"`
	Username     string       `json:"username"`
	Aud          string       `json:"aud"`
	Iss          string       `json:"iss"`
	Role         int          `json:"role"`
	Email        string       `json:"email"`
	Organization Organization `json:"organization"`
	Name         string       `json:"name"`
	Surname      string       `json:"surname"`
	Resources    []Resource   `json:"resources"`
}

type Resource struct {
	Name  string   `json:"name"`
	Perms []string `json:"perms"`
}

type Organization struct {
	Id   string `json:"id"`
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
