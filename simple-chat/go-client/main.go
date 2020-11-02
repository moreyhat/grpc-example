package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	pb "github.com/moreyhat/grpc-example/simple-chat/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func postMessage(conn *grpc.ClientConn, message string) error {
	client := pb.NewChatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.PostMessage(ctx, &pb.PostMessageRequest{Message: message})
	if err != nil {
		log.Fatalf("could not post message: %v", err)
		return err
	}
	if r.Result == true {
		fmt.Println("Successfully posted message")
		return nil
	}
	fmt.Println("Failed to post message")
	return errors.New("Failed to post message")
}

func listMessages(conn *grpc.ClientConn) ([]*pb.MessageItem, error) {
	client := pb.NewChatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.ListMessages(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	sort.Slice(r.Item, func(i, j int) bool { return r.Item[i].GetTimeStamp() < r.Item[j].GetTimeStamp() })
	return r.Item, nil
}

func main() {
	postMsgCmd := flag.NewFlagSet("post-message", flag.ExitOnError)
	postMsg := postMsgCmd.String("m", "", "Message to post")

	if len(os.Args) < 2 {
		fmt.Println("Subcommands: list-messages or post-message is necessary")
		return
	}

	conn, err := grpc.Dial("localhost:80", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "list-messages":
		messages, err := listMessages(conn)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to get messages")
		}
		for _, message := range messages {
			fmt.Printf("%v: %v\n", message.TimeStamp, message.Message)
		}
		return
	case "post-message":
		postMsgCmd.Parse(os.Args[2:])
		if *postMsg != "" {
			postMessage(conn, *postMsg)
			return
		}
		fmt.Println("Message is required with '-m' option")
		return
	default:
		fmt.Println("Subcommands: list-messages or post-message is necessary")
		return
	}
}
