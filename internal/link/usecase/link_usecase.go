package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/thinhlu123/shortener/config"
	"github.com/thinhlu123/shortener/internal/link"
	"github.com/thinhlu123/shortener/internal/models"
	"github.com/thinhlu123/shortener/pkg/auth"
	"github.com/thinhlu123/shortener/pkg/cache"
	"github.com/thinhlu123/shortener/pkg/logger"
	"github.com/thinhlu123/shortener/pkg/utils"
	"math/big"
)

type LinkUsecase struct {
	linkRepo link.LinkRepository
	redis    *cache.Cache
	logger   logger.Logger
}

func NewLinkUsecase(linkRepo link.LinkRepository, log logger.Logger, redis *cache.Cache) *LinkUsecase {
	return &LinkUsecase{
		linkRepo: linkRepo,
		logger:   log,
		redis:    redis,
	}
}

func (l *LinkUsecase) CreateShortLink(ctx context.Context, oriLink string) (string, error) {
	var (
		shortLink = ""
		err       error
	)

	token := utils.GetFromMetadata(ctx, "Authorization")
	usr, err := auth.GetUsernameFromToken(token)
	if err != nil {
		return "", err
	}

	shortLink, err = generateShortLink(oriLink)
	if err != nil {
		return "", err
	}

	if err = l.linkRepo.CreateShortLink(ctx, models.Link{
		OriginalUrl: oriLink,
		ShortUrl:    shortLink,
		Username:    usr,
	}); err != nil {
		return "", err
	}

	return shortLink, nil
}

func (l *LinkUsecase) GetLink(ctx context.Context, shortLink string) (string, error) {
	// TODO: get from cache

	// get from db
	filter := models.Link{
		ShortUrl: shortLink,
	}
	orLink, err := l.linkRepo.GetLink(ctx, filter)
	if err != nil {
		return "", err
	}

	// increase number view short link
	err = l.linkRepo.IncreaseClickCount(ctx, filter)
	if err != nil {
		return "", err
	}

	return orLink, nil
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func generateShortLink(initialLink string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	urlHashBytes := sha256Of(initialLink + u.String())
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}
	return config.Conf.Url.Host + finalString[:8], nil
}
