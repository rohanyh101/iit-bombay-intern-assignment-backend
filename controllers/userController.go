package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/roh4nyh/iit_bombay/database"
	helper "github.com/roh4nyh/iit_bombay/helpers"
	"github.com/roh4nyh/iit_bombay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserCollectionName = "users"
)

var userValidate = validator.New()
var UserCollection *mongo.Collection = database.OpenCollection(DatabaseName, UserCollectionName)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword, foundUserPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(foundUserPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintln("username or password is incorrect")
		check = false
	}

	return check, msg
}

func UserSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := userValidate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if *user.Role != models.ROLE_LIBRARIAN && *user.Role != models.ROLE_MEMBER {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role provided. Must be either 'LIBRARIAN' or 'MEMBER'."})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this user already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		*user.IsActive = false

		token, _ := helper.GenerateUserToken(*user.Username, user.UserID, *user.Role, *user.IsActive)
		user.Token = &token

		resultInsertionNumber, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintln("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusCreated, resultInsertionNumber)
	}
}

func UserLogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		token, err := helper.GenerateUserToken(*foundUser.Username, foundUser.ID.Hex(), *foundUser.Role, *foundUser.IsActive)
		if err != nil || token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		helper.UpdateUserToken(token, foundUser.UserID)

		err = UserCollection.FindOne(ctx, bson.M{"_id": foundUser.ID}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
