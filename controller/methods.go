package controller

import (
	"context"

	pb "github.com/embarkstudios/cassini/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	log.WithField("context", ctx).Debug("Ping called")
	reply := &pb.PingReply{
		Version: "master",
	}
	return reply, nil
}

func (c *Controller) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	log.WithField("key", req.Key).Debug("Get called")
	return nil, status.Error(codes.NotFound, "Key not found")
}

func (c *Controller) Announce(ctx context.Context, req *pb.AnnounceRequest) (*pb.AnnounceReply, error) {
	log.WithField("node", req.Node).Debug("Announce called")
	reply := &pb.AnnounceReply{
		ExpireTimeSeconds: 60,
	}
	return reply, nil
}
