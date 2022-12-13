package services

import (
	"context"
	"net/http"
)

// RequestFunc may take information from an HTTP request and put it into a
// request context. In Servers, RequestFuncs are executed prior to invoking the
// endpoint. In Clients, RequestFuncs are executed after creating the request
// but prior to invoking the HTTP client.
type RequestFunc func(context.Context, *http.Request) context.Context

// ServerResponseFunc may take information from a request context and use it to
// manipulate a ResponseWriter. ServerResponseFuncs are only executed in
// servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(context.Context, http.ResponseWriter) context.Context

// ClientResponseFunc may take information from an HTTP request and make the
// response available for consumption. ClientResponseFuncs are only executed in
// clients, after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(context.Context, *http.Response) context.Context

// ServerFinalizerFunc can be used to perform work at the end of an HTTP
// request, after the response has been written to the client. The principal
// intended use is for request logging. In addition to the response code
// provided in the function signature, additional response parameters are
// provided in the context under keys with the ContextKeyResponse prefix.
type ServerFinalizerFunc func(ctx context.Context, code int, r *http.Request)

type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)

type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)

// EncodeRequestFunc encodes the passed request object into the HTTP request
// object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward EncodeRequestFunc could be something that JSON
// encodes the object directly to the request body.
type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error

type interceptingWriter struct {
	http.ResponseWriter
	code    int
	written int64
}

func ServerBefore(before ...RequestFunc) ServerOption {
	return func(s *HttpSever) { s.before = append(s.before, before...) }
}

func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context {
	// ipRequest, _ := util.GetIP(r)
	for k, v := range map[contextKey]string{
		ContextKeyRequestMethod: r.Method,
		// ContextKeyRequestURI:             r.RequestURI,
		// ContextKeyRequestPath:            r.URL.Path,
		// ContextKeyRequestProto:           r.Proto,
		// ContextKeyRequestHost:            r.Host,
		// ContextKeyRequestRemoteAddr:      r.RemoteAddr,
		// ContextKeyRequestXForwardedFor:   r.Header.Get("X-Forwarded-For"),
		// ContextKeyRequestXForwardedProto: r.Header.Get("X-Forwarded-Proto"),
		// ContextKeyRequestAuthorization:   r.Header.Get("Authorization"),
		// ContextKeyRequestReferer:         r.Header.Get("Referer"),
		// ContextKeyRequestUserAgent:       r.Header.Get("User-Agent"),
		// ContextKeyRequestXRequestID:      r.Header.Get("X-Request-Id"),
		// ContextKeyRequestAccept:          r.Header.Get("Accept"),
		// ContextKeyAccessToken:            r.Header.Get(HeaderAccessToken),
		// ContextKeyAppID:                  r.Header.Get(HeaderAppID),
		// ContextKeyAppKey:                 r.Header.Get(HeaderAppKey),
		// ContextKeyAppVersion:             r.Header.Get(HeaderAppVersion),
		// ContextKeyDeviceID:               r.Header.Get(HeaderDeviceID),
		// ContextKeyDeviceModel:            r.Header.Get(HeaderDeviceModel),
		// ContextKeyDeviceOSName:           r.Header.Get(HeaderDeviceOSName),
		// ContextKeyLanguage:               r.Header.Get(HeaderLanguage),
		// ContextKeyTimestamp:              r.Header.Get(HeaderTimestamp),
		// ContextKeyUserID:                 r.Header.Get(HeaderUserID),
		// ContextKeyUsername:               r.Header.Get(HeaderUsername),
		// ContextKeyFullName:               r.Header.Get(HeaderFullName),
		// ContextKeyMobile:                 r.Header.Get(HeaderMobile),
		// ContextKeyEmail:                  r.Header.Get(HeaderEmail),
		// ContextKeyClientNo:               r.Header.Get(HeaderClientNo),
		// ContextKeyOtpType:                r.Header.Get(HeaderOtpType),
		// ContextKeySessionID:              r.Header.Get(HeaderSessionID),
		// ContextKeyEncryptKey:             r.Header.Get(HeaderEncryptKey),
		// ContextKeyGender:                 r.Header.Get(HeaderGender),
		// ContextKeyVendorUsers:            r.Header.Get(HeaderVendorUsers),
		// ContextKeyClientInfo:             r.Header.Get(HeaderClientInfo),
		// ContextKeyChannel:                r.Header.Get(HeaderChannel),
		// ContextKeyOrganization:           r.Header.Get(HeaderOrganization),
		// ContextKeyLoginChannel:           r.Header.Get(HeaderLoginChannel),
		// ContextKeyActionBy:               r.Header.Get(HeaderActionBy),
		// ContextKeyVIBRequestID:           r.Header.Get(HeaderVIBRequestID),
		// ContextKeyIsDummyUser:            r.Header.Get(HeaderIsDummyUser),
		// ContextKeySmoDeviceId:            r.Header.Get(HeaderSmoDeviceId),
		// ContextKeyDeviceOSVersion:        r.Header.Get(HeaderDeviceOSVersion),
		// ContextKeyIPRequest:              ipRequest,
	} {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}
