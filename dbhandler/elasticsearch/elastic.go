package elastic

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/vickydk/utl/config"
	"github.com/vickydk/utl/log"
	"go.uber.org/zap"
	logs "log"
	"net/http"
	"os"
	"sync"
	"time"
)

var once sync.Once
var instance *elastic.Client

func GetClient() *elastic.Client {
	once.Do(func() {
		instance = new()
	})
	return instance
}

func new() *elastic.Client {
	// Pass the connection settings as variables instead of hardcoding
	var sec time.Duration = 5

	// Convert port integer to a string
	timeOut := sec * time.Second

	// Concatenate domain string
	setURL := config.Env.ElasticUrl

	// Instantiate a client instance of the elastic library
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetErrorLog(logs.New(os.Stderr, "ELASTIC ", logs.LstdFlags)),
		elastic.SetURL(setURL),
		elastic.SetHealthcheckInterval(timeOut), // quit trying after 5 seconds
	)

	if err != nil {
		log.Panic("elastic.NewClient() ERROR", zap.Error(err))
	}

	ctx := context.Background()
	// Ping the Elasticsearch server to get e.g. the version number
	_, code, err := client.Ping(setURL).Do(ctx)
	if err != nil {
		log.Panic("elastic.Ping() ERROR", zap.Error(err))
	}
	if code != http.StatusOK {
		log.Panic("elastic.Ping() resp not 200", zap.Int("code", code))
	}

	return client
}
