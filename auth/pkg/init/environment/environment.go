package environment

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
	"time"
)

type AppConfig struct {
	FiberPort              int64  `env:"FIBER_PORT" envDefault:"3000"`
	SyncJobIntervalMinutes int64  `env:"SYNC_JOB_INTERVAL_MINUTES" envDefault:"1"`
	LogLevel               string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath            string `env:"LOG_FILE_PATH" envDefault:"./auth.log"`
	ConsoleLogEnabled      bool   `env:"CONSOLE_LOG_ENABLED" envDefault:"true"`
	AES256EncryptionKey    string `env:"AES256_ENCRYPTION_KEY" envDefault:"EBCpGxd23e5CJSC2DpgT8w1tcDtiAaqE"`
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Init() (*AppConfig, error) {
	var config = AppConfig{}
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	config.LogFilePath = fmt.Sprintf("%s_%s.log", strings.Split(config.LogFilePath, ".log")[0], time.Now().UTC().Format(time.RFC3339))
	return &config, nil
}
