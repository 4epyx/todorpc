package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/4epyx/todorpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthUserInterceptor struct {
	authUser *auth.AuthUser
}

func NewAuthUserInterceptor(authUser *auth.AuthUser) *AuthUserInterceptor {
	return &AuthUserInterceptor{authUser: authUser}
}

func (a *AuthUserInterceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	header, err := a.getHeaderValue(ctx, "authorization")
	if err != nil {
		return nil, err
	}

	data, err := a.authUser.AuthorizeUser(metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"authorization": header})))
	if err != nil {
		return nil, err
	}

	return handler(metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"user_id": strconv.FormatInt(data.Id, 10)})), req)
}

func (*AuthUserInterceptor) parseMetadataFromCtx(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error while getting metadata from context")
	}
	return md, nil
}

func (a *AuthUserInterceptor) getHeaderValue(ctx context.Context, headerName string) (string, error) {
	md, err := a.parseMetadataFromCtx(ctx)
	if err != nil {
		return "", err
	}
	fmt.Println(md)
	header := md[headerName]
	if len(header) == 0 {
		return "", fmt.Errorf("header %s not found", headerName)
	}

	return header[0], nil
}
