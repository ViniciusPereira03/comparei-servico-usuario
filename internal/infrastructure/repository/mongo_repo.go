package repository

import (
	"context"
	"time"

	"comparei-servico-usuario/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implementa o repository usando MongoDB
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository cria um novo MongoRepository
func NewMongoRepository(client *mongo.Client, dbName, collectionName string) *MongoRepository {
	coll := client.Database(dbName).Collection(collectionName)
	return &MongoRepository{collection: coll}
}

// CreateUser insere um novo usuário em MongoDB
func (r *MongoRepository) CreateUser(u *user.User) (*user.User, error) {
	u.Status = 1
	u.Level = 1
	u.RayDistance = 5
	now := time.Now()
	u.CreatedAt = now
	u.ModifiedAt = now
	// Insere o documento de usuário
	res, err := r.collection.InsertOne(context.Background(), u)
	if err != nil {
		return nil, err
	}
	// Converter ObjectID para string
	oid := res.InsertedID.(primitive.ObjectID)
	u.ID = oid.Hex()
	return u, nil
}

// GetUser busca um usuário por ID
func (r *MongoRepository) GetUser(id string) (*user.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var u user.User
	err = r.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}

// GetUsers lista todos os usuários ativos
func (r *MongoRepository) GetUsers(order string) ([]*user.User, error) {

	findOptions := options.Find()
	if order == "ranking" {
		findOptions.SetSort(bson.D{{Key: "level", Value: -1}})
	}

	cursor, err := r.collection.Find(context.Background(), bson.M{"status": 1}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var list []*user.User
	for cursor.Next(context.Background()) {
		var u user.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		list = append(list, &u)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateUser atualiza os dados de um usuário
func (r *MongoRepository) UpdateUser(u *user.User) error {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}

	set := bson.M{}

	if u.Name != "" {
		set["name"] = u.Name
	}
	if u.Username != "" {
		set["username"] = u.Username
	}
	if u.Email != "" {
		set["email"] = u.Email
	}
	if u.Photo != "" {
		set["photo"] = u.Photo
	}
	if u.Status > 0 {
		set["status"] = u.Status
	}
	if u.RayDistance > 0 {
		set["ray_distance"] = u.RayDistance
	}
	if u.Level > 0 {
		set["level"] = u.Level
	}

	if len(set) == 0 {
		return nil
	}

	update := bson.M{
		"$set": set,
	}

	_, err = r.collection.UpdateByID(context.Background(), oid, update)
	return err
}

// DeleteUser marca o usuário como inativo
func (r *MongoRepository) UpdateLevelUser(user_id string, level int) error {
	oid, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{
		"level":       level,
		"modified_at": time.Now(),
	}}
	_, err = r.collection.UpdateByID(context.Background(), oid, update)
	return err
}

// DeleteUser marca o usuário como inativo
func (r *MongoRepository) DeleteUser(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"status": 0, "deleted_at": time.Now()}}
	_, err = r.collection.UpdateByID(context.Background(), oid, update)
	return err
}

// GetUserByUsername busca usuário pelo username
func (r *MongoRepository) GetUserByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
