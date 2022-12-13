package controller

import (
	"context"

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

type RouteHandler interface {
	HandleLogin(ctx context.Context, req models.LoginRequest) (*models.BaseResponse, error)
}

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type Controller struct {
	e            base_enpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []RequestFunc
	after        []ServerResponseFunc
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	errorEncoder ErrorEncoder
}

func (s *Controller) HandleLogin(ctx context.Context, req models.LoginRequest) (*models.BaseResponse, error) {
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
