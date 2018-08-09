package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"web-service-users/src/user/delivery/gin_http"
	"web-service-users/src/user/repository"
	"web-service-users/src/user/usercase"

	"web-service-users/src/user/searchdb"

	"web-service-users/src/user/webservice"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/vanng822/go-solr/solr"
)

func main() {

	modeType := os.Getenv("MODE")

	isProduction := false
	if modeType == "PRODUCTION" {
		isProduction = true
	}

	mongodbUser := os.Getenv("MONGODB_USER")
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")
	mongodbHost := os.Getenv("MONGODB_HOST")
	mongodbPort := os.Getenv("MONGODB_PORT")
	mongodbDBName := os.Getenv("MONGODB_DB_NAME")
	webAddr := os.Getenv("WEB_ADDRESS")
	solrHost := os.Getenv("SOLR_HOST_URL")
	solrCore := os.Getenv("SOLR_CORE")
	kudosHost := os.Getenv("KUDOS_WS_HOST")
	kudosPort := os.Getenv("KUDOS_WS_PORT")

	configureLog(isProduction)

	db := connectDataBase(mongodbUser, mongodbPassword, mongodbHost, mongodbPort, mongodbDBName)
	solrConnection := createSolrConnection(solrHost, solrCore)
	solrService := searchdb.NewSolrSearchDB(solrConnection)
	kudosWs := webservice.NewKudosWebService(kudosHost, kudosPort)
	userRepo := repository.NewMongodbUserRepository(db)
	userUC := usercase.NewUserUserCase(userRepo, kudosWs)
	engine := createGinHTTPDelivery()
	gin_http.NewUserGinHTTPHandler(engine, userUC, solrService)
	engine.Run(webAddr)

}

func connectDataBase(user string, password string, host string, port string, dataBaseName string) *mongo.Database {

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, dataBaseName)
	client, err := mongo.NewClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dataBaseName)
	return db

}

func createGinHTTPDelivery() *gin.Engine {

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true
	server.Use(cors.New(config))
	return server

}

func configureLog(isProduction bool) {

	if isProduction {
		file, err := os.OpenFile("/var/log/web-service-users/users.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}

}

func createSolrConnection(solrHost string, solrCore string) *solr.SolrInterface {
	solrInterface, err := solr.NewSolrInterface(solrHost, solrCore)
	if err != nil {
		log.Fatal(err)
	}
	return solrInterface
}
