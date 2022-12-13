package endpoint

import (
	"context"

	base_endpoint "github.com/vandong9/go-grpc-auth-svc/pkg/base_enpoint"
	"github.com/vandong9/go-grpc-auth-svc/pkg/models"
	"github.com/vandong9/go-grpc-auth-svc/pkg/services"
)

type Endpoints struct {
	Login base_endpoint.Endpoint
}

func MakeEndpoints(s services.HttpSever) Endpoints {
	return Endpoints{
		Login: makeLogin(s),
	}
}

func makeLogin(c services.HttpSever) base_endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.LoginRequest)
		return c.HandleLogin(ctx, req)
	}
}
