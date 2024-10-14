package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/roh4nyh/iit_bombay/helpers"
	"github.com/roh4nyh/iit_bombay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookValidate = validator.New()

func AddBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var book models.Book
		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := bookValidate.Struct(book)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		book.CreatedAt = time.Now()
		book.UpdatedAt = time.Now()

		_, err := BookCollection.InsertOne(ctx, book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while adding book"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "book added successfully"})
	}
}

func UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var book models.Book

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateObj := bson.M{}

		if book.Title != nil {
			updateObj["title"] = book.Title
		}

		if book.Author != nil {
			updateObj["author"] = book.Author
		}

		if book.Status != nil {
			updateObj["status"] = book.Status
		}

		if book.Qty != 0 {
			updateObj["qty"] = book.Qty
		}

		if book.ISBN != nil {
			updateObj["isbn"] = book.ISBN
		}

		updateObj["updated_at"] = time.Now()

		filter := bson.M{"isbn": bson.M{"$eq": isbn}}
		update := bson.M{"$set": updateObj}

		_, err := BookCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating book"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "book updated successfully"})
	}
}

func DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		filter := bson.M{"isbn": bson.M{"$eq": isbn}}

		_, err := BookCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting book"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var users []models.User

		cursor, err := UserCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing users"})
			return
		}

		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding user data"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no user available"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIdStr := c.Param("user_id")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var user models.User
		err = UserCollection.FindOne(ctx, bson.M{"_id": memberId, "role": models.ROLE_MEMBER}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if user.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

		count, err := UserCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this user already exists"})
			return
		}

		isActive := false
		mRole := models.ROLE_MEMBER

		user.Role = &mRole
		user.IsActive = &isActive
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		password := HashPassword(*user.Password)
		user.Password = &password

		token, _ := helpers.GenerateUserToken(*user.Username, user.UserID, *user.Role, *user.IsActive)
		user.Token = &token

		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user occurred while adding customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user added successfully"})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("user_id")
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateObj := bson.M{}

		if user.Username != nil {
			updateObj["username"] = user.Username
		}

		if user.Role != nil {
			updateObj["role"] = user.Role
		}

		if user.IsActive != nil {
			updateObj["is_active"] = user.IsActive
		}

		if user.Password != nil {
			updateObj["password"] = user.Password
		}

		updateObj["updated_at"] = time.Now()

		filter := bson.M{"_id": bson.M{"$eq": userId}}
		update := bson.M{"$set": updateObj}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
	}
}

func DeActivateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIdStr := c.Param("user_id")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err = UserCollection.FindOne(ctx, bson.M{"_id": memberId, "role": models.ROLE_MEMBER}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if user.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		updateObj := bson.M{}

		updateObj["is_active"] = false

		updateObj["updated_at"] = time.Now()

		filter := bson.M{"_id": bson.M{"$eq": memberId}}
		update := bson.M{"$set": updateObj}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user de-activated successfully"})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("user_id")
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Check if the user is a librarian
		var user models.User
		err = UserCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
			}
			return
		}

		if *user.Role == models.ROLE_LIBRARIAN {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete a librarian"})
			return
		}

		filter := bson.M{"_id": userId}
		result, err := UserCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while deleting user"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func GetActiveUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var users []models.User

		cursor, err := UserCollection.Find(ctx, bson.M{"is_active": true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing customers"})
			return
		}

		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding users data"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusOK, []models.User{})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetNonActiveUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var users []models.User

		cursor, err := UserCollection.Find(ctx, bson.M{"is_active": false})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing users"})
			return
		}

		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding users data"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no users available"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetTransactionHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("user_id")
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var borrowHistory []models.BorrowHistory

		cursor, err := BorrowHistoryCollection.Find(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing borrowed books"})
			return
		}

		if err = cursor.All(ctx, &borrowHistory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding borrowed books data"})
			return
		}

		if len(borrowHistory) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no borrowed books available"})
			return
		}

		c.JSON(http.StatusOK, borrowHistory)
	}
}
