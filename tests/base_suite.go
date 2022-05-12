//go:build tests
// +build tests

package tests

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/khaiql/dbcleaner"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/services/session"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/storage"
	"github.com/stretchr/testify/suite"
)

func init() {
	fmt.Println("this is test build...")
}

type BaseTestSuite struct {
	suite.Suite
	dbcleaner.DbCleaner
}

func (suite *BaseTestSuite) SetupSuite() {
	suite.DbCleaner = dbcleaner.New()
	// TODO the env should read through api
	duration, _ := time.ParseDuration("24h")
	env.Envs = &env.Env{
		DBName:           "mises_unit_test",
		DBUser:           "",
		DBPass:           "",
		MongoURI:         "mongodb://localhost:27017",
		DebugMisesPrefix: "1001",
		TokenDuration:    duration,
	}
	db.SetupMongo(context.Background())
	models.EnsureIndex()
	session.SetupMisesClient()
	storage.SetupImageStorage("127.0.0.1", "xxx", "xx")

}

func (suite *BaseTestSuite) TearDownSuite() {
}

func (suite *BaseTestSuite) Clean(collections []string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
	for _, collection := range collections {
		_ = db.DB().Collection(collection).Drop(context.Background())
	}
	models.EnsureIndex()
}

func (suite *BaseTestSuite) Acquire(collections ...string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
}

func (suite *BaseTestSuite) CreateTestUsers(count int) {
	for i := 0; i < count; i++ {
		u := &models.User{
			UID:      uint64(i + 1),
			Username: fmt.Sprintf("user%d", i),
			Misesid:  fmt.Sprintf("mises%d", i),
		}
		db.ODM(context.Background()).Create(u)

	}
}
