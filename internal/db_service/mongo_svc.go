package db_service

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    "sync"
    "sync/atomic"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type DbService[DocType interface{}] interface {
	Connect(ctx context.Context) (*mongo.Client, error)
    CreateDocument(ctx context.Context, id string, document *DocType, collectionName string) error
    FindDocument(ctx context.Context, id string, collectionName string)  (*DocType, error)
	FindDocuments(ctx context.Context, collectionName string) ([]DocType, error)
	FindDocumentsByQuery(ctx context.Context, filter bson.D, collectionName string) ([]DocType, error)
    UpdateDocument(ctx context.Context, id string, document *DocType, collectionName string) error
    DeleteDocument(ctx context.Context, id string, collectionName string) error
    Disconnect(ctx context.Context) error
}

var ErrNotFound = fmt.Errorf("document not found")
var ErrConflict = fmt.Errorf("conflict: document already exists")

type MongoServiceConfig struct {
    ServerHost string
    ServerPort int
    UserName   string
    Password   string
    DbName     string
    Timeout    time.Duration
}

type MongoSvc[DocType interface{}] struct {
    MongoServiceConfig
    client     atomic.Pointer[mongo.Client]
    clientLock sync.Mutex
}

func NewMongoService[DocType interface{}](config MongoServiceConfig) DbService[DocType] {
	enviro := func(name string, defaultValue string) string {
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return defaultValue
	}

	svc := &MongoSvc[DocType]{}
	svc.MongoServiceConfig = config

	if svc.ServerHost == "" {
		svc.ServerHost = enviro("DEPARTMENT_API_MONGODB_HOST", "localhost")
	}

	if svc.ServerPort == 0 {
		port := enviro("DEPARTMENT_API_MONGODB_PORT", "27017")
		if port, err := strconv.Atoi(port); err == nil {
			svc.ServerPort = port
		} else {
			log.Printf("Invalid port value: %v", port)
			svc.ServerPort = 27017
		}
	}

	if svc.UserName == "" {
		svc.UserName = enviro("DEPARTMENT_API_MONGODB_USERNAME", "")
	}

	if svc.Password == "" {
		svc.Password = enviro("DEPARTMENT_API_MONGODB_PASSWORD", "")
	}

	if svc.DbName == "" {
		svc.DbName = enviro("DEPARTMENT_API_MONGODB_DATABASE", "cernica-department")
	}

	if svc.Timeout == 0 {
		seconds := enviro("DEPARTMENT_API_MONGODB_TIMEOUT_SECONDS", "10")
		if seconds, err := strconv.Atoi(seconds); err == nil {
			svc.Timeout = time.Duration(seconds) * time.Second
		} else {
			log.Printf("Invalid timeout value: %v", seconds)
			svc.Timeout = 10 * time.Second
		}
	}

	log.Printf(
		"MongoDB config: //%v@%v:%v/%v/%v",
		svc.UserName,
		svc.ServerHost,
		svc.ServerPort,
		svc.DbName,
		svc.Timeout.String(),
	)
	return svc
}

func (this *MongoSvc[DocType]) Connect(ctx context.Context) (*mongo.Client, error) {
    // optimistic check
    client := this.client.Load()
    if client != nil {
        return client, nil
    }

    this.clientLock.Lock()
    defer this.clientLock.Unlock()
    // pesimistic check
    client = this.client.Load()
    if client != nil {
        return client, nil
    }

    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()

    var uri = fmt.Sprintf("mongodb://%v:%v", this.ServerHost, this.ServerPort)
    log.Printf("Using URI: " + uri)

    if len(this.UserName) != 0 {
        uri = fmt.Sprintf("mongodb://%v:%v@%v:%v", this.UserName, this.Password, this.ServerHost, this.ServerPort)
    }

    if client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetConnectTimeout(10*time.Second)); err != nil {
        return nil, err
    } else {
        this.client.Store(client)
        return client, nil
    }
}

func (this *MongoSvc[DocType]) Disconnect(ctx context.Context) error {
    client := this.client.Load()

    if client != nil {
        this.clientLock.Lock()
        defer this.clientLock.Unlock()

        client = this.client.Load()
        defer this.client.Store(nil)
        if client != nil {
            if err := client.Disconnect(ctx); err != nil {
                return err
            }
        }
    }
    return nil
}

func (this *MongoSvc[DocType]) CreateDocument(ctx context.Context, id string, document *DocType, collectionName string) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.Connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil: // no error means there is conflicting document
        return ErrConflict
    case mongo.ErrNoDocuments:
        // do nothing, this is expected
    default: // other errors - return them
        return result.Err()
    }

    _, err = collection.InsertOne(ctx, document)
    return err
}

func (this *MongoSvc[DocType]) FindDocument(ctx context.Context, id string, collectionName string) (*DocType, error) {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.Connect(ctx)
    if err != nil {
        return nil, err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return nil, ErrNotFound
    default: // other errors - return them
        return nil, result.Err()
    }
    var document *DocType
    if err := result.Decode(&document); err != nil {
        return nil, err
    }
    return document, nil
}

func (this *MongoSvc[DocType]) FindDocuments(ctx context.Context, collectionName string) ([]DocType, error) {
    ctxWithTimeout, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()

    client, err := this.Connect(ctxWithTimeout)
    if err != nil {
        return nil, err
    }

    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)

    // Define empty filter to get all documents
    filter := bson.D{{}}

    // Define options to set any options needed for the query
    findOptions := options.Find()

    // Perform the find operation
    cursor, err := collection.Find(ctxWithTimeout, filter, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctxWithTimeout)

    // Iterate over the cursor and decode documents into a slice
    var documents []DocType
    for cursor.Next(ctxWithTimeout) {
        var document DocType
        if err := cursor.Decode(&document); err != nil {
            return nil, err
        }
        documents = append(documents, document)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return documents, nil
}

func (this *MongoSvc[DocType]) FindDocumentsByQuery(ctx context.Context, filter bson.D, collectionName string) ([]DocType, error) {
    ctxWithTimeout, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()

    client, err := this.Connect(ctxWithTimeout)
    if err != nil {
        return nil, err
    }

    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)

    // Define options to set any options needed for the query
    findOptions := options.Find()

    // Perform the find operation
    cursor, err := collection.Find(ctxWithTimeout, filter, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctxWithTimeout)

    // Iterate over the cursor and decode documents into a slice
    var documents []DocType
    for cursor.Next(ctxWithTimeout) {
        var document DocType
        if err := cursor.Decode(&document); err != nil {
            return nil, err
        }
        documents = append(documents, document)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return documents, nil
}


func (this *MongoSvc[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType, collectionName string) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.Connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return ErrNotFound
    default: // other errors - return them
        return result.Err()
    }
    _, err = collection.ReplaceOne(ctx, bson.D{{Key: "id", Value: id}}, document)
    return err
}

func (this *MongoSvc[DocType]) DeleteDocument(ctx context.Context, id string, collectionName string) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.Connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(collectionName)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return ErrNotFound
    default: // other errors - return them
        return result.Err()
    }
    _, err = collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}