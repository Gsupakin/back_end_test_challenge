package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user User) (primitive.ObjectID, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByName(ctx context.Context, name string) (User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Count(ctx context.Context) (int64, error)
}

// LogRepository defines the interface for request log operations
type LogRepository interface {
	Create(ctx context.Context, log RequestLog) error
}
