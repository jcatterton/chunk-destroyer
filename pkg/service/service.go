package service

import (
	"chunk-destroyer/config"
	"chunk-destroyer/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Config   *config.Config
	DBClient *mongo.Client
}

func InitializeService(ctx context.Context, conf *config.Config) (*Service, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoURI))
	if err != nil {
		return nil, err
	}

	return &Service{
		Config:   conf,
		DBClient: client,
	}, nil
}

func (s *Service) Run(ctx context.Context) (int, error) {
	fsFiles, err := s.GetFiles(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error retrieving files from database")
		return 0, err
	}

	deletedChunks, err := s.DeleteOrphanChunks(ctx, fsFiles)
	if err != nil {
		return 0, err
	}

	return deletedChunks, nil
}

func (s *Service) GetFiles(ctx context.Context) ([]primitive.ObjectID, error) {
	cursor, err := s.getFsFilesCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var results []models.FsFile
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, 0)
	for _, f := range results {
		ids = append(ids, f.ID)
	}

	return ids, nil
}

func (s *Service) DeleteOrphanChunks(ctx context.Context, fileIDs []primitive.ObjectID) (int, error) {
	query := bson.M{"files_id": bson.M{"$nin": fileIDs}}

	result, err := s.getFsChunksCollection().DeleteMany(ctx, query)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (s *Service) getFsFilesCollection() *mongo.Collection {
	return s.DBClient.Database(s.Config.Database).Collection(s.Config.FsFilesCollection)
}

func (s *Service) getFsChunksCollection() *mongo.Collection {
	return s.DBClient.Database(s.Config.Database).Collection(s.Config.FsChunksCollection)
}
