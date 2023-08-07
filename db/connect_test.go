package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/4epyx/todorpc/db"
	"github.com/stretchr/testify/suite"
)

type TestConnectToDb struct {
	suite.Suite
	correctDbUrl   string
	incorrectDbUrl string
}

func (t *TestConnectToDb) SetupTest() {
	var ok bool
	t.correctDbUrl, ok = os.LookupEnv("TEST_DB_URL")
	if !ok {
		t.T().Fatal("test db url not found in environment variable")
	}

	t.T().Log(t.correctDbUrl)

	t.incorrectDbUrl = "pg://fakeuser:incorrectpassword@incorrecthost:1111/db"
}

func (t *TestConnectToDb) TestCorrectUrl() {
	_, err := db.ConnectToDB(context.Background(), t.correctDbUrl)
	t.Nil(err)
}

func (t *TestConnectToDb) TestIncorrectUrl() {
	_, err := db.ConnectToDB(context.Background(), t.incorrectDbUrl)
	t.NotNil(err)
}

func TestConnectToDbSuite(t *testing.T) {
	suite.Run(t, new(TestConnectToDb))
}
