package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/vickydk/utl/config"
	"github.com/vickydk/utl/log"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"path"
	"sync"
)

type Connections struct {
	Connection *pubsub.Client
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
	client := new(pubsub.Client)
	client, err := pubsub.NewClient(ctx, config.Env.GcpProjectId)
	if err != nil {
		filename := path.Join("./", config.Env.GcpCredentialFileName)
		client, err = pubsub.NewClient(ctx, config.Env.GcpProjectId, option.WithCredentialsFile(filename))
		if err != nil {
			log.Panic("got an error while connecting pubsub, ", zap.Error(err))
		}
	}

	return &Connections{Connection: client, Context: ctx}
}
