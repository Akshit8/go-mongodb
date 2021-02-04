package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
	// InsertingDocuments(ctx)

	fmt.Println("Listing all documents")
	listAllDocuments(ctx)

	fmt.Println("Getting single document")
	getSingleDocument(ctx)

	fmt.Println("Querying database")
	queryingDocuments(ctx)

	fmt.Println("Sorting documents in query")
	sortingDocumentsInQuery(ctx)

	fmt.Println("update one document in query")
	// updateOneDocument(ctx)

	fmt.Println("updated many documents in query")
	// updateManyDocuments(ctx)

	fmt.Println("replacing one documents in query")
	replaceOne(ctx)
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

func updateOneDocument(ctx context.Context) {
	id, err := primitive.ObjectIDFromHex("601ad97ed8880f320dd66dc8")
	if err != nil {
		log.Fatal("error converting string to id: ", err)
	}
	result, err := podcastCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"author", "updated Akshit Sadana"},
			}},
		},
	)
	if err != nil {
		log.Fatal("error updating the document:", err)
	}
	fmt.Printf("updated %d Documents\n", result.ModifiedCount)
}

func updateManyDocuments(ctx context.Context) {
	result, err := episodeCollection.UpdateMany(
		ctx,
		bson.M{"duration": 25},
		bson.D{
			{"$set", bson.D{{"duration", 34}}},
		},
	)
	if err != nil {
		log.Fatal("error updating many docs: ", err)
	}
	fmt.Printf("updated %d Documents\n", result.ModifiedCount)
}

func replaceOne(ctx context.Context) {
	result, err := podcastCollection.ReplaceOne(
		ctx,
		bson.M{"author": "updated Akshit Sadana"},
		bson.M{
			"title":  "The Akshit Sadana Show",
			"author": "Akshit Sadana",
		},
	)
	if err != nil {
		log.Fatal("error replacing doc: ", err)
	}
	fmt.Printf("Replaced %v Document\n", result.ModifiedCount)
}
