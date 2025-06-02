package infrastructure

import (
	"context"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserRepository implements domain.UserRepository
type MongoUserRepository struct {
	collection *mongo.Collection
}

// MongoLogRepository implements domain.LogRepository
type MongoLogRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new instance of MongoUserRepository
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		collection: collection,
	}
}

// NewMongoLogRepository creates a new instance of MongoLogRepository
func NewMongoLogRepository(collection *mongo.Collection) *MongoLogRepository {
	return &MongoLogRepository{
		collection: collection,
	}
}

// Create implements domain.UserRepository
func (r *MongoUserRepository) Create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// FindByEmail implements domain.UserRepository
func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

// FindByName implements domain.UserRepository
func (r *MongoUserRepository) FindByName(ctx context.Context, name string) (domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	return user, err
}

// FindByID implements domain.UserRepository
func (r *MongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"deleted_at": nil,
	}).Decode(&user)
	return user, err
}

// FindAll implements domain.UserRepository
func (r *MongoUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Update implements domain.UserRepository
func (r *MongoUserRepository) Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":        id,
			"deleted_at": nil,
		},
		bson.M{"$set": update},
	)
	return err
}

// Delete implements domain.UserRepository
func (r *MongoUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":        id,
			"deleted_at": nil,
		},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
	)
	return err
}

// Count implements domain.UserRepository
func (r *MongoUserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// Create implements domain.LogRepository
func (r *MongoLogRepository) Create(ctx context.Context, log domain.RequestLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	return err
}
