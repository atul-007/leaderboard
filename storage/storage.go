package storage

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/atul-007/leaderboard/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	scoreCollection *mongo.Collection
	mu              sync.Mutex // Mutex for locking during database operations
)

// InitMongoDB initializes MongoDB connection
func InitMongoDB() error {
	//clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	clientOptions := options.Client().ApplyURI("mongodb+srv://atulranjan789:atul1234@cluster0.xr7e6vt.mongodb.net/")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	scoreCollection = client.Database("gaming").Collection("scores")
	log.Println("Connected to MongoDB")
	return nil
}

// SaveScore saves score to MongoDB
func SaveScore(score *models.Score) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := scoreCollection.InsertOne(context.Background(), score)
	return err
}

// GetRank fetches the rank of a user based on userName and scope
func GetRank(userName, scope string) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	var scores []models.Score

	// Filter scores based on scope
	filter := bson.M{}

	if scope != "globally" {
		if scope == "country" {
			filter["country"] = userName
		} else if scope == "state" {
			filter["state"] = userName
		}
	}

	cursor, err := scoreCollection.Find(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())

	// Decode scores from cursor
	if err := cursor.All(context.Background(), &scores); err != nil {
		return 0, err
	}

	// Sort scores by score value in descending order
	sortScores(&scores)

	// Find rank of the user
	rank := 0
	for i, s := range scores {
		if s.UserName == userName {
			rank = i + 1
			break
		}
	}

	return rank, nil
}

// ListTopN lists the top N ranks based on the scope
func ListTopN(n int, scope string) ([]models.Score, error) {
	mu.Lock()
	defer mu.Unlock()

	var scores []models.Score

	// Define the filter based on scope
	filter := bson.M{}
	switch scope {
	case "":
		// Default to globally if scope is not provided
		// No additional filter needed for global scope
	case "globally":
		// No additional filter needed for global scope
	case "country":
		filter["country"] = bson.M{"$exists": true}
	case "state":
		filter["state"] = bson.M{"$exists": true}
	case "India":
		filter["country"] = "India"
	default:
		return nil, fmt.Errorf("invalid scope: %s", scope)
	}

	// Find top N scores based on filter
	limit := int64(n)
	options := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(limit)
	cursor, err := scoreCollection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode scores from cursor
	if err := cursor.All(context.Background(), &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

// sortScores sorts scores by score value in descending order
func sortScores(scores *[]models.Score) {
	sort.Slice(*scores, func(i, j int) bool {
		return (*scores)[i].Score > (*scores)[j].Score
	})
}
