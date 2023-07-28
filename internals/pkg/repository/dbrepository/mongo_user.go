package dbrepository

import (
	"context"
	"errors"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mongoDbRepo) GetUsers(limit, page int64) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit
	opt := options.FindOptions{}
	opt.SetLimit(limit)
	opt.SetSkip(skip)
	query := bson.M{}

	result, err := m.db.GetUserCollection().Find(ctx, query, &opt)
	if err != nil {
		return nil, errors.New("error in fetching all users")
	}
	defer result.Close(ctx)
	users := []*models.User{}
	for result.Next(ctx) {
		user := &models.User{}
		if err := result.Decode(&user); err != nil {
			return nil, errors.New("error in scanning user")
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *mongoDbRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"username": username}
	user := &models.User{}
	if err := m.db.GetUserCollection().FindOne(ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("data does not exists")
		}
		return nil, errors.New("error in fetching the user")
	}
	return user, nil
}

func (m *mongoDbRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"email": email}
	user := &models.User{}
	if err := m.db.GetUserCollection().FindOne(ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("data does not exists")
		}
		return nil, errors.New("error in fetching the user")
	}
	return user, nil
}

func (m *mongoDbRepo) UsernameExists(username string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"username": username}
	existingUser := &models.User{}
	if err := m.db.GetUserCollection().FindOne(ctx, query).Decode(&existingUser); err == nil {
		return errors.New("username exists")
	}
	return nil
}

func (m *mongoDbRepo) EmailExists(email string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"email": email}
	existingUser := &models.User{}
	if err := m.db.GetUserCollection().FindOne(ctx, query).Decode(&existingUser); err == nil {
		return errors.New("email exists")
	}
	return nil
}

func (m *mongoDbRepo) DeleteUser(username string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"username": username}
	res, err := m.db.GetUserCollection().DeleteOne(ctx, query)
	if err != nil {
		return errors.New("error in deleteing the user")
	}
	if res.DeletedCount == 0 {
		return errors.New("no user deleted")
	}
	return nil

}

func (m *mongoDbRepo) CreateUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	res, err := m.db.GetUserCollection().InsertOne(ctx, user)
	if err != nil {
		return nil, errors.New("errors in inserting new user")
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *mongoDbRepo) UpdateUser(username string, update *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"username": username}
	_, err := m.db.GetUserCollection().UpdateOne(ctx, query, update)
	if err != nil {
		return nil, errors.New("error in updating user")
	}
	user, err := m.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("errors in fetching the updated record")
	}
	return user, nil
}
