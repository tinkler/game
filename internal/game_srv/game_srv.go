package game_srv

import (
	"context"
	"strconv"

	"sfs.ink/liang/game/pkg/game"
)

type GameSrv struct {
	game.Game
}

func NewGameSrv(g *game.Game) *GameSrv {
	return &GameSrv{
		Game: *game.NewGame(context.Background(), 0),
	}
}

func (s *GameSrv) ID() string {
	return strconv.Itoa(s.Game.ID())
}
