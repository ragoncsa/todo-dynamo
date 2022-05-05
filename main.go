package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ragoncsa/todo/config"
	"github.com/ragoncsa/todo/gorm"
	"github.com/ragoncsa/todo/http"
	"github.com/ragoncsa/todo/mock"

	awsconf "github.com/aws/aws-sdk-go-v2/config"

	"github.com/spf13/viper"

	docs "github.com/ragoncsa/todo/docs"
)

// @title        Tasks service
// @version      1.0
// @description  This is a sample server that manages tasks.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// Added workaround due to issues with environment variables in Viper
// https://github.com/spf13/viper/issues/761
func overrideUsingEnvVars(config *config.Config) {
	if host, present := os.LookupEnv("DB_HOST"); present {
		config.Database.Host = host
	}
}

func loadConfig() *config.Config {
	viper.SetConfigName("local-env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var conf config.Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
	overrideUsingEnvVars(&conf)
	return &conf
}

func main() {
	conf := loadConfig()

	docs.SwaggerInfo.BasePath = "/"

	server := http.InitServer(conf)
	sdkConfig, err := awsconf.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	tsDB := &gorm.TaskService{DynamoDbClient: dynamodb.NewFromConfig(sdkConfig)}
	tsHTTP := http.TaskService{
		Service: tsDB,
		// Removed authorization piece, since it's not relevant for this demo.
		// AuthzClient: authz.New(conf),
		AuthzClient: &mock.AlwaysAllow{},
	}
	server.RegisterRoutes(&tsHTTP)
	server.Start()
}
