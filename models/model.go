package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FsFile struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type FsChunk struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type Log struct {
	AppName   string
	Event     string
	Message   string
	IsError   bool
	TimeStamp time.Time
}
