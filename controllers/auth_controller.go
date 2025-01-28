package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Emibrown/E-commerce-API/config"
	"github.com/Emibrown/E-commerce-API/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminRegisterInput struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	AdminSecret string `json:"admin_secret" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   RegisterInput     true  "Register Input"
// @Success      201   {object} RegisterResponse
// @Failure      400   {object} ErrorResponse
// @Router       /api/auth/register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Unable to hash password"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPwd),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "User registered successfully",
		User: UserPayload{
			ID:      user.ID,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
	})
}

// Login godoc
// @Summary      Login a user
// @Description  Logs in an existing user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   LoginInput  true  "Login Input"
// @Success      200   {object} LoginResponse
// @Failure      400,401 {object} ErrorResponse
// @Router       /api/auth/login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// RegisterAdmin godoc
// @Summary      Register a new admin user
// @Description  Creates a new admin user if a valid admin secret is provided
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   AdminRegisterInput  true  "Admin Register Input"
// @Success      201   {object} RegisterResponse
// @Failure      400,401 {object} ErrorResponse
// @Router       /api/auth/register-admin [post]
func RegisterAdmin(c *gin.Context) {
	var input AdminRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check admin secret
	if input.AdminSecret != os.Getenv("ADMIN_SECRET") {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid admin secret"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Unable to hash password"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPwd),
		IsAdmin:  true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Unable to create admin user"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "Admin user registered successfully",
		User: UserPayload{
			ID:      user.ID,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
	})
}
