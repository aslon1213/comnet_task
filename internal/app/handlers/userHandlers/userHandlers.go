package userHandlers

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	ctx context.Context
	db  *sql.DB
}

func New(ctx context.Context, db *sql.DB) *UserHandlers {
	return &UserHandlers{
		ctx: ctx,
		db:  db,
	}
}

func (u *UserHandlers) HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (uh *UserHandlers) CreateUser(c *gin.Context) {

}
