package delete

import (
	perm_db "auth/pkg/database/perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	permId int
}

func NewService(db *gorm.DB, permId int) *Service {
	return &Service{
		db:     db,
		permId: permId,
	}
}

func (s *Service) Delete() (status int, perm *perm_db.Perm, err error) {
	perm = new(perm_db.Perm)
	if err = s.db.Where("PermId = ?", s.permId).Delete(perm).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, perm, nil
}
