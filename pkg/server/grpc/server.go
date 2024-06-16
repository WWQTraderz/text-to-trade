package grpc

import (
	"context"
	"net"

	chatpb "github.com/tjons/text-to-trade/pkg/api/chat"
	userpb "github.com/tjons/text-to-trade/pkg/api/user"
	watchlistpb "github.com/tjons/text-to-trade/pkg/api/watchlist"
	"github.com/tjons/text-to-trade/pkg/model"
	"github.com/tjons/text-to-trade/pkg/server/chat"
	"github.com/tjons/text-to-trade/pkg/server/user"
	"github.com/tjons/text-to-trade/pkg/server/watchlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Run starts the example gRPC service.
// "network" and "address" are passed to net.Listen.
func Run(ctx context.Context, network, address string) error {
	db, err := model.Connect()
	if err != nil {
		return err
	}

	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			grpclog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	s := grpc.NewServer()
	watchlistpb.RegisterWatchlistServiceServer(s, watchlist.NewWatchlistServer(db))
	chatpb.RegisterChatServer(s, chat.NewChatServer(db))
	userpb.RegisterUserServiceServer(s, user.NewUserService(db))
	// TODO: register services here

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}
