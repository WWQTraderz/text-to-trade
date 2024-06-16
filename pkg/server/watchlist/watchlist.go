package watchlist

import (
	"context"

	"log"

	"github.com/tjons/text-to-trade/pkg/api/watchlist"
	"github.com/tjons/text-to-trade/pkg/model"
	"gorm.io/gorm"
)

type WatchlistServer struct {
	watchlist.UnimplementedWatchlistServiceServer
	db *gorm.DB
}

func NewWatchlistServer(db *gorm.DB) watchlist.WatchlistServiceServer {
	return &WatchlistServer{db: db}
}

func (s *WatchlistServer) GetWatchlist(ctx context.Context, req *watchlist.WatchlistRequest) (*watchlist.WatchlistResponse, error) {
	log.Println("Getting watchlist")

	watchlistData := &model.Watchlist{}
	res := s.db.Where("user_id = ?", req.UserId).First(watchlistData)
	if res.Error != nil {
		log.Println(res.Error)
		return nil, res.Error
	}
	resp := &watchlist.WatchlistResponse{
		Watchlist: &watchlist.Watchlist{
			Id:      uint32(watchlistData.ID),
			UserId:  uint32(watchlistData.UserID),
			Symbols: watchlistData.Symbols,
			Name:    watchlistData.Name,
		},
	}

	return resp, nil
}

func (s *WatchlistServer) GetWatchlists(ctx context.Context, req *watchlist.WatchlistRequest) (*watchlist.WatchlistListResponse, error) {
	return nil, nil
}

func (s *WatchlistServer) CreateWatchlist(ctx context.Context, req *watchlist.Watchlist) (*watchlist.WatchlistResponse, error) {
	watchlistData := &model.Watchlist{
		UserID:  uint(req.UserId),
		Symbols: req.Symbols,
		Name:    req.Name,
	}
	res := s.db.Create(watchlistData)
	if res.Error != nil {
		log.Println(res.Error)
		return nil, res.Error
	}
	resp := &watchlist.WatchlistResponse{
		Watchlist: &watchlist.Watchlist{
			Id:      uint32(watchlistData.ID),
			UserId:  uint32(watchlistData.UserID),
			Symbols: watchlistData.Symbols,
			Name:    watchlistData.Name,
		},
	}

	return resp, nil
}

func (s *WatchlistServer) UpdateWatchlist(ctx context.Context, req *watchlist.Watchlist) (*watchlist.WatchlistResponse, error) {
	watchlistData := &model.Watchlist{
		UserID:  uint(req.UserId),
		Symbols: req.Symbols,
		Name:    req.Name,
	}
	watchlistData.ID = uint(req.Id)

	res := s.db.Save(watchlistData)
	if res.Error != nil {
		log.Println(res.Error)
		return nil, res.Error
	}
	resp := &watchlist.WatchlistResponse{
		Watchlist: &watchlist.Watchlist{
			Id:      uint32(watchlistData.ID),
			UserId:  uint32(watchlistData.UserID),
			Symbols: watchlistData.Symbols,
			Name:    watchlistData.Name,
		},
	}

	return resp, nil
}
