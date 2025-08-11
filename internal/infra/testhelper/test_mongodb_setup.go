package testhelper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestDatabase struct {
	DbInstance *mongo.Database
	DbAddress  string
	container  testcontainers.Container
}

// SetupTestDatabase sets up a MongoDB TestContainer for testing purposes.
func SetupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("Setting up MongoDB TestContainer...")
	container, dbInstance, dbAddr, err := createMongoContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func (tdb *TestDatabase) CleanUp() error {
	collections, err := tdb.DbInstance.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	for _, name := range collections {
		if err := tdb.DbInstance.Collection(name).Drop(context.Background()); err != nil {
			return fmt.Errorf("drop collection %s: %w", name, err)
		}
	}
	log.Println("✅ Database cleaned up")
	return nil
}

func createMongoContainer(ctx context.Context) (testcontainers.Container, *mongo.Database, string, error) {
	var env = map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": "test",
		"MONGO_INITDB_ROOT_PASSWORD": "tester",
		"MONGO_INITDB_DATABASE":      "wordrop_test",
	}
	port := "27017/tcp"
	log.Printf("Using port %s for MongoDB", port)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}

	log.Printf("Starting MongoDB TestContainer with request: %+v", req)
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to map port: %v", err)
	}

	log.Printf("MongoDB TestContainer started with mapped port: %s", mappedPort.Port())
	uri := fmt.Sprintf("mongodb://test:tester@localhost:%s", mappedPort.Port())

	log.Printf("Connecting to MongoDB at %s", uri)
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return container, nil, uri, fmt.Errorf("connect to mongo: %w", err)
	}

	log.Println("✅ Connected to MongoDB")

	db := client.Database("wordrop_test")
	log.Println("✅ Mongo TestContainer running at", uri)
	return container, db, uri, nil
}
