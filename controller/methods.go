package controller

import (
	"context"

	pb "github.com/embarkstudios/cassini/api"
	log "github.com/sirupsen/logrus"
)

func (c *Controller) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	log.WithField("context", ctx).Debug("Ping called")
	reply := &pb.PingReply{
		Version: "master",
	}
	return reply, nil
}

func (c *Controller) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	log.WithField("context", ctx).Debug("Get called")
	reply := &pb.GetReply{
		Key: req.Key,
		Url: "https://somewhere/over/the/rainbow",
	}
	return reply, nil
}
