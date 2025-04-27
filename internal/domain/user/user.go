package user

import "time"

type User struct {
	// ID será o hex do ObjectID gerado pelo Mongo
	ID          string     `bson:"_id,omitempty" json:"id"`
	Name        string     `bson:"name"       json:"name"`
	Username    string     `bson:"username"   json:"username"`
	Email       string     `bson:"email"      json:"email"`
	Password    string     `bson:"password"   json:"password"`
	Status      int        `bson:"status"     json:"status"`
	Photo       string     `bson:"photo,omitempty" json:"photo,omitempty"`
	RayDistance int        `bson:"ray_distance"     json:"ray_distance"`
	Level       int        `bson:"level"            json:"level"`
	CreatedAt   time.Time  `bson:"created_at"       json:"created_at"`
	ModifiedAt  time.Time  `bson:"modified_at"      json:"modified_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
