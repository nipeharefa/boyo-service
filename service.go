package boyo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type (
	boyoService struct {
		e           *echo.Echo
		serviceName string
		vip         *viper.Viper
	}
)

type ServiceOptions func(*boyoService)

func NewBoyoService(name string, vip *viper.Viper, opts ...ServiceOptions) *boyoService {
	b := &boyoService{}
	b.vip = vip

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(b)
		}
	}

	e := echo.New()
	e.HideBanner = true

	b.e = e

	return b
}

func (s boyoService) Run() error {
	port := 8000
	if p := s.vip.GetInt("app.port"); p != 0 {
		port = p
	}

	sb := strings.Builder{}
	fmt.Fprintf(&sb, ":%d", port)

	if err := s.e.Start(sb.String()); err != nil && err != http.ErrServerClosed {
		s.e.Logger.Error("shutting down the server")
		return err
	}

	return nil
}
