package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbname                = "quickstart"
	podcastCollectionName = "podcasts"
	episodeCollectionName = "episodes"
)

func main() {
	mongoURI := os.Getenv("MONGODB_CONNECTION_STRING")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("error creating mongo client: ", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}
	defer client.Disconnect(ctx)

	// Listing all databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("error retrieving databases list")
	}
	fmt.Println(databases)

	// connecting to db quickstart
	quickStartDatabse := client.Database(dbname)
	podcastCollection := quickStartDatabse.Collection(podcastCollectionName)
	episodeCollection := quickStartDatabse.Collection(episodeCollectionName)

	podcastResult, err := podcastCollection.InsertOne(ctx, bson.D{
		{"title", "The Developer Podcast"},
		{"author", "Akshit Sadana"},
		{"tags", bson.A{"development", "programming", "coding"}},
	})
	if err != nil {
		log.Fatal("error inserting document: ", err)
	}

	episodeResult, err := episodeCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with React."},
			{"duration", 32},
		},
	})
	if err != nil {
		log.Fatal("error inserting multiple documents: ", err)
	}

	fmt.Printf("Inserted %d documents to episode collection\n", len(episodeResult.InsertedIDs))
}
