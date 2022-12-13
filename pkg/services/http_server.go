package services

import (
	"context"
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

type RouteHandler interface {
	HandleLogin(ctx context.Context, req models.LoginRequest) (*models.BaseResponse, error)
}

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
