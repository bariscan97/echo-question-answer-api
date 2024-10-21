package post

import (
	controller "articles-api/controllers/post"
	repository "articles-api/database/repository/post"
	"articles-api/middleware"
	service "articles-api/services/post"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func PostRouter(e *echo.Echo, pool *pgxpool.Pool) {

	post_repo := repository.NewUserRepo(pool)
	post_service := service.NewPostService(post_repo)
	post_controller := controller.NewAuthController(post_service)

	posts := e.Group("/posts")

	posts.Use(middleware.AuthMiddleware)

	posts.POST("/", post_controller.CreatePost)
	posts.GET("/all/:category", post_controller.GetAllPost)
	posts.GET("/:id", post_controller.GetSinglePostById)
	posts.DELETE("/:id", post_controller.DeletePostById)
	posts.PUT("/:id", post_controller.UpdatePostById)
	posts.POST("/:parentid", post_controller.CreateComment)
	posts.GET("/:id/comments", post_controller.GetPostCommentsById)
	posts.GET("/:userid/articles", post_controller.GetUserPostsById)
	posts.PUT("/:id/like", post_controller.LikePostById)
}
