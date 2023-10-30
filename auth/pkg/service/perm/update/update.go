package update

import (
	perm_db "auth/pkg/database/perm"
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

func (s *Service) Update() (status int, perm *perm_db.Perm, err error) {
	if err = s.db.Where("PermId = ?", s.perm.PermId).Updates(s.perm).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.perm, nil
}
