package grpc

import (
	context "context"
	"github.com/funnyecho/code-push/daemon/session"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/log"
)

func New(endpoints Endpoints, logger log.Logger) *sessionServer {
	return &sessionServer{
		endpoints,
		logger,
	}
}

type sessionServer struct {
	endpoints Endpoints
	log.Logger
}

func (s *sessionServer) GenerateAccessToken(ctx context.Context, request *pb.GenerateAccessTokenRequest) (*pb.StringResponse, error) {
	res, err := s.endpoints.GenerateAccessToken(unmarshalAccessTokenClaims(request.GetClaims()))

	return marshalBytesToStringResponse(res), err
}

func (s *sessionServer) VerifyAccessToken(ctx context.Context, request *pb.VerifyAccessTokenRequest) (*pb.VerifyAccessTokenResponse, error) {
	res, err := s.endpoints.VerifyAccessToken([]byte(request.GetToken()))
	return &pb.VerifyAccessTokenResponse{
		Claims: marshalAccessTokenClaims(res),
	}, err
}

func marshalAccessTokenClaims(claims *session.AccessTokenClaims) *pb.AccessTokenClaims {
	if claims == nil {
		return nil
	}

	return &pb.AccessTokenClaims{
		Issuer:   pb.AccessTokenIssuer(claims.Issuer),
		Subject:  claims.Subject,
		Audience: claims.Audience,
	}
}

func marshalBytesToStringResponse(bytes []byte) *pb.StringResponse {
	if bytes == nil {
		return nil
	}

	return &pb.StringResponse{Data: string(bytes)}
}

func unmarshalAccessTokenClaims(request *pb.AccessTokenClaims) *session.AccessTokenClaims {
	if request == nil {
		return nil
	}

	return &session.AccessTokenClaims{
		Issuer:   session.AccessTokenIssuer(request.Issuer),
		Subject:  request.Subject,
		Audience: request.Audience,
	}
}
