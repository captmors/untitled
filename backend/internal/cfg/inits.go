package cfg

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"untitled/internal/interfaces"
	ms "untitled/internal/musicstorage"
	"untitled/internal/users"
	"untitled/internal/users/mdl"
	tus "untitled/internal/utils/tusuploader"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	initLogging()
}

func Init() *gin.Engine {
	r := gin.Default()

	pgDB := initPostgres()
    mongoDB := initMongoDB()
    esClient := initElasticSearch()

	initApps(r, pgDB, mongoDB, esClient)

	return r
}

var (
	apps []interfaces.App
)

func initPostgres() *gorm.DB {
    PostgresUser := os.Getenv("POSTGRES_USER")
    PostgresPassword := os.Getenv("POSTGRES_PASSWORD")
    PostgresDB := os.Getenv("POSTGRES_DB")

    dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
        PostgresUser, PostgresPassword, PostgresDB)

    DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Error connecting to database:", err)
		return nil
    }

    if err := DB.AutoMigrate(&mdl.User{}); err != nil {
        log.Fatal("Error migrating database:", err)
    }

	log.Info("Connected to PostgreSQL!")

    return DB
}


func initApps(r *gin.Engine, pgDB *gorm.DB, mongoDB *mongo.Client, esClient *elasticsearch.Client) {
    // Users (auth)
    usersApp := users.NewUserApp(pgDB, []byte(JwtSecret)) 
    apps = append(apps, usersApp)

    // MusicStorage (Mongo и ElasticSearch)
    tusHandler := tus.InitTusUploader(r, tus.TusHandlerCfg{
        MaxFileSize: MaxUploadFileSize,
        UploadDir:   UploadDir,
    })

    musicStorageApp := ms.NewMusicStorageApp(mongoDB, pgDB, esClient, tusHandler) 
    apps = append(apps, musicStorageApp)

    for _, app := range apps {
        app.Init(r)
    }
}

// ENV:
// - LOG_TO_FILE: bool
func initLogging() *os.File {
	logDir := filepath.Join(LogDir)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	var file *os.File
	var err error

	if LogToFile {
		logFile := filepath.Join(logDir, LogDefaultFile)
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		multiWriter := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(os.Stdout)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line), ""
		},
	})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)

	return file
}

func initMongoDB() *mongo.Client {
    MongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
    MongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

    uri := fmt.Sprintf("mongodb://%s:%s@localhost:27017", MongoUser, MongoPassword)
    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Info("Connected to MongoDB!")

    return client
}

func initElasticSearch() *elasticsearch.Client {
    cfg := elasticsearch.Config{
        Username: ElasticUser,
        Password: ElasticPassword,
    }

    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
        log.Fatalf("Error creating the ElasticSearch client: %s", err)
    }

    res, err := es.Info()
    if err != nil {
        log.Fatalf("Error getting ElasticSearch info: %s", err)
    }
    defer res.Body.Close()

    log.Info("Connected to ElasticSearch!")
    return es
}
