package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardForHeader          = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *GrpcServer) extraMetadata(ctx context.Context) *Metadata {
	mtdata := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) != 0 {
			mtdata.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) != 0 {
			mtdata.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwardForHeader); len(clientIps) != 0 {
			mtdata.ClientIp = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdata.ClientIp = p.Addr.String()
	}
	return mtdata
}
