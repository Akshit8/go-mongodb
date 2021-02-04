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

var (
	podcastCollection *mongo.Collection
	episodeCollection *mongo.Collection
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

	fmt.Println("Listing all databases")
	listDatabases(ctx, client)

	// connecting to db quickstart
	quickStartDatabase := client.Database(dbname)
	podcastCollection = quickStartDatabase.Collection(podcastCollectionName)
	episodeCollection = quickStartDatabase.Collection(episodeCollectionName)

	fmt.Println("Inserting documents")

	fmt.Println("Listing all documents")
	listAllDocuments(ctx)

	fmt.Println("Getting single document")
	getSingleDocument(ctx)

	fmt.Println("Querying database")
	queryingDocuments(ctx)

	fmt.Println("Sorting Documents In Query")
	sortingDocumentsInQuery(ctx)
}

func listDatabases(ctx context.Context, client *mongo.Client) {
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("error retrieving databases list")
	}
	fmt.Println(databases)
}

func InsertingDocuments(ctx context.Context) {
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

func listAllDocuments(ctx context.Context) {
	cursor, err := episodeCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("error listing all documents: ", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var episode bson.M
		err = cursor.Decode(&episode)
		if err != nil {
			log.Fatal("error decoding cursor: ", err)
		}
		fmt.Println(episode)
	}
}

func getSingleDocument(ctx context.Context) {
	var podcast bson.M
	err := podcastCollection.FindOne(ctx, bson.M{}).Decode(&podcast)
	if err != nil {
		log.Fatal("error decoding single podcast: ", err)
	}
	fmt.Println(podcast)
}

func queryingDocuments(ctx context.Context) {
	filterCursor, err := episodeCollection.Find(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal("err querying db: ", err)
	}
	defer filterCursor.Close(ctx)
	for filterCursor.Next(ctx) {
		var episode bson.M
		err = filterCursor.Decode(&episode)
		if err != nil {
			log.Fatal("error decoding filtered episode: ", err)
		}
		fmt.Println(episode)
	}
}

func sortingDocumentsInQuery(ctx context.Context) {
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", -1}})
	sortCursor, err := episodeCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}})
	if err != nil {
		log.Fatal("err sorting query db: ", err)
	}
	defer sortCursor.Close(ctx)
	for sortCursor.Next(ctx) {
		var episode bson.M
		err = sortCursor.Decode(&episode)
		if err != nil {
			log.Fatal("error decoding filtered episode: ", err)
		}
		fmt.Println(episode)
	}
}
