package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user
type User struct {
	ID             string         `json:"id" bson:"id"`
	Name           string         `json:"name" bson:"name"`
	Created        time.Time      `json:"created" bson:"created"`
	PieceCount     map[string]int `json:"-" bson:"-"`
	LifetimePieces int            `json:"lifetimePieces" bson:"lifetimePieces"`
	PasswordHash   string         `json:"passwordHash" bson:"passwordHash"`
}

// NewUser creates a new user
func NewUser(name string, password string) *User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	user := &User{
		ID:             uuid.New().String(),
		Name:           name,
		Created:        time.Now(),
		PieceCount:     make(map[string]int),
		LifetimePieces: 0,
		PasswordHash:   string(hash)}
	SaveUser(user)
	return user
}

// storeClient is the mongoClient
var userCollection *mongo.Collection

// InitStore inits the mongoDB storage
func InitStore() error {
	connString := os.Getenv("MONGODB_PUZZLE_CONN_STRING")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	storeClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return err
	}
	userCollection = storeClient.Database("puzzle").Collection("users")
	return nil
}

// SaveUser saves a user in the store
func SaveUser(u *User) error {
	_, err := userCollection.InsertOne(context.TODO(), u)
	return err
}

// UpdateUser updates a user in the store
func UpdateUser(u *User) error {
	_, err := userCollection.UpdateOne(context.TODO(), bson.M{"id": u.ID}, *u)
	return err
}

// GetUser retrieves a user from mongodb based on id
func GetUser(id string) (u User, err error) {
	err = userCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&u)
	fmt.Print(err)
	return
}

// AuthUser authenticates a user from mongodb based on username and password
func AuthUser(name string, password string) (*User, error) {
	cursor, err := userCollection.Find(context.TODO(), bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	for _, user := range users {
		pwError := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if pwError == nil {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("invalid password")
}
