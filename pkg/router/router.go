package router

import (
	"net/http"

	"github.com/gorilla/mux"
	logruskit "github.com/sirupsen/logrus"
	"github.com/vandong9/go-grpc-auth-svc/pkg/endpoint"
	"github.com/vandong9/go-grpc-auth-svc/pkg/middleware"
	"github.com/vandong9/go-grpc-auth-svc/pkg/services"
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

	return services.HttpSever{
		h.endpoints.Login,h.options...,
	}
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