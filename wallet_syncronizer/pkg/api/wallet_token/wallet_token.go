package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-syncronizer/pkg/api/wallet_token/get"
)

type Env struct {
	DB *gorm.DB
}

func NewEnv(db *gorm.DB) *Env {
	return &Env{DB: db}
}

func (e *Env) Get(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}

func (e *Env) List(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}
