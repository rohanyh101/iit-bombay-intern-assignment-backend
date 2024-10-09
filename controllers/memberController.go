package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roh4nyh/iit_bombay/database"
	"github.com/roh4nyh/iit_bombay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName                = "Cluster0"
	BookCollectionName          = "books"
	BorrowHistoryCollectionName = "borrowHistory"
)

var BookCollection *mongo.Collection = database.OpenCollection(DatabaseName, BookCollectionName)
var BorrowHistoryCollection *mongo.Collection = database.OpenCollection(DatabaseName, BorrowHistoryCollectionName)

func GetBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var books []models.Book
		cursor, err := BookCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing users"})
			return
		}

		if err = cursor.All(ctx, &books); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding user data"})
			return
		}

		if len(books) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no users available"})
			return
		}

		c.JSON(http.StatusOK, books)
	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book models.Book
		err := BookCollection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found"})
			return
		}

		if book.Title == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found"})
			return
		}

		c.JSON(http.StatusOK, book)
	}
}

func DeActivateMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIdStr := c.GetString("uid")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// make sure user returns all borrowed books before deactivating
		var borrowedBooks []models.Book
		cursor, err := BookCollection.Find(ctx, bson.M{"borrowed_by": memberId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing borrowed books"})
			return
		}

		if err = cursor.All(ctx, &borrowedBooks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding borrowed books data"})
			return
		}

		if len(borrowedBooks) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "return all borrowed books before de-activating account"})
			return
		}

		// Deactivate user
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

func BorrowedBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIdStr := c.GetString("uid")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var borrowedBooks []models.Book

		cursor, err := BookCollection.Find(ctx, bson.M{"borrowed_by": memberId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing borrowed books"})
			return
		}

		if err = cursor.All(ctx, &borrowedBooks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding borrowed books data"})
			return
		}

		// Return an empty array if no books have been borrowed
		if len(borrowedBooks) == 0 {
			c.JSON(http.StatusOK, []models.Book{})
			return
		}

		c.JSON(http.StatusOK, borrowedBooks)
	}
}

func BorrowBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")

		memberIdStr := c.GetString("uid")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book models.Book
		err = BookCollection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found"})
			return
		}

		if book.Status == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found"})
			return
		}

		if *book.Status == models.STATUS_OUT_OF_STOCK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book is out of stock"})
			return
		}

		updateObj := bson.M{}

		updateObj["status"] = models.STATUS_BORROWED
		updateObj["borrowed_by"] = memberId
		updateObj["updated_at"] = time.Now()
		updateObj["qty"] = book.Qty - 1

		if updateObj["qty"].(int) <= 0 {
			updateObj["status"] = models.STATUS_OUT_OF_STOCK
		}

		filter := bson.M{"isbn": bson.M{"$eq": isbn}}
		update := bson.M{"$set": updateObj}

		_, err = BookCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating book"})
			return
		}

		bbook := models.STATUS_BORROWED
		borrowHistory := models.BorrowHistory{
			UserID:     memberId,
			BookID:     book.ID,
			BorrowedAt: time.Now(),
			Status:     &bbook,
		}

		_, err = BorrowHistoryCollection.InsertOne(ctx, borrowHistory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting borrow history"})
			return
		}

		// Activate user if deactivated
		updateObj = bson.M{}

		updateObj["is_active"] = true
		updateObj["updated_at"] = time.Now()

		filter = bson.M{"_id": bson.M{"$eq": memberId}}
		update = bson.M{"$set": updateObj}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "book borrowed successfully"})
	}
}

func ReturnBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")

		memberIdStr := c.GetString("uid")
		memberId, err := primitive.ObjectIDFromHex(memberIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book models.Book
		err = BookCollection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found in your borrowed list"})
			return
		}

		if book.Status == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found in your borrowed list"})
			return
		}

		if *book.Status == models.STATUS_AVAILABLE {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book is not yet borrowed by you"})
			return
		}

		// Update book
		updateObj := bson.M{}

		updateObj["status"] = models.STATUS_AVAILABLE
		updateObj["borrowed_by"] = nil
		updateObj["updated_at"] = time.Now()
		updateObj["qty"] = book.Qty + 1

		filter := bson.M{"isbn": bson.M{"$eq": isbn}}
		update := bson.M{"$set": updateObj}

		_, err = BookCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating book"})
			return
		}

		// Update borrow history
		var borrowHistory models.BorrowHistory
		err = BorrowHistoryCollection.FindOne(ctx, bson.M{"book_id": book.ID, "user_id": memberId, "status": models.STATUS_BORROWED}).Decode(&borrowHistory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "borrow history not found"})
			return
		}

		updateObj = bson.M{}

		updateObj["status"] = models.STATUS_RETURNED
		updateObj["returned_at"] = time.Now()

		filter = bson.M{"_id": bson.M{"$eq": borrowHistory.ID}}
		update = bson.M{"$set": updateObj}

		_, err = BorrowHistoryCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating borrow history"})
			return
		}

		// // Deactivate user if all books are returned
		// var borrowedBooks []models.Book
		// cursor, err := BookCollection.Find(ctx, bson.M{"borrowed_by": memberId})
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing borrowed books"})
		// 	return
		// }

		// if err = cursor.All(ctx, &borrowedBooks); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding borrowed books data"})
		// 	return
		// }

		// if len(borrowedBooks) == 0 {
		// 	updateObj = bson.M{}

		// 	updateObj["is_active"] = false

		// 	updateObj["updated_at"] = time.Now()

		// 	filter = bson.M{"_id": bson.M{"$eq": memberId}}
		// 	update = bson.M{"$set": updateObj}

		// 	_, err = UserCollection.UpdateOne(ctx, filter, update)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user"})
		// 		return
		// 	}
		// }

		c.JSON(http.StatusOK, gin.H{"message": "book returned successfully"})
	}
}
