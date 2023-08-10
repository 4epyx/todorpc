package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/4epyx/todorpc/auth"
	"github.com/4epyx/todorpc/db"
	"github.com/4epyx/todorpc/interceptor"
	"github.com/4epyx/todorpc/pb"
	"github.com/4epyx/todorpc/repository"
	pgxrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/service"
	"google.golang.org/grpc"
)

func main() {
	host, port := getHostAndPort()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		panic(err)
	}

	authInterceptor, err := setupAuthInterceptor()
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.AuthInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)

	repo, err := getPgxRepo()
	if err != nil {
		panic(err)
	}

	registerService(grpcServer, repo)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func getHostAndPort() (string /*host*/, string /*port*/) {
	host, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		host = "localhost"
	}

	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		port = "8080"
	}
	return host, port
}

func getPgxRepo() (*pgxrepo.PgxTaskRepository, error) {
	dbUrl, err := getDbUrl()
	if err != nil {
		return nil, err
	}

	conn, err := db.ConnectToDB(context.TODO(), dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.MigrateTaskTable(context.TODO(), conn); err != nil {
		return nil, err
	}

	return pgxrepo.NewTaskRepository(conn), nil
}

func getDbUrl() (string, error) {
	dbUrl, ok := os.LookupEnv("DB_URL")
	if !ok {
		return "", errors.New("DB_URL not found in environment variables")
	}

	return dbUrl, nil
}

func setupAuthInterceptor() (*interceptor.AuthUserInterceptor, error) {
	url, ok := os.LookupEnv("AUTH_SERVER_URL")
	if !ok {
		return nil, errors.New("could not find auth server url environment")
	}
	authUser, err := auth.NewAuthUser(url)
	if err != nil {
		return nil, err
	}

	return interceptor.NewAuthUserInterceptor(authUser), nil
}

func registerService(grpcServer *grpc.Server, repo repository.TaskRepository) {
	pb.RegisterTaskServiceServer(grpcServer, service.NewTaskService(repo))
}
