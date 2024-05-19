package grpcApi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent" // if the request is coming from an http client
	xForwardedForHeader        = "x-forwarded-for"        // if the request is coming from an http client
	userAgentHeader            = "user-agent"             // if the request is coming from a grpc client
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

// We are using GRPC gateway, so the request coming to the rpc handler functions can
// be both from a GRPC client like (evans cli) or as well as an HTTP client (like postman).
// Therefore the metadata they send can be stored in different formats in the context.
func (server *server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	// md is of type MD map[string][]string
	if md, ok := metadata.FromIncomingContext(ctx); ok {

		// to get the user agent if the request is coming from an HTTP client
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// to get the user agent if the request is coming from a GRPC client
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// to get the client IP address if the request is coming from an HTTP client
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	// to get the client IP address if the request is coming from a GRPC client
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
