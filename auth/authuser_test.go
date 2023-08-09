package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/4epyx/todorpc/auth"
	"github.com/4epyx/todorpc/pb/authpb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
)

type TestAuthUser struct {
	suite.Suite
	ctx  context.Context
	auth *auth.AuthUser
}

func (t *TestAuthUser) SetupTest() {
	t.ctx = context.Background()
	url, ok := os.LookupEnv("AUTH_SERVER_URL")
	if !ok {
		t.T().Fatal("could not find auth server url environment")
	}

	var err error
	t.auth, err = auth.NewAuthUser(url)
	if err != nil {
		t.T().Fatal(err)
	}
}

func (t *TestAuthUser) TestAuth() {
	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTE2Njk0MjksInVzZXJfZW1haWwiOiJhYm9iYTJAZXhhbXBsZS5jb20iLCJ1c2VyX2lkIjoyfQ.3IXuD2tE4EEZsOIkD54czvCH5JkbgowGBaQ1cIYwRcE"
	ctx := t.getOutgiongContext(token)
	expected := &authpb.AuthUserData{
		Id:    2,
		Email: "aboba2@example.com",
	}

	data, err := t.auth.AuthorizeUser(ctx)
	t.Nil(err)
	t.Equal(expected.Id, data.Id)
	t.Equal(expected.Email, data.Email)
}

func (t *TestAuthUser) TestInvalidTokenAuht() {
	ctx := t.getOutgiongContext("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalidtoken")
	_, err := t.auth.AuthorizeUser(ctx)

	t.NotNil(err)
}

func (t *TestAuthUser) getOutgiongContext(token string) context.Context {
	md := metadata.New(map[string]string{"authorization": token})
	return metadata.NewOutgoingContext(t.ctx, md)
}

func TestAuthUserSuite(t *testing.T) {
	suite.Run(t, new(TestAuthUser))
}
