package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Username    string             `json:"username,omitempty" bson:"username,omitempty,unique"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	AccessLevel int                `json:"access_level,omitempty" bson:"access_level,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	LastLogin   time.Time          `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

type CreateUser struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty,unique"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type File struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FilePath   string             `json:"filePath,omitempty" bson:"file,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Extension  string             `json:"extension,omitempty" bson:"extension,omitempty"`
	Size       int64              `json:"size,omitempty" bson:"size,omitempty"`
	Uploader   *User              `json:"uplaoder,omitempty" bson:"uploader,omitempty"`
	UploadedAt time.Time          `json:"uploaded_at,omitempty" bsom:"uploaded_at,omitempty"`
	ModifiedAt time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}
