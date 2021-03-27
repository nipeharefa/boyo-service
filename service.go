package boyo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type (
	boyoService struct {
		e           *echo.Echo
		serviceName string
		vip         *viper.Viper
		connStr     string
		db          *sqlx.DB
	}
)

func WithDB(key string) ServiceOptions {
	return func(s *boyoService) {
		conn := s.vip.GetString(key)
		s.connStr = conn
	}
}

type ServiceOptions func(*boyoService)

func NewBoyoService(name string, vip *viper.Viper, opts ...ServiceOptions) *boyoService {
	b := &boyoService{}
	b.vip = vip
	b.db = new(sqlx.DB)

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

	if err := s.connectDB(); err != nil {
		return err
	}

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

func (s *boyoService) GetDB() *sqlx.DB {

	return s.db
}

func (s *boyoService) GetRouter() *echo.Echo {

	return s.e
}

func (s *boyoService) connectDB() error {

	if s.connStr != "" {
		db, err := sqlx.Open("postgres", s.connStr)
		if err != nil {
			return err
		}

		if err := db.Ping(); err != nil {
			return err
		}

		*s.db = *db
	}

	return nil
}
