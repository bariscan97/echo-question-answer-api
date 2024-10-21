package routes

import (
	auth "articles-api/routes/auth"
	post "articles-api/routes/post"
	user "articles-api/routes/user"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Routers struct {
	e    *echo.Echo
	pool *pgxpool.Pool
}

func NewRouters(e *echo.Echo, pool *pgxpool.Pool) *Routers {
	return &Routers{
		e:    e,
		pool: pool,
	}
}

func (routers *Routers) InitRouter() {
	auth.AuthRouter(routers.e, routers.pool)
	post.PostRouter(routers.e, routers.pool)
	user.UserRouter(routers.e, routers.pool)
}

func (routers *Routers) RoutersRun(Addr string) error {
	return routers.e.Start(Addr)
}
