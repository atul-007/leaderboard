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

	// Default scope to globally if not provided
	if scope == "" {
		scope = "globally"
	}

	// Filter scores based on scope
	filter := bson.M{}

	switch scope {
	case "globally":
		// No additional filter needed for global scope
	case "country":
		userCountry, err := GetUserCountry(userName)
		if err != nil {
			return 0, err
		}
		filter["country"] = userCountry
	case "state":
		userState, err := GetUserState(userName)
		if err != nil {
			return 0, err
		}
		filter["state"] = userState
	case "city":
		userCity, err := GetUserCity(userName)
		if err != nil {
			return 0, err
		}
		filter["city"] = userCity
	default:
		return 0, fmt.Errorf("invalid scope: %s", scope)
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

	// Find rank of the user among users with the same country or city
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
func ListTopN(n int, scope string, scopeName string) ([]models.Score, error) {
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
		filter["country"] = scopeName
	case "state":
		filter["state"] = scopeName
	case "city":
		filter["city"] = scopeName
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

// GetUserByUsername retrieves the user document based on the provided username
func GetUserByUsername(userName string) (*models.Score, error) {
	var user models.Score
	filter := bson.M{"username": userName}

	err := scoreCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, return nil without error
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetUserCity retrieves the city of the specified user from the database
func GetUserCity(userName string) (string, error) {
	var user models.Score
	filter := bson.M{"username": userName}

	err := scoreCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, return appropriate error message
			return "", fmt.Errorf("user %s not found", userName)
		}
		return "", err
	}

	return user.City, nil
}

// GetUserCountry retrieves the country of the specified user from the database
func GetUserCountry(userName string) (string, error) {
	var user models.Score
	filter := bson.M{"username": userName}

	err := scoreCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, return appropriate error message
			return "", fmt.Errorf("user %s not found", userName)
		}
		return "", err
	}

	return user.Country, nil
}
func GetUserState(userName string) (string, error) {
	var user models.Score
	filter := bson.M{"username": userName}

	err := scoreCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, return appropriate error message
			return "", fmt.Errorf("user %s not found", userName)
		}
		return "", err
	}

	return user.State, nil
}

// UpdateScore updates the score of the specified user in the database
func UpdateScore(score *models.Score) error {
	filter := bson.M{"username": score.UserName}
	update := bson.M{"$set": bson.M{"score": score.Score}}

	_, err := scoreCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
