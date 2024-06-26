package service

import (
	"context"
	"fmt"
	
	"github.com/superwhys/venkit/v2/service/example/grpc/examplepb"
	"google.golang.org/grpc/metadata"
)

type ExampleService struct {
	examplepb.UnimplementedExampleHelloServiceServer
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (es *ExampleService) SayHello(ctx context.Context, in *examplepb.HelloRequest) (*examplepb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, ok)
	
	return &examplepb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}
