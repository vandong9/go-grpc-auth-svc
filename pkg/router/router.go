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
	// options   []httptransport.ServerOption
}

func NewHandler(endpoints endpoint.Endpoints, logger logruskit.FieldLogger) *Handler {
	return &Handler{
		endpoints: endpoints,
		logger:    logger,
		// options: []httptransport.ServerOption{
		// 	httptransport.ServerBefore(
		// 		httptransport.PopulateRequestContext,
		// 	),
		// },
	}
}

func (h *Handler) MakeHandlers(m middleware.Middleware) http.Handler {
	var r = mux.NewRouter()

	cp := r.PathPrefix("/authen").Subrouter()
	{
		cp.Methods(http.MethodPost).Path("/api/v1/login").Handler(h.Login())
	}
	return r

}

func (h *Handler) Login() http.Handler {
	return services.HttpSever{}
}
