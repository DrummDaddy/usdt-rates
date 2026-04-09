package grpc

import (
	"context"

	usdtpb "github.com/DrummDaddy/usdt-rates/gen/usdt/v1"
	"github.com/DrummDaddy/usdt-rates/internal/service"
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
	out, err := h.svc.GetRates(ctx, service.GetRatesInput{
		N: int(req.N),
		M: int(req.M),
	})
	if err != nil {
		return nil, err
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
