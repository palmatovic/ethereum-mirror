package create

import (
	perm_db "auth/pkg/database/perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	perm *perm_db.Perm
}

func NewService(db *gorm.DB, perm *perm_db.Perm) *Service {
	return &Service{
		db:   db,
		perm: perm,
	}
}

func (s *Service) Create() (status int, perm *perm_db.Perm, err error) {
	if err = s.db.Create(s.perm).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.perm, nil
}
