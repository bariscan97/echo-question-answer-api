package auth

import (
	controller "articles-api/controllers/auth"
	repository "articles-api/database/repository/user"
	service "articles-api/services/user"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func AuthRouter(e *echo.Echo, pool *pgxpool.Pool) {

	user_repo := repository.NewUserRepo(pool)
	user_service := service.NewAuthService(user_repo)
	auth_controller := controller.NewAuthController(user_service)

	auth := e.Group("/auth")

	auth.POST("/login", auth_controller.Login)
	auth.POST("/register", auth_controller.Register)
}
