package user

import (
	controller "articles-api/controllers/user"
	repository "articles-api/database/repository/user"
	service "articles-api/services/user"
	"articles-api/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	
)

func UserRouter(e *echo.Echo, pool *pgxpool.Pool) {
	
	user_repo := repository.NewUserRepo(pool)
	user_service := service.NewAuthService(user_repo)
	user_controller := controller.NewUserController(user_service)
	
	users := e.Group("/users")
	
	users.Use(middleware.AuthMiddleware)
	
	users.GET("/", user_controller.GetAllUsers)
	users.GET("/me", user_controller.GetMe)
	users.POST("/uploadImage",user_controller.UpdateUserProfileImgById)
	users.GET("/:id", user_controller.GetUserById)
	users.DELETE("/:id", user_controller.DeleteUserById)
	users.PUT("/:id", user_controller.UpdateUserById)
	users.PUT("/:id/block", user_controller.BlockUserById)
	users.GET("/:id/blocklist" ,user_controller.GetMyBlockList)
	users.PUT("/:id/follow" , user_controller.FollowUserById)
	users.GET("/:id/following" ,user_controller.GetUserFollowingListById)
	users.GET("/:id/followers" ,user_controller.GetUserFollowersListById)
}

