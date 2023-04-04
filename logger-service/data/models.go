package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo


	return Models{
		LogEntry: LogEntry{},
	}
}


type Models struct{
	LogEntry LogEntry
}

type LogEntry struct{
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	col := client.Database("logs").Collection("logs")
	_, err := col.InsertOne(context.TODO(), LogEntry{
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs ", err )
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel() 

	col := client.Database("logs").Collection("logs")
	options := options.Find()
	options.SetSort(bson.D{{"created_at", -1}})

	cursor, err := col.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Println("finding all docs error ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil{
			log.Println("error decoding log into slice ", err)
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}


func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel() 

	col := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return nil, err
	}

	var entry LogEntry
	err = col.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil{
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel() 

	col := client.Database("logs").Collection("logs")
	if err := col.Drop(ctx); err != nil {
		return err
	}

	return nil
}



func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel() 

	col := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := col.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}



	return result, nil
}
