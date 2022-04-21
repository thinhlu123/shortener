package server

import (
	"context"
	"github.com/thinhlu123/shortener/config"
	"github.com/thinhlu123/shortener/internal/link/delivery"
	linkService "github.com/thinhlu123/shortener/internal/link/pb"
	"github.com/thinhlu123/shortener/internal/link/repository"
	"github.com/thinhlu123/shortener/internal/link/usecase"
	delivery2 "github.com/thinhlu123/shortener/internal/user/delivery"
	pd "github.com/thinhlu123/shortener/internal/user/pb"
	repository2 "github.com/thinhlu123/shortener/internal/user/repository"
	usecase2 "github.com/thinhlu123/shortener/internal/user/usecase"
	"github.com/thinhlu123/shortener/pkg/cache"
	"github.com/thinhlu123/shortener/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	logger logger.Logger
	redis  *cache.Cache
}

func NewServer(logger logger.Logger, cache *cache.Cache) *Server {
	return &Server{logger: logger, redis: cache}
}

func (s *Server) Run() error {
	linkRepo := repository.NewLinkRepo()
	linkUsecase := usecase.NewLinkUsecase(linkRepo, s.logger, s.redis)

	userRepo := repository2.NewUserRepository()
	userUsecase := usecase2.NewUserUsecase(userRepo, s.logger)

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: config.Conf.Server.MaxConnectionIdle * time.Minute,
		Timeout:           config.Conf.Server.Timeout * time.Second,
		MaxConnectionAge:  config.Conf.Server.MaxConnectionAge * time.Minute,
		Time:              config.Conf.Server.Timeout * time.Minute,
	}),
	//grpc.UnaryInterceptor(s.Logger),
	//grpc.ChainUnaryInterceptor(
	//	grpc_ctxtags.UnaryServerInterceptor(),
	//	grpc_prometheus.UnaryServerInterceptor,
	//	grpcrecovery.UnaryServerInterceptor(),
	//),
	)

	linkGrpcService := delivery.NewLinkMicroservice(s.logger, linkUsecase)
	linkService.RegisterLinkServiceServer(server, linkGrpcService)

	userGrpcService := delivery2.NewUserMicroservice(userUsecase, s.logger)
	pd.RegisterUserServiceServer(server, userGrpcService)

	lis, err := net.Listen("tcp", config.Conf.Server.PprofPort)
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		s.logger.Infof("Serve gRPC serve at port: ", config.Conf.Server.PprofPort)
		s.logger.Fatal(server.Serve(lis))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
