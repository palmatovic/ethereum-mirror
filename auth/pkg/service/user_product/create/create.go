package create

import (
	user_product_db "auth/pkg/database/user_product"
	user_product_model "auth/pkg/model/api/user_product/create"
	"auth/pkg/service_util/aes"
	"crypto/sha256"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m1/go-generate-password/generator"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"time"
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

func (s *Service) Create() (status int, userProduct *user_product_db.UserProduct, err error) {

	// generate two fa key
	masterPwdKey, err := aes.NewService().NewAES256Key()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	masterTwoFAKey, err := aes.NewService().NewAES256Key()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "auth",
		AccountName: "auth",
		Algorithm:   otp.AlgorithmSHA256,
	})
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	config := generator.Config{
		Length:                     8,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   false,
		ExcludeAmbiguousCharacters: true,
	}
	g, err := generator.New(&config)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	pwd, err := g.Generate()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	sha256Pwd := sha256.Sum256([]byte(*pwd))

	sha256MasterPwdKey := sha256.Sum256(*masterPwdKey)

	sha256MasterTwoFAKey := sha256.Sum256(*masterTwoFAKey)

	sha256FAKey := sha256.Sum256([]byte(key.Secret()))

	userProduct = &user_product_db.UserProduct{
		UserId:               s.userProduct.UserId,
		ProductId:            s.userProduct.ProductId,
		Enabled:              false,
		ChangePassword:       true,
		Password:             string(sha256Pwd[:]),
		PasswordExpirationAt: time.Now().Add(time.Hour * 24 * 365),
		PasswordExpired:      false,
		MasterPasswordKey:    string(sha256MasterPwdKey[:]),
		TwoFAKey:             string(sha256FAKey[:]),
		MasterTwoFAKey:       string(sha256MasterTwoFAKey[:]),
	}

	if err = s.db.Create(userProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, userProduct, nil
}
