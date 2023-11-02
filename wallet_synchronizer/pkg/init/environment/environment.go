package environment

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
	"time"
)

type AppConfig struct {
	PlaywrightHeadless      bool   `env:"PLAYWRIGHT_HEADLESS" envDefault:"true"`
	AlchemyAPIKey           string `env:"ALCHEMY_API_KEY" envDefault:"owUCVigVvnHA63o0C6mh3yrf3jxMkV7b"`
	FiberPort               int64  `env:"FIBER_PORT" envDefault:"3000"`
	BrowserPath             string `env:"BROWSER_PATH" envDefault:"/usr/bin/brave-browser"`
	ScrapeIntervalMinutes   int64  `env:"SCRAPE_INTERVAL_MINUTES" envDefault:"1"`
	LogLevel                string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath             string `env:"LOG_FILE_PATH" envDefault:"./wallet_synchronizer.log"`
	ConsoleLogEnabled       bool   `env:"CONSOLE_LOG_ENABLED" envDefault:"true"`
	OwnWallet               string `env:"OWN_WALLET" envDefault:"0x251e929c9b5887e2d30b38dec708b7e40fb8c492"`
	ServerSSLCertFilepath   string `env:"SERVER_SSLCERT_FILEPATH,required"`
	ServerSSLKeyFilepath    string `env:"SERVER_SSLKEY_FILEPATH,required"`
	RSA256PublicKeyFilepath string `env:"RSA_256_PUBLIC_KEY_FILEPATH,required"`
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
