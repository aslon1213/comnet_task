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
	utilshelpers "github.com/aslon1213/comnet_task/internal/app/utils-helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserHandlers struct contains context and db connection
type UserHandlers struct {
	ctx context.Context
	db  *sql.DB
}

// initialize user handlers with context and db connection
func New(ctx context.Context, db *sql.DB) *UserHandlers {
	return &UserHandlers{
		ctx: ctx,
		db:  db,
	}
}

// POST - user/register - handler.
// input should contain age(int), name(string), login(string) and password(string) in form data.
// when succesfull returns 201 status code and message.
func (u *UserHandlers) Register(c *gin.Context) {
	os.Setenv("TZ", "Asia/Tashkent")
	// fmt.Println("User is going to be registered")
	var user models.User

	// get form data
	a := c.PostForm("age")
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
	// perform checks for valid data
	if err := utilshelpers.Check_for_route_names(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if user exists
	rows, err := u.db.QueryContext(u.ctx, "SELECT * FROM users WHERE login = ?", user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.JSON(http.StatusFound, gin.H{
			"error": "user already exists",
		})
		return
	}

	// create user
	tx, err := u.db.BeginTx(u.ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	stmt, err := tx.Prepare("INSERT INTO users(login, password, name, age) VALUES(?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	_, err = stmt.Exec(user.Login, user.Password, user.Name, user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer stmt.Close()
	err = tx.Commit()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// if succesfull return message with 200 code
	c.JSON(http.StatusCreated, gin.H{
		"message": "user has registered successfully",
	})
}

// GET - /user/auth - handler.
// input should contains login and password in json format.
// when succesfull sets cookie SESSTOKEN with jswt token and returns 200 status code.
func (u *UserHandlers) Auth(c *gin.Context) {
	// get json data and bind it with struct
	var user_input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user_input); err != nil {
		// if not valid return error
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// search for user
	rows, err := u.db.QueryContext(u.ctx, "SELECT id, login, password FROM users WHERE login = ?", user_input.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while preparing query --- " + err.Error(),
		})
		return
	}
	user := models.User{}
	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	defer rows.Close()
	// hash password and compare with user's hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_input.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong password or login",
		})
		return
	}

	// expire tiem 1 day
	// time.Now().In(time.FixedZone("UTC+5", 5*60*60))
	// set timezone to UTC+5
	// os.Setenv("TZ", "Asia/Tashkent")
	location, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	expire_time := time.Now().In(location).Add(24 * time.Hour)
	// create jwt token with expire_time and user info
	token_string, err := utilshelpers.CreateSessionCookieToken(user, expire_time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// set cookie SESSTOKEN
	c.SetCookie("SESSTOKEN", token_string, int(expire_time.Unix())-int(time.Now().In(location).Unix()), "/", "", false, true)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Welcome " + user.Login,
	})
}

// GET - user/:name - handler.
// returns user info by name.
// sample output contains age, name and id of user.
func (u *UserHandlers) GetUserByName(c *gin.Context) {

	username := c.Param("name")
	// qeury
	rows, err := u.db.QueryContext(u.ctx, "SELECT id, login, name, age FROM users WHERE name = ?", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// append all results to slice
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
	output := []gin.H{}
	for _, user := range users {
		output = append(output, gin.H{
			"id": user.ID,

			"age":  user.Age,
			"name": user.Name})

	}
	c.JSON(200, output)
}

// POST /user/phone - handler.
// creates phone for user.
// input should be in json format and contain phone(string <= 12), is_fax(bool) and description(string) fileds.
// returns 201 status code and inserted id.
func (u *UserHandlers) CreateUserPhone(c *gin.Context) {

	var phone_input struct {
		Phone       string `json:"phone"`
		IsFax       bool   `json:"is_fax"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&phone_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	// fmt.Println(phone_input)

	// no need to check
	user_id, ok := c.Get("user_id")
	// fmt.Println("USERID", user_id)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id not found",
		})
		return
	}

	tx, err := u.db.BeginTx(u.ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	stmt, err := tx.PrepareContext(u.ctx, "INSERT INTO phones(user_id, phone, is_fax, description) VALUES(?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while preparing query --- " + err.Error(),
		})
		return
	}

	res, err := stmt.Exec(user_id, phone_input.Phone, phone_input.IsFax, phone_input.Description)
	if err != nil {
		// stmt.Close()
		tx.Commit()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer stmt.Close()
	err = tx.Commit()
	if err != nil {
		// stmt.Close()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting ID --- " + err.Error(),
		})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(201, gin.H{
		"message":     "Phone created succesfully",
		"inserted_id": id,
	})
}

// home page handler.
func (u *UserHandlers) HomePage(c *gin.Context) {
	fmt.Println(os.Getenv("SIGNING_SECRET"))
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

// GET /user/phone?q= - handler.
// query phones.
// should contain q query param.
func (uh *UserHandlers) GetPhonesByQuery(c *gin.Context) {

	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "q is required",
		})
		return
	}

	rows, err := uh.db.QueryContext(uh.ctx, "SELECT user_id, phone, is_fax, description FROM phones WHERE phone LIKE ?", "%"+q+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(200, output)

}

// PUT - user/phone - handler.
// updates phone info - should contain phone_id(int),description(string), is_fax(bool) and phone(string) in json body.
// returns 200 status code and message.
func (uh *UserHandlers) UpdatePhone(c *gin.Context) {

	phone_input := models.PhoneUpdateInput{}
	if err := c.ShouldBindJSON(&phone_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	// fmt.Println(phone_input)

	// create statement
	res, err := uh.db.ExecContext(uh.ctx, "UPDATE phones SET phone = ?, is_fax = ?, description = ? WHERE id = ?", phone_input.Phone, phone_input.IsFax, phone_input.Description, phone_input.PhoneId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	rows_affected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if rows_affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "phone not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "phone updated successfully",
	})

}

// DELETE - user/phone/:phone_id - handler.
// deletes phone by phone_id.
// returns 200 status code and message if succesfull.
func (uh *UserHandlers) DeletePhone(c *gin.Context) {

	phone_id := c.Param("phone_id")
	if phone_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone_id is required",
		})
		return
	}

	res, err := uh.db.ExecContext(uh.ctx, "DELETE FROM phones WHERE id = ?", phone_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "phone not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "phone deleted successfully",
	})

}
