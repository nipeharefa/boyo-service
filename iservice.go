package boyo

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Service interface {
	GetDB() *sqlx.DB
	GetRouter() *echo.Echo
	Run() error
}
