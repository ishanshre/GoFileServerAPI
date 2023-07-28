package repository

import "github.com/ishanshre/GoFileServerAPI/internals/pkg/models"

type Repository interface {

	//user interface
	GetUsers(limit, page int) ([]*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UsernameExists(username string) error
	EmailExists(email string) error
	DeleteUser(username string) error
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(username string, update *models.User) (*models.User, error)
}
