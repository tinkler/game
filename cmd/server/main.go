package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	pb "sfs.ink/liang/game/api/proto/game"
	"sfs.ink/liang/game/pkg/attr"
	"sfs.ink/liang/game/pkg/game"
	"sfs.ink/liang/game/pkg/logger"
	"sfs.ink/liang/game/pkg/model/tank"
	"sfs.ink/liang/game/pkg/util"
)

var serverNodeID string

func init() {
	serverNodeID = string(util.Krand(10, util.KC_RAND_KIND_NUM))
}

type serverNode struct {
	log  *logger.Logger
	stub *game.Game
	pb.UnimplementedGameServer
}

func (n *serverNode) ID() string {
	return "ServerNode" + serverNodeID
}

func (n *serverNode) TypeName() string {
	return "ServerNode"
}

func newServerNode(ctx context.Context) *serverNode {
	node := &serverNode{}
	node.log = logger.NewLogger(node)
	g := game.NewGame(ctx, 1)
	node.stub = g
	return node
}

// LoginChar 登陆角色
func (s *serverNode) LoginChar(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	ses := util.Krand(32, util.KC_RAND_KIND_ALL)
	s.stub.LoginChar(int(in.Id))
	return &pb.LoginResponse{
		Session: s.ID() + "_" + string(ses),
	}, nil
}

// MoveChar 移动角色
func (s *serverNode) MoveChar(ctx context.Context, in *pb.MoveRequest) (*pb.StepFrame, error) {
	res := &pb.StepFrame{T: timestamppb.Now()}
	switch in.Direction {
	case pb.MoveRequest_UP:
		res.S = int64(s.stub.MoveUp())
	case pb.MoveRequest_RIGHT:
		res.S = int64(s.stub.MoveRight())
	case pb.MoveRequest_DOWN:
		res.S = int64(s.stub.MoveDown())
	case pb.MoveRequest_LEFT:
		res.S = int64(s.stub.MoveLeft())
	}
	return res, nil
}

// RelayTank 坦克帧同步
// 服务器流数据形式,在属性输入后持续返回游戏中所有坦克的位置[position]
// 状态:
// bodyAngle - 身体转向角度
// targetBodyAngle - 目标身体转向角度
// turretAngle - 炮台转向角度
// targetTurretAngle - 目标炮台转向角度
func (s *serverNode) RelayTank(in *pb.TankAttr, stream pb.Game_RelayTankServer) error {

	t := tank.NewTank(in.Name, attr.Offset{Dx: in.Position.Dx, Dy: in.Position.Dy}, float64(in.BodyAngle), float64(in.TargetBodyAngle), float64(in.TurretAngle), float64(in.TargetTurretAngle))
	frames := s.stub.AddClient(stream.Context(), &t)
	if frames == nil {
		return nil
	}
	for frame := range frames {
		var ts []*pb.TankAttr
		for _, t := range frame.V.([]tank.Tank) {
			dx, dy, ba, ta := t.Attr()
			ts = append(ts, &pb.TankAttr{
				Name:        t.Name(),
				Position:    &pb.Offset{Dx: dx, Dy: dy},
				BodyAngle:   float32(ba),
				TurretAngle: float32(ta),
			})
		}
		if err := stream.Send(&pb.TanksAttr{Tanks: ts, Step: &pb.StepFrame{S: int64(frame.S), T: timestamppb.New(frame.T)}}); err != nil {
			return err
		}
	}
	return nil
}

// UpdateTank 更新坦克属性
func (s *serverNode) UpdateTank(ctx context.Context, in *pb.TankAttr) (*pb.StepFrame, error) {
	tk := tank.NewTank(in.Name, attr.Offset{Dx: in.Position.Dx, Dy: in.Position.Dy}, float64(in.BodyAngle), float64(in.TargetBodyAngle), float64(in.TurretAngle), float64(in.TargetTurretAngle))
	f := s.stub.UpdateClient(tk)
	pf := &pb.StepFrame{
		S: int64(f.S),
		T: timestamppb.New(f.T),
	}
	return pf, nil
}

func runGrpc() {
	sn := newServerNode(context.Background())
	lis, err := net.Listen("tcp", ":9301")
	if err != nil {
		sn.log.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGameServer(s, sn)
	if err := s.Serve(lis); err != nil {
		sn.log.Error("failed to serve: %v", err)
	}
}

func main() {
	runGrpc()
}
