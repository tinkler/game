package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "sfs.ink/liang/game/api/proto/game"
	"sfs.ink/liang/game/pkg/util"
)

func main() {
	// conn, err := rpc.DialHTTP("tcp", "127.0.0.1:9301")
	// if err != nil {
	// 	log.Fatalln("dailing error: ", err)
	// }
	// req := myrpc.LoginRequest{ID: 10}
	// var res myrpc.LoginResponse
	// err = conn.Call("GameRpc.LoginChar", req, &res)
	// if err != nil {
	// 	log.Fatalln("Login game error: ", err)
	// }
	// fmt.Printf("Returnning session %s\n", res.Session)

	conn, err := grpc.Dial("ssdf.sedns.cn:9301", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGameClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := c.LoginChar(ctx, &pb.LoginRequest{Id: 10})
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	log.Printf("Logged in: %s", r.GetSession())
	m, err := c.MoveChar(ctx, &pb.MoveRequest{Direction: pb.MoveRequest_UP})
	if err != nil {
		log.Fatalf("could not move: %v", err)
	}
	log.Printf("Move success in step frame %d at %s", m.S, m.T.AsTime().Format("2006-01-02 15:04:05"))
	name := string(util.Krand(10, util.KC_RAND_KIND_ALL))
	go func() {
		of := pb.Offset{}
		for {
			time.Sleep(time.Second * 5)
			of.Dx += 20
			_, err := c.UpdateTank(ctx, &pb.TankAttr{Name: name, Position: &of})
			if err != nil {
				return
			}
		}
	}()

	stream, err := c.RelayTank(ctx, &pb.TankAttr{Name: name, Position: &pb.Offset{Dx: 10, Dy: 20}, BodyAngle: 30, TurretAngle: 100})
	if err != nil {
		log.Fatalf("%v.RelayTank(_) = _, %v", c, err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.RelayTank(_) = _, %v", c, err)
		}
		for _, t := range res.Tanks {
			log.Printf("Name %s, BodyAngle %f, TurretAngle %f, Position Dx %f Dy %f\n Frame: %d, %s\n", t.Name, t.BodyAngle, t.TurretAngle, t.Position.Dx, t.Position.Dy, res.Step.S, res.Step.T.AsTime().Format("15:04:05"))
		}

	}
}
