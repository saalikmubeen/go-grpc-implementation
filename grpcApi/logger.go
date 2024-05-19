package grpcApi

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ** GRPC INTERCEPTOR **
// GRPC interceptor to log requests in the terminal
// UnaryServerInterceptor:
// func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	// To print the logs in the terminal without any JSON formatting
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	startTime := time.Now()
	result, err := handler(ctx, req) // call the service handler to process the request
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

/*
* HTTP MIDDLEWARE
The above GrpcLogger middleware will  work only for gRPC requests.
If we try to send an HTTP request to the gateway server, we won't see any
logs written. That's because we're using in-process translation on our
gateway server . So the gateway will call the RPC handler function directly
without going through any interceptors. If we run the gateway as a separate
server, And use cross-process translation to call the gRPC server via a network
call, Then the logs will show up in the gRPC server as normal.
But that's another network hop that can increase the duration of the request.
So, if we still want to keep using in-process translation, we will have to write a
separate HTTP middleware to log the HTTP requests.
*/
func HttpLogger(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// To print the logs in the terminal without any JSON formatting
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		startTime := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).

			// status code and status text are tricky
			// to get because handler.ServeHTTP(rec, req) doesn't return
			// anything so we don't know abou the status code and status text.
			// So we need to track the status code and status text manually
			// in the ResponseRecorder struct and use it here.

			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).

			//
			Dur("duration", duration).
			Msg("received a HTTP request")
	})

}

// http.ResponseWriter is an interface:
// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

// ResponseRecorder struct that implements the http.ResponseWriter interface
// to track the status code and response body of the HTTP response from the
// handler.
type ResponseRecorder struct {
	http.ResponseWriter // embed the oroginal http.ResponseWriter
	// this StatusCode field will be updated to the correct value when the
	// WriteHeader method is called by the handler function.
	StatusCode int
	// In case of the error, it will be written to the response body
	// so we need to track the response body as well.
	Body []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode) // this will be called by the handler function
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body) // this will be called by the handler function
}
