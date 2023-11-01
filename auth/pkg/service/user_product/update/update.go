package update

import (
	user_product_db "auth/pkg/database/user_product"
	user_product_model "auth/pkg/model/api/user_product/update"
	"auth/pkg/service/user_product/get"
	"crypto/sha256"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"time"
	"unicode"
)

type Service struct {
	db          *gorm.DB
	userProduct *user_product_model.UserProduct
}

func NewService(db *gorm.DB, userProduct *user_product_model.UserProduct) *Service {
	return &Service{
		db:          db,
		userProduct: userProduct,
	}
}

func (s *Service) Update() (status int, group *user_product_db.UserProduct, err error) {

	findStatus, dbUserProduct, err := get.NewService(s.db, s.userProduct.UserProductId).Get()
	if err != nil {
		if findStatus == fiber.StatusNotFound {
			return fiber.StatusNotFound, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	var count int
	if s.userProduct.RenewPassword != nil {
		count++
	}

	if s.userProduct.ForgotPassword != nil {
		count++
	}

	if s.userProduct.ForgotTwoFA != nil {
		count++
	}

	if count != 1 {
		return fiber.StatusBadRequest, nil, errors.New("only one of renew_password, forgot_password, or forgot_two_fa can be set")
	}

	switch {
	case s.userProduct.RenewPassword != nil:
		if s.userProduct.RenewPassword.NewPassword != s.userProduct.RenewPassword.RepeatNewPassword {
			return fiber.StatusBadRequest, nil, errors.New("new_password and repeat_new_password do not match")
		}
		if !validPassword(s.userProduct.RenewPassword.NewPassword) {
			return fiber.StatusBadRequest, nil, errors.New("new_password is not valid")
		}
		sha256Password := sha256.Sum256([]byte(s.userProduct.RenewPassword.NewPassword))
		dbUserProduct.Password = string(sha256Password[:])
		dbUserProduct.PasswordExpired = false
		dbUserProduct.PasswordExpirationAt = time.Now().Add(time.Hour * 24 * 365)
		dbUserProduct.ChangePassword = false
		break
	case s.userProduct.ForgotPassword != nil:
		sha256MasterPasswordKey := sha256.Sum256([]byte(s.userProduct.ForgotPassword.MasterPasswordKey))
		if string(sha256MasterPasswordKey[:]) != dbUserProduct.MasterPasswordKey {
			return fiber.StatusBadRequest, nil, errors.New("master_password_key does not match")
		}
		if s.userProduct.ForgotPassword.NewPassword != s.userProduct.ForgotPassword.RepeatNewPassword {
			return fiber.StatusBadRequest, nil, errors.New("new_password and repeat_new_password do not match")
		}
		sha256Password := sha256.Sum256([]byte(s.userProduct.RenewPassword.NewPassword))
		dbUserProduct.Password = string(sha256Password[:])
		dbUserProduct.PasswordExpired = false
		dbUserProduct.PasswordExpirationAt = time.Now().Add(time.Hour * 24 * 365)
		dbUserProduct.ChangePassword = false
		break
	case s.userProduct.ForgotTwoFA != nil:
		sha256MasterTwoFAKey := sha256.Sum256([]byte(s.userProduct.ForgotTwoFA.MasterTwoFAKey))
		if string(sha256MasterTwoFAKey[:]) != dbUserProduct.MasterTwoFAKey {
			return fiber.StatusBadRequest, nil, errors.New("master_two_fa_key does not match")
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "auth",
			AccountName: dbUserProduct.User.Username,
			Algorithm:   otp.AlgorithmSHA256,
		})
		if err != nil {
			return fiber.StatusInternalServerError, nil, err
		}
		sha256FAKey := sha256.Sum256([]byte(key.Secret()))
		dbUserProduct.TwoFAKey = string(sha256FAKey[:])
	}

	if err = s.db.Where("UserProductRoleId = ?", dbUserProduct.UserProductId).Updates(dbUserProduct).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, dbUserProduct, nil
}

func validPassword(password string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}

	return true

}
