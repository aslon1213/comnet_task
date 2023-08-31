package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aslon1213/comnet_task/internal/app/handlers"
	"github.com/aslon1213/comnet_task/internal/app/handlers/userHandlers"
	"github.com/aslon1213/comnet_task/internal/app/initializers"
	middlewares "github.com/aslon1213/comnet_task/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	// middlewares
	Midw *middlewares.Middlewares
	H    *handlers.Handlers
	Db   *sql.DB
	// gin engine
	Gin *gin.Engine
}

func New() *App {

	// prepare gin engine
	initializers.Init()
	ctx := context.Background()
	app := &App{}
	db, err := initializers.ConnectSqlite()
	app.Db = db
	if err != nil {
		panic(err)
	}
	app.Gin = gin.Default()
	app.Midw = middlewares.New()
	app.H = handlers.New(userHandlers.New(ctx, app.Db))

	// register routes
	// user routes
	app.UserRoutes()

	app.Gin.GET("/", app.H.UserHandlers.HomePage)

	return app
}

func (a *App) Run() {

	fmt.Println("Running server on port 8080")

	a.Gin.Run(":8080")
}

func (a *App) UserRoutes() {
	a.Gin.GET("/users", a.H.UserHandlers.HomePage)
}
