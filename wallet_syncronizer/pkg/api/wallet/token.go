package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-syncronizer/pkg/api/token/get"
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

func (e *Env) GetList(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}

func (e *Env) Create(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}

func (e *Env) Update(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}

func (e *Env) Delete(ctx *fiber.Ctx) error {
	return get.Get(ctx, e.DB)
}
