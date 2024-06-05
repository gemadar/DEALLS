package controllers

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/models"
	"echo-mongo-api/responses"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

// Get token key from .env file
func tokenKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("JWT_SECRET")
}

// Generate token after sign in
func tokenGen(data models.User) string {
	claims := &models.Claims{
		Name:     data.Name,
		Username: data.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := []byte(tokenKey())
	tokenString, _ := token.SignedString(jwtKey)

	return tokenString
}

// Sign up controller
func SignUp(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// Validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	// Validate username uniqueness
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := userCollection.Indexes().CreateMany(context.Background(), indexes)
	println("a")
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: "Username already exist!"})
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Username: user.Username,
		Password: string(hashed),
		Name:     user.Name,
		Location: user.Location,
		Age:      user.Age,
		Gender:   user.Gender,
		Liked:    []string{},
		Passed:   []string{},
	}

	// Insert new user
	_, errs := userCollection.InsertOne(ctx, newUser)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: errs.Error()})
	}

	return c.JSON(http.StatusCreated, responses.MessageResponse{Status: http.StatusCreated, Message: "success", Data: "Welcome!"})
}

// Sign in controller
func SignIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var auth models.Auth
	var data models.User
	defer cancel()

	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	if validationErr := validate.Struct(&auth); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	// Get user data from DB
	err := userCollection.FindOne(ctx, bson.M{"username": auth.Username}).Decode(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(auth.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: "Wrong Credentials"})
	}

	token := tokenGen(data)

	return c.JSON(http.StatusOK, responses.MessageResponse{Status: http.StatusOK, Message: "success", Data: token})
}
