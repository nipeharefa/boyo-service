package boyo

import (
	"github.com/jmoiron/sqlx"
)

type Service interface {
	GetDB() *sqlx.DB
	// GetRouter() *echo.Echo
	Run() error
}
