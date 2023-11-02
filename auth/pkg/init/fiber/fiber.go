package fiber

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	ctx     context.Context
	port    int64
	apiList []Api
	cert    tls.Certificate
}

type Api struct {
	method  string
	path    string
	handler []fiber.Handler
}

func NewApi(
	method string,
	path string,
	handler []fiber.Handler) Api {
	return Api{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func NewService(
	cert tls.Certificate,
	ctx context.Context,
	port int64,
	apiList []Api,
) *Service {
	return &Service{
		ctx:     ctx,
		port:    port,
		apiList: apiList,
		cert:    cert,
	}
}

func (s *Service) Init() error {
	app := initializeFiberApp(s.apiList)
	app.Server()
	address := fmt.Sprintf(":%d", s.port)
	err := app.ListenTLSWithCertificate(address, s.cert)
	if err != nil {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			return err
		}
	}
	return nil
}

func initializeFiberApp(apiList []Api) *fiber.App {
	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
		ContextKey: "uuid",
	}))

	app.Server().WriteTimeout = 300 * time.Second
	app.Server().ReadTimeout = 300 * time.Second
	app.Server().ReadBufferSize = 100 * 1024 * 1024
	app.Server().MaxRequestBodySize = 100 * 1024 * 1024

	registerAPIRoutes(app, apiList)

	return app
}

func registerAPIRoutes(app *fiber.App, apiList []Api) {
	for _, api := range apiList {
		app.Add(api.method, api.path, api.handler...)
	}

}
