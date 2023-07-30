package dbrepository

import (
	"context"
	"errors"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mongoDbRepo) InsertFileData(file *models.File) (*models.File, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	res, err := m.db.GetFileCollection().InsertOne(ctx, file)
	if err != nil {
		return nil, err
	}
	file.ID = res.InsertedID.(primitive.ObjectID)
	return file, nil
}

func (m *mongoDbRepo) GetFileByFileName(filename string) (*models.File, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	filter := bson.M{"name": filename}
	file := &models.File{}
	if err := m.db.GetFileCollection().FindOne(ctx, filter).Decode(&file); err != nil {
		return nil, err
	}
	return file, nil
}

func (m *mongoDbRepo) AllFilesByUser(username string, limit, page int) ([]*models.File, error) {
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
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	filter := bson.M{
		"uploader.username": username,
	}
	res, err := m.db.GetFileCollection().Find(ctx, filter, &opt)
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)
	files := []*models.File{}
	for res.Next(ctx) {
		file := &models.File{}
		if err := res.Decode(&file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (m *mongoDbRepo) FileNameExists(fileName string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()
	existingFile := &models.File{}
	err := m.db.GetFileCollection().FindOne(ctx, bson.M{"name": fileName}).Decode(&existingFile)
	if err == nil {
		return errors.New("username already exists")
	}
	return nil
}
