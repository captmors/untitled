package dockermocker

import (
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	user     = "postgres"
	password = "postgres"
	dbName   = "postgres"
	port     = "5432"
	dsn      = "host=%s port=%s user=%s sslmode=disable password=%s dbname=%s"
)

// DockerTest represents the test environment with Docker and PostgreSQL.
type DockerTest struct {
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
	Port     string
}

// NewDockerTest initializes a Docker environment for testing with PostgreSQL.
// It sets up the Docker pool and resource with the specified exposed port.
func NewDockerTest(exposedPort string, containerLifetime uint) *DockerTest {
	pool, resource := initTestDocker(exposedPort, containerLifetime)
	return &DockerTest{
		Pool:     pool,
		Resource: resource,
		Port:     exposedPort,
	}
}
// initTestDocker starts a PostgreSQL container for testing with the given port.
// It configures Docker options and sets an expiration time for the resource.
func initTestDocker(exposedPort string, containerLifetime uint) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + user,
		},
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {
				{HostIP: "0.0.0.0", HostPort: exposedPort},
			},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := resource.Expire(containerLifetime); err != nil {
		log.Fatalf("Could not set expiration for resource: %s", err)
	}

	log.Println("Trying to connect to the database...")
	
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db := connectToDatabase(resource, exposedPort)
		if db == nil {
			return fmt.Errorf("failed to connect to the database")
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return pool, resource
}

func (d *DockerTest) OpenDatabaseConnection() *gorm.DB {
	return connectToDatabase(d.Resource, d.Port)
}

// connectToDatabase establishes a connection to the PostgreSQL database
// with retry logic. It constructs the DSN and retries connection attempts if necessary.
func connectToDatabase(resource *dockertest.Resource, exposedPort string) *gorm.DB {
	host := resource.GetBoundIP(fmt.Sprintf("%s/tcp", port))
	dsn := fmt.Sprintf(dsn, host, exposedPort, user, password, dbName)

	var gdb *gorm.DB
	var err error
	retries := 5

	for retries > 0 {
		gdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		retries--
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Printf("Could not establish database connection: %s", err)
		return nil
	}

	return gdb
}

// Close purges the Docker resource after tests are completed.
// This ensures the container is removed and cleaned up properly.
func (d *DockerTest) Close() {
	if err := d.Pool.Purge(d.Resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
