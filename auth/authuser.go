package auth

import (
	"context"

	"github.com/4epyx/todorpc/pb/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthUser struct {
	client authpb.AuthorizationServiceClient
}

func NewAuthUser(serverUrl string) (*AuthUser, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(serverUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := authpb.NewAuthorizationServiceClient(conn)
	return &AuthUser{client}, nil
}

func (a *AuthUser) AuthorizeUser(ctx context.Context) (*authpb.AuthUserData, error) {
	authdata, err := a.client.AuthorizeUser(ctx, &authpb.Empty{})
	if err != nil {
		return nil, err
	}

	return authdata, nil
}
