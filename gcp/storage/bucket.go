package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/vickydk/utl/config"
	"github.com/vickydk/utl/log"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"path"
	"sync"
)

type Connections struct {
	Connection *storage.Client
	Context    context.Context
}

var (
	connection *Connections
	mutex      sync.Mutex
)

func GetConnection() *Connections {
	if connection == nil {
		mutex.Lock()
		defer mutex.Unlock()
		connection = newConnection()
	}

	return connection
}

func newConnection() *Connections {
	ctx := context.Background()
	client := new(storage.Client)
	client, err := storage.NewClient(ctx)
	if err != nil {
		filename := path.Join("./", config.Env.GcpCredentialFileName)
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(filename))
		if err != nil {
			log.Panic("got an error while connecting storage, ", zap.Error(err))
		}
	}

	bucket := client.Bucket(config.Env.GcpStorageBucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		log.Panic("got an error while connecting bucket: ", zap.Error(err))
	}

	return &Connections{Connection: client, Context: ctx}
}
