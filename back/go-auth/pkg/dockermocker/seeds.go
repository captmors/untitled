package dockermocker

import (
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type Seeder interface {
	// Run generates and inserts n instances of the model into the database.
	Run(db *gorm.DB, n int) error
}

type UniformSeeder struct {
	modelPrototype interface{} // The prototype of the model that will be used for seeding
}

// NewUniformSeeder creates a new UniformSeeder for the given model.
func NewUniformSeeder(model interface{}) *UniformSeeder {
	return &UniformSeeder{
		modelPrototype: model,
	}
}

// Run generates and inserts n instances of the model into the database.
func (s *UniformSeeder) Run(db *gorm.DB, n int) error {
	modelType := reflect.TypeOf(s.modelPrototype)

	for i := 0; i < n; i++ {
		modelInstance := reflect.New(modelType.Elem()).Interface()

		// Populate the model instance with fake data
		if err := faker.FakeData(modelInstance); err != nil {
			return err
		}

		// Insert the populated model into the database
		if err := db.Create(modelInstance).Error; err != nil {
			return err
		}
	}
	return nil
}

// SeedDatabase executes the provided seeders, mocking the database
func SeedDatabase(db *gorm.DB, seeders []Seeder, n int) {
	for _, seed := range seeders {
		if err := seed.Run(db, n); err != nil {
			log.Fatalf("Error running seeder: %s", err)
		}
	}
}
