package userHandlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aslon1213/comnet_task/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserHandlers) Register(c *gin.Context) {
	fmt.Println("User is going to be registered")
	var user models.User

	// get form data

	// fmt.Println(user_input)
	// fmt.Println()
	fmt.Println("AGE:", c.PostForm("age"))
	fmt.Println("NAME:", c.PostForm("name"))
	fmt.Println("LOGIN:", c.PostForm("login"))
	fmt.Println("PASSWORD:", c.PostForm("password"))
	a := c.PostForm("age")
	// fmt.Println(a)
	if a == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	age, err := strconv.Atoi(a)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	password_hashed, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user.Age = age
	user.Name = c.PostForm("name")
	user.Login = c.PostForm("login")
	user.Password = string(password_hashed)
	// user.Email = ""
	// user.Phone = ""

	// user.FirstName = c.GetString("firstname")
	// user.LastName = c.GetString("lastname")
	fmt.Println(user)
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check if user exists

	rows, err := u.db.QueryContext(u.ctx, "SELECT * FROM users WHERE login = ?", user.Login)
	// := u.db.Query("SELECT * FROM users WHERE login = ?", user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.JSON(http.StatusFound, gin.H{
			"message": "user already exists",
		})
		return
	}

	// create user
	// tx, err := u.db.Begin()
	tx, err := u.db.BeginTx(u.ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	stmt, err := tx.Prepare("INSERT INTO users(login, password, name, age) VALUES(?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}
	res, err := stmt.Exec(user.Login, user.Password, user.Name, user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer stmt.Close()
	err = tx.Commit()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(res.LastInsertId())
	c.JSON(http.StatusCreated, gin.H{
		"message": "user has registered successfully",
		"result":  res,
	})
}

func (u *UserHandlers) Auth(c *gin.Context) {
	var user_input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(user_input)
	rows, err := u.db.QueryContext(u.ctx, "SELECT id, login, password FROM users WHERE login = ?", user_input.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := models.User{}
	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	defer rows.Close()
	fmt.Println(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_input.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong password",
		})
		return
	}

	token_string, err := CreateSessionCookieToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.SetCookie("token", token_string, 24*3600, "/", "localhost", false, true)

	// check if user exists

	c.JSON(200, gin.H{
		"message": "Welcome User",
	})
}

func CreateSessionCookieToken(user models.User) (string, error) {

	// one day for expiration
	expirations_time := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Login,
		"user_id":  user.ID,
		"expires":  expirations_time,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_SECRET")))
	if err != nil {
		return "", err
	}

	fmt.Println(tokenString)
	return tokenString, nil
}

func (u *UserHandlers) GetUserByName(c *gin.Context) {

	username := c.Param("name")
	fmt.Println(username)
	rows, err := u.db.QueryContext(u.ctx, "SELECT id, login, name, age FROM users WHERE name = ?", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	users := []models.User{}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.Login, &user.Name, &user.Age)
		users = append(users, user)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	fmt.Println(users)
	output := []gin.H{}
	for _, user := range users {
		output = append(output, gin.H{
			"id": user.ID,

			"age":  user.Age,
			"name": user.Name})

	}
	c.JSON(200, output)
}

func (u *UserHandlers) CreateUserPhone(c *gin.Context) {

	var phone_input struct {
		Phone       string `json:"phone"`
		IsFax       bool   `json:"is_fax"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&phone_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	fmt.Println(phone_input)

	// no need to check
	user_id, ok := c.Get("user_id")
	fmt.Println("USERID", user_id)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id not found",
		})
		return
	}

	tx, err := u.db.BeginTx(u.ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	stmt, err := tx.PrepareContext(u.ctx, "INSERT INTO phones(user_id, phone, is_fax, description) VALUES(?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while preparing query --- " + err.Error(),
		})
		return
	}

	res, err := stmt.Exec(user_id, phone_input.Phone, phone_input.IsFax, phone_input.Description)
	if err != nil {
		// stmt.Close()
		tx.Commit()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer stmt.Close()
	err = tx.Commit()
	if err != nil {
		// stmt.Close()

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while inserting ID --- " + err.Error(),
		})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(201, gin.H{
		"message":     "Phone created succesfully",
		"inserted_id": id,
	})
}

func (u *UserHandlers) HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (uh *UserHandlers) GetPhonesByQuery(c *gin.Context) {

	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "q is required",
		})
		return
	}

	rows, err := uh.db.QueryContext(uh.ctx, "SELECT user_id, phone, is_fax, description FROM phones WHERE phone LIKE ?", "%"+q+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	output := []models.PhoneOutput{}
	defer rows.Close()
	for rows.Next() {
		phone := models.PhoneOutput{}
		err = rows.Scan(&phone.UserID, &phone.Phone, &phone.IsFax, &phone.Description)
		output = append(output, phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	c.JSON(200, output)

}

func (uh *UserHandlers) UpdatePhone(c *gin.Context) {

	phone_input := models.PhoneUpdateInput{}
	if err := c.ShouldBindJSON(&phone_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	fmt.Println(phone_input)

	// create statement
	res, err := uh.db.ExecContext(uh.ctx, "UPDATE phones SET phone = ?, is_fax = ?, description = ? WHERE id = ?", phone_input.Phone, phone_input.IsFax, phone_input.Description, phone_input.PhoneId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	rows_affected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if rows_affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "phone not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "phone updated successfully",
	})

}

func (uh *UserHandlers) DeletePhone(c *gin.Context) {

	phone_id := c.Param("phone_id")
	if phone_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "phone_id is required",
		})
		return
	}

	res, err := uh.db.ExecContext(uh.ctx, "DELETE FROM phones WHERE id = ?", phone_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "phone not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "phone deleted successfully",
	})

}
