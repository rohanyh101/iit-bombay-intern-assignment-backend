package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ROLE_LIBRARIAN      = "LIBRARIAN"
	ROLE_MEMBER         = "MEMBER"
	STATUS_AVAILABLE    = "AVAILABLE"
	STATUS_OUT_OF_STOCK = "OUT_OF_STOCK"
	STATUS_BORROWED     = "BORROWED"
	STATUS_RETURNED     = "RETURNED"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  *string            `bson:"username" json:"username" validate:"required"`
	Password  *string            `bson:"password" json:"password" validate:"required,min=4"`
	Role      *string            `bson:"role" json:"role" validate:"required,eq=LIBRARIAN|eq=MEMBER"`
	IsActive  *bool              `bson:"is_active" json:"is_active"` // Marks if user is active or deleted
	Token     *string            `bson:"token,omitempty" json:"token,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ISBN   *string            `bson:"isbn" json:"isbn" validate:"required"`
	Title  *string            `bson:"title" json:"title" validate:"required"`
	Author *string            `bson:"author" json:"author" validate:"required"`
	Status *string            `bson:"status" json:"status" validate:"required,eq=AVAILABLE|eq=OUT_OF_STOCK"`
	Qty    int                `bson:"qty" json:"qty" validate:"required"`
	// BorrowedBy *primitive.ObjectID `bson:"borrowed_by,omitempty" json:"borrowed_by,omitempty"` // User ID of the member borrowing the book
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	BookID    string    `bson:"book_id,omitempty" json:"book_id,omitempty"`
}

type BorrowHistory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"` // The member who borrowed the book
	BookID     primitive.ObjectID `bson:"book_id" json:"book_id"` // The book being borrowed
	BorrowedAt time.Time          `bson:"borrowed_at" json:"borrowed_at"`
	ReturnedAt time.Time          `bson:"returned_at,omitempty" json:"returned_at,omitempty"` // Nullable if not yet returned
	Status     string             `bson:"status,omitempty" json:"status,omitempty" validate:"eq=RETURNED|eq=BORROWED"`
	BorrowID   string             `bson:"borrow_id,omitempty" json:"borrow_id,omitempty"`
}
