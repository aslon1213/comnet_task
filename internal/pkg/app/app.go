package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aslon1213/comnet_task/internal/app/handlers"
	"github.com/aslon1213/comnet_task/internal/app/handlers/userHandlers"
	initializers "github.com/aslon1213/comnet_task/internal/app/initializers"
	middlewares "github.com/aslon1213/comnet_task/internal/app/middleware"
	authmiddleware "github.com/aslon1213/comnet_task/internal/app/middleware/AuthMiddleware"
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

	ctx := context.Background()
	app := &App{}

	db, err := initializers.Init()
	if err != nil {
		panic(err)
	}

	app.Db = db

	app.Gin = gin.Default()
	app.Midw = middlewares.New(authmiddleware.New(db, ctx))
	app.H = handlers.New(userHandlers.New(ctx, app.Db))

	// register routes
	// user routes
	// app.UserRoutes()
	user := app.Gin.Group("/user")
	//2 register - done
	user.POST("/register", app.H.UserHandlers.Register)
	//3 - done
	user.GET("/auth", app.H.UserHandlers.Auth)
	//4 - done - middleware done
	user.GET("/:name", app.Midw.Auth.AuthMiddleware, app.H.UserHandlers.GetUserByName)
	//5
	user.POST("/phone", app.Midw.Auth.AuthMiddleware, app.H.UserHandlers.CreateUserPhone)
	//6
	user.GET("/phone", app.Midw.Auth.AuthMiddleware, app.H.UserHandlers.GetPhonesByQuery)
	//7
	user.PUT("/phone", app.Midw.Auth.AuthMiddleware, app.H.UserHandlers.UpdatePhone)
	//8
	user.DELETE("/phone/:phone_id", app.Midw.Auth.AuthMiddleware, app.H.UserHandlers.DeletePhone)

	app.Gin.GET("/", app.H.UserHandlers.HomePage)

	return app
}

func (a *App) Run() {

	fmt.Println("Running server on port 8080")

	a.Gin.Run(":8080")
}

func (a *App) UserRoutes() {

}
