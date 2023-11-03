package playwright

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

type Service struct {
	headless bool
}

func NewService(headless bool) *Service {
	return &Service{headless: headless}
}

func (s *Service) Init() (*playwright.Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	err = playwright.Install(&playwright.RunOptions{Verbose: false})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := pw.Stop(); err != nil {
			logrus.WithError(err).Fatal("error during playwright stop")
		}
	}()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(s.headless),
	})
	if err != nil {
		return nil, err
	}
	return &browser, nil
}
