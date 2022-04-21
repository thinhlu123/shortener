package delivery

import (
	"context"
	"github.com/thinhlu123/shortener/internal/link"
	linkService "github.com/thinhlu123/shortener/internal/link/pb"
	"github.com/thinhlu123/shortener/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewLinkMicroservice(logger logger.Logger, linkUc link.LinkUsecases) *LinkMicroservice {
	return &LinkMicroservice{
		logger: logger,
		linkUC: linkUc,
	}
}

type LinkMicroservice struct {
	logger logger.Logger
	linkUC link.LinkUsecases
}

func (l LinkMicroservice) GetLink(ctx context.Context, req *linkService.GetLinkReq) (*linkService.GetLinkResp, error) {
	shortLink := req.GetShortLink()
	if shortLink == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid param input")
	}

	oriLink, err := l.linkUC.GetLink(ctx, shortLink)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &linkService.GetLinkResp{
		OriginalLink: oriLink,
	}, nil
}

func (l LinkMicroservice) CreateLink(ctx context.Context, req *linkService.CreateLinkReq) (*linkService.CreateLinkResp, error) {
	// TODO: implement count number of anonymous user
	// TODO: implement devkey

	oriLink := req.GetLink()
	if oriLink == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid param input")
	}

	short, err := l.linkUC.CreateShortLink(ctx, oriLink)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &linkService.CreateLinkResp{
		OriginalLink: oriLink,
		ShortLink:    short,
	}, nil
}
