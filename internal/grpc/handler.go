package grpc

import (
	"context"
	"fmt"

	usdtpb "github.com/DrummDaddy/usdt-rates/gen/gen/usdt/v1"
	"github.com/DrummDaddy/usdt-rates/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	usdtpb.UnimplementedRateServiceServer
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetRates(ctx context.Context, req *usdtpb.GetRatesRequest) (*usdtpb.GetRatesResponse, error) {
	// Validate inputs (1-indexed)
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	if req.N < 1 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("n must be >= 1, got %d", req.N))
	}
	if req.M < 1 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("m must be >= 1, got %d", req.M))
	}
	if req.N > req.M {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("n must be <= m, got n=%d m=%d", req.N, req.M))
	}

	out, err := h.svc.GetRates(ctx, service.GetRatesInput{
		N: int(req.N),
		M: int(req.M),
	})
	if err != nil {
		// For service/internal errors you can map to Internal
		return nil, status.Errorf(codes.Internal, "failed to get rates: %v", err)
	}

	return &usdtpb.GetRatesResponse{
		FetchedAt: timestamppb.New(out.FetchedAt),
		AskTopN:   out.AskTopN.String(),
		AskAvgNM:  out.AskAvgNM.String(),
		BidTopN:   out.BidTopN.String(),
		BidAvgNM:  out.BidAvgNM.String(),
	}, nil
}

func (h *Handler) Healthcheck(ctx context.Context, req *usdtpb.HealthcheckRequest) (*usdtpb.HealthcheckResponse, error) {
	return &usdtpb.HealthcheckResponse{Status: "ok"}, nil
}
