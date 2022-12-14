package router

import (
	"context"
	"github.com/gorilla/mux"
	logruskit "github.com/sirupsen/logrus"
	"github.com/vandong9/go-grpc-auth-svc/pkg/endpoint"
	"github.com/vandong9/go-grpc-auth-svc/pkg/middleware"
	"github.com/vandong9/go-grpc-auth-svc/pkg/models"
	"github.com/vandong9/go-grpc-auth-svc/pkg/services"
	"net/http"
)

type Handler struct {
	endpoints endpoint.Endpoints
	logger    logruskit.FieldLogger
	options   []services.ServerOption
}

func NewHandler(endpoints endpoint.Endpoints, logger logruskit.FieldLogger) *Handler {
	return &Handler{
		endpoints: endpoints,
		logger:    logger,
		options: []services.ServerOption{
			services.ServerBefore(
				services.PopulateRequestContext,
			),
		},
	}
}

func (h *Handler) MakeHandlers(m middleware.Middleware) http.Handler {
	var r = mux.NewRouter()

	cp := r.PathPrefix("/authen").Subrouter()
	{
		cp.Methods(http.MethodGet).Path("/api/v1/login").Handler(h.Login())
	}
	return r

}

func (h *Handler) Login() http.Handler {
	return services.NewServer(h.endpoints.Login, decodeLoginRequest, services.EncodeJSONResponse, h.options...)
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	email := vars["email"]
	//if ok == "" {
	//	return nil, errorkit.BadRequest("The 'cardseno' path param is required.")
	//}
	password := vars["password"]

	req := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	return req, nil
}

// return httptransport.NewServer(
// 	h.endpoints.GetDummyHomeInfo,
// 	decodeGetDummyHomeInfo,
// 	httptransport.EncodeJSONResponse,
// 	h.options...,
// )

// func decodeGetDummyHomeInfo(_ context.Context, r *http.Request) (request interface{}, err error) {
// 	req := model.GetDummyHomeInfoRequest{
// 		ID: r.URL.Query().Get("id"),
// 	}
// 	if err := validator.New().Struct(req); err != nil {
// 		return nil, errorkit.BadRequest(err)
// 	}

// 	return &req, nil
// }
