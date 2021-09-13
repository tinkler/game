package game

import (
	"context"
	"sync"
	"time"

	"sfs.ink/liang/game/pkg/logger"
	"sfs.ink/liang/game/pkg/model/tank"
	"sfs.ink/liang/game/pkg/util"
)

const fps = time.Second / 10

type RelaySrc struct {
	ctx       context.Context
	id        string
	stepFrame int
	sync.Mutex
	log     *logger.Logger
	clients map[string]*RelayClient
	driving chan struct{} //主动执行帧同步
}

func newRelaySrv(ctx context.Context) *RelaySrc {
	s := &RelaySrc{ctx: ctx}
	s.id = string(util.Krand(10, util.KC_RAND_KIND_NUM))
	s.log = logger.NewLogger(s)
	s.clients = make(map[string]*RelayClient)
	s.driving = make(chan struct{})
	return s
}

func (s *RelaySrc) ID() string {
	return s.id
}

func (s *RelaySrc) TypeName() string {
	return "RelaySrv"
}

func (s *RelaySrc) Serve() {
	var ts *time.Timer
	var times int64 = 0
LOOP:
	for {
		s.Lock()
		if len(s.clients) == 0 {
			s.Unlock()
			// 每次等待时间加长
			ts = time.NewTimer(time.Second * time.Duration(times))
			s.log.Info("空值等待 %d", times)
			select {
			case <-ts.C:
				if times < 10000 {
					times++
				}

				continue LOOP
			case <-s.driving:
				times = 0
				continue LOOP
			case <-s.ctx.Done():
				break LOOP
			}

		}
		// 收集
		var tanks []tank.Tank
		for _, c := range s.clients {
			tanks = append(tanks, *c.tank)
		}
		f := Frame{}
		f.S = s.stepFrame
		f.T = time.Now()
		f.V = tanks
		for _, c := range s.clients {
			go func(client *RelayClient) {
				timeout := time.NewTimer(time.Microsecond * 500)
				select {
				case client.FrameCh <- f:
					return
				case <-timeout.C:
					return
				}
			}(c)
		}
		s.Unlock()
		time.Sleep(fps)
	}
}

func (s *RelaySrc) UpdateClient(data tank.Tank) Frame {
	if data.Name() == "" {
		return Frame{}
	}
	name := data.Name()
	s.Lock()
	defer s.Unlock()
	tk, ok := s.clients[name]
	if !ok {
		s.log.Warn("update unkown tank ", name)
		return Frame{}
	}
	s.stepFrame++
	tk.updateData <- data
	return Frame{S: s.stepFrame, T: time.Now()}
}

func (s *RelaySrc) AddClient(ctx context.Context, tk *tank.Tank) <-chan Frame {

	if tk.Name() == "" {
		return nil
	}
	name := tk.Name()
	s.Lock()
	defer s.Unlock()
	needDriving := len(s.clients) == 0
	_, ok := s.clients[name]
	if ok {
		s.log.Warn("duplicate add tank name ", name)
		return nil
	}
	c := &RelayClient{
		ctx:        ctx,
		Name:       name,
		tank:       tk,
		updateData: make(chan tank.Tank),
		FrameCh:    make(chan Frame),
	}
	s.clients[name] = c
	go func() {
		ticker := time.NewTicker(fps)
		for {
			select {
			case <-ticker.C:
				s.log.Info("update on server ", name)
				c.tank.Update(float64(fps))
			case data, ok := <-c.updateData:
				if !ok {
					return
				}
				c.tank.UpdateTo(data)
			case <-ctx.Done():
				s.RemoveClient(name)
				return
			}
		}
	}()
	if needDriving {
		s.driving <- struct{}{}
	}
	return s.clients[name].FrameCh
}

func (s *RelaySrc) RemoveClient(name string) {
	s.Lock()
	defer s.Unlock()
	tk, ok := s.clients[name]
	if !ok {

		s.log.Warn("remove unkown tank ", name)
		return
	}
	tk.Close()
	delete(s.clients, name)
	s.log.Info("remove tank ", name)
}

func (s *RelaySrc) MoveUp() int {
	s.Lock()
	defer s.Unlock()
	s.stepFrame++
	s.log.Info("Up")
	return s.stepFrame
}

func (s *RelaySrc) MoveRight() int {
	s.Lock()
	defer s.Unlock()
	s.stepFrame++
	s.log.Info("Right")
	return s.stepFrame
}

func (s *RelaySrc) MoveDown() int {
	s.Lock()
	defer s.Unlock()
	s.stepFrame++
	s.log.Info("Down")
	return s.stepFrame
}

func (s *RelaySrc) MoveLeft() int {
	s.Lock()
	defer s.Unlock()
	s.stepFrame++
	s.log.Info("Left")
	return s.stepFrame
}

type RelayClient struct {
	// client context
	ctx        context.Context
	Name       string
	tank       *tank.Tank
	updateData chan tank.Tank
	FrameCh    chan Frame
}

func (c *RelayClient) Close() {
	close(c.updateData)
	close(c.FrameCh)
}
