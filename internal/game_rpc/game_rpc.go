package game_rpc

import (
	"context"
	"strconv"

	"sfs.ink/liang/game/internal/game_srv"
	"sfs.ink/liang/game/pkg/game"
	"sfs.ink/liang/game/pkg/logger"
	"sfs.ink/liang/game/pkg/util"

	myrpc "sfs.ink/liang/game/pkg/rpc"
)

type GameRpc struct {
	sub *game_srv.GameSrv
	sta rpcStatus
	log *logger.Logger
}

type rpcStatus struct {
	session string
}

func (sta rpcStatus) ID() string {
	return sta.session
}

func (r rpcStatus) TypeName() string {
	return "GameRpc"
}

func NewGameRpc() *GameRpc {
	id := 1
	r := new(GameRpc)
	l := logger.NewLogger(r.sta)
	r.sub = game_srv.NewGameSrv(game.NewGame(context.Background(), id))
	r.log = l
	r.sta.session = strconv.Itoa(id)
	return r
}

func (r *GameRpc) LoginChar(req myrpc.LoginRequest, res *myrpc.LoginResponse) error {
	res.Session = string(util.Krand(12, util.KC_RAND_KIND_ALL))
	r.sub.LoginChar(req.ID)
	return nil
}
