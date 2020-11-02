package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/moreyhat/grpc-example/simple-chat/pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedChatServer
}

func (s *server) PostMessage(ctx context.Context, in *pb.PostMessageRequest) (*pb.PostMessageResponse, error) {
	log.Printf("Received: %v", in.GetMessage())
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	t := time.Now()
	err := rdb.HSet(context.Background(), "message", t.String(), in.GetMessage()).Err()
	if err != nil {
		return &pb.PostMessageResponse{Result: false}, nil
	}
	return &pb.PostMessageResponse{Result: true}, nil
}

func (s *server) ListMessages(ctx context.Context, in *empty.Empty) (*pb.ListMessagesResponse, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	messages, err := rdb.HGetAll(context.Background(), "message").Result()
	if err != nil {
		return nil, err
	}
	res := []*pb.MessageItem{}
	for timeStamp, message := range messages {
		res = append(res, &pb.MessageItem{
			TimeStamp: timeStamp,
			Message:   message,
		})
	}
	return &pb.ListMessagesResponse{Item: res}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
