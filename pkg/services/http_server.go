package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vandong9/go-grpc-auth-svc/pkg/base_enpoint"
	"github.com/vandong9/go-grpc-auth-svc/pkg/models"
)

type contextKey int

const (
	ContextKeyRequestMethod contextKey = iota
	ContextKeyResponseHeaders
	ContextKeyResponseSize
	ContextKeyAppID
)

const (
	ResponseCodeSuccess = "000000"
)

type HttpSever struct {
	e            base_enpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []RequestFunc
	after        []ServerResponseFunc
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	errorEncoder ErrorEncoder
}

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type ServerOption func(*HttpSever)

func (s HttpSever) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if len(s.finalizer) > 0 {
		iw := &interceptingWriter{w, http.StatusOK, 0}
		defer func() {
			ctx = context.WithValue(ctx, ContextKeyResponseHeaders, iw.Header())
			ctx = context.WithValue(ctx, ContextKeyResponseSize, iw.written)
			for _, f := range s.finalizer {
				f(ctx, iw.code, r)
			}
		}()
		w = iw
	}

	for _, f := range s.before {
		ctx = f(ctx, r)
	}

	request, err := s.dec(ctx, r)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	response, err := s.e(ctx, request)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	for _, f := range s.after {
		ctx = f(ctx, w)
	}

	if err = s.enc(ctx, w, response); err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}
}

// type RouteHandler interface {
// 	HandleLogin(ctx context.Context, req models.LoginRequest) (*models.BaseResponse, error)
// }

func (s *HttpSever) HandleLogin(ctx context.Context, req models.LoginRequest) (*models.BaseResponse, error) {
	// var reqCtx = GetRequestContext(ctx)
	// req.ClientId = reqCtx.UserID
	// resp, err := s.amqpAdapter.RegisterToken(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.Code != httptransport.ResponseCodeSuccess {
	// 	return nil, errorkit.ServiceError(resp.Code, errorkit.TitleCodeCommonInvalid)
	// }

	// go s.addDeviceToken(reqCtx.DeviceID, req.NotificationToken)
	fmt.Println("pac pac")
	return &models.BaseResponse{
		Code: ResponseCodeSuccess,
		Message: models.BaseMessage{
			Title: "Register token",
			Text:  "Success",
		},
	}, nil
}

type RequestContext struct {
	AccessToken     string
	AppID           string
	AppKey          string
	Channel         string
	Organization    string
	AppVersion      string
	UserID          string
	Username        string
	FullName        string
	Mobile          string
	Email           string
	ClientNo        string
	DeviceID        string
	DeviceModel     string
	DeviceOSName    string
	IPRequest       string
	Language        string
	Timestamp       string
	OtpType         string
	SessionID       string
	EncryptKey      string
	Gender          string
	LoginChannel    string
	VendorUsers     string
	ClientInfo      string
	XRequestId      string
	ActionBy        string
	VIBRequestID    string
	IsDummyUser     string
	SmoDeviceId     string
	DeviceOSVersion string
}

func GetRequestContext(ctx context.Context) RequestContext {
	return RequestContext{
		AppID: GetString(ctx, ContextKeyAppID),
	}
}

func GetString(ctx context.Context, key contextKey) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}

func NewServer(
	e base_enpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *HttpSever {
	s := &HttpSever{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorEncoder: DefaultErrorEncoder,
		errorHandler: NewLogErrorHandler(NewNopLogger()),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

type LogErrorHandler struct {
	logger Logger
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Log("err", err)
}

type Logger interface {
	Log(keyvals ...interface{}) error
}

type nopLogger struct{}

func NewNopLogger() Logger { return nopLogger{} }

func (nopLogger) Log(...interface{}) error { return nil }

func NewLogErrorHandler(logger Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

type MessageLoader interface {
	Populate(ctx context.Context)
}

type StatusCoder interface {
	StatusCode() int
}

type Headerer interface {
	Headers() http.Header
}

func DefaultErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	if msgLoader, ok := err.(MessageLoader); ok {
		msgLoader.Populate(ctx)
	}
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}

	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

func NopRequestDecoder(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
