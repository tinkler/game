package game

import (
	"context"
	"sync"

	"sfs.ink/liang/game/pkg/logger"
	"sfs.ink/liang/game/pkg/model/char"
	"sfs.ink/liang/game/pkg/model/tank"
	"sfs.ink/liang/game/pkg/model/world"
)

type Game struct {
	ctx      context.Context
	id       int
	world    world.World
	charSrv  *CharSrv
	relaySrv *RelaySrc
}

func NewGame(ctx context.Context, id int) *Game {
	cs := &CharSrv{
		chars: make(map[int]*char.Char),
	}
	cs.log = logger.NewLogger(cs)
	g := &Game{
		ctx:      ctx,
		id:       id,
		world:    world.World{},
		charSrv:  cs,
		relaySrv: newRelaySrv(ctx),
	}
	go g.relaySrv.Serve()
	return g
}

func (g *Game) ID() int {
	return g.id
}

func (g *Game) LoginChar(id int) {
	g.charSrv.addChar(id)
}

func (g *Game) AddClient(ctx context.Context, tk *tank.Tank) <-chan Frame {
	return g.relaySrv.AddClient(ctx, tk)
}

func (g *Game) UpdateClient(tk tank.Tank) Frame {
	return g.relaySrv.UpdateClient(tk)
}

func (g *Game) MoveUp() int {
	return g.relaySrv.MoveUp()
}

func (g *Game) MoveRight() int {
	return g.relaySrv.MoveRight()
}

func (g *Game) MoveDown() int {
	return g.relaySrv.MoveDown()
}

func (g *Game) MoveLeft() int {
	return g.relaySrv.MoveLeft()
}

type CharSrv struct {
	mu    sync.Mutex
	chars map[int]*char.Char
	log   *logger.Logger
}

// clear char from chars map
func clearCharData(chars map[int]*char.Char, c *char.Char) {

}

func (s *CharSrv) ID() string {
	return "1"
}

func (s *CharSrv) TypeName() string {
	return "CharSrv"
}

func (s *CharSrv) addChar(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, exists := s.chars[id]
	if exists {
		// TODO: clear existing char
		clearCharData(s.chars, c)
	}
	s.chars[id] = &char.Char{}
	s.log.Info("%d added", id)
}
