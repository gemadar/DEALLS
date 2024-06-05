package controllers

import (
	"echo-mongo-api/models"
	"echo-mongo-api/responses"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// View controller
func View(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var data models.User
	var user *models.User
	defer cancel()

	// Get token from request header
	authHeader := c.Request().Header.Get("Authorization")

	// Get data from token
	claims, claimErr := jwt.ParseWithClaims(authHeader, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey()), nil
	})

	if claims.Valid {
		claim := claims.Claims.(*models.Claims)

		userCollection.FindOne(ctx, bson.M{"username": claim.Username}).Decode(&data)

		// Validate user status and view in a day
		if !data.Premium && data.Viewed > 10 {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: "That's all for today! Come back tomorrow or get premium!"})
		}

		if data.Liked == nil {
			data.Liked = []string{}
		}

		// Filter user
		filter := bson.M{"$and": []bson.M{{"username": bson.M{"$nin": data.Liked}}, {"username": bson.M{"$nin": data.Passed}}, {"username": bson.M{"$ne": claim.Username}}}}
		option := options.FindOne().SetProjection(bson.M{"password": 0, "viewed": 0, "liked": 0, "id": 0, "premium": 0, "passed": 0})

		err := userCollection.FindOne(ctx, filter, option).Decode(&user)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		if user == nil {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: "That's all!"})
		}

	} else {
		return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: claimErr.Error()})
	}

	return c.JSON(http.StatusOK, responses.DataResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

// Get premium controller
func GetPremium(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	authHeader := c.Request().Header.Get("Authorization")

	claims, claimErr := jwt.ParseWithClaims(authHeader, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey()), nil
	})

	if claims.Valid {
		claim := claims.Claims.(*models.Claims)

		update := bson.M{"premium": true}
		_, err := userCollection.UpdateOne(ctx, bson.M{"username": claim.Username}, bson.M{"$set": update})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

	} else {
		return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: claimErr.Error()})
	}

	return c.JSON(http.StatusOK, responses.MessageResponse{Status: http.StatusOK, Message: "success", Data: "Congratulation! You are a premium user now!"})
}

// Like controller
func Like(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	authHeader := c.Request().Header.Get("Authorization")

	claims, claimErr := jwt.ParseWithClaims(authHeader, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey()), nil
	})

	if claims.Valid {
		claim := claims.Claims.(*models.Claims)

		update := bson.M{"liked": user.Username}
		increment := bson.M{"viewed": 1}
		_, err := userCollection.UpdateOne(ctx, bson.M{"username": claim.Username}, bson.M{"$push": update, "$inc": increment})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

	} else {
		return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: claimErr.Error()})
	}

	return c.JSON(http.StatusOK, responses.MessageResponse{Status: http.StatusOK, Message: "success", Data: "Liked! Hope it's a match!"})
}

// Pass controller
func Pass(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.MessageResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	authHeader := c.Request().Header.Get("Authorization")

	claims, claimErr := jwt.ParseWithClaims(authHeader, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey()), nil
	})

	if claims.Valid {
		claim := claims.Claims.(*models.Claims)

		update := bson.M{"passed": user.Username}
		increment := bson.M{"viewed": 1}
		_, err := userCollection.UpdateOne(ctx, bson.M{"username": claim.Username}, bson.M{"$push": update, "$inc": increment})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

	} else {
		return c.JSON(http.StatusInternalServerError, responses.MessageResponse{Status: http.StatusInternalServerError, Message: "error", Data: claimErr.Error()})
	}

	return c.JSON(http.StatusOK, responses.MessageResponse{Status: http.StatusOK, Message: "success", Data: "Passed!"})
}
