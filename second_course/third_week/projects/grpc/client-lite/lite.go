package main

import (
	"context"
	"fmt"
	"go_mail_Coursera/second_course/third_week/projects/grpc/session"
	"log"

	"google.golang.org/grpc"
)

func main() {
	grpcCon, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("couldn't connect to grpc")
	}
	defer grpcCon.Close()

	sessManager := session.NewAuthCheckerClient(grpcCon)

	ctx := context.Background()

	sessId, err := sessManager.Create(ctx, &session.Session{
		Login:     "rvasily",
		Useragent: "chrome",
	})
	fmt.Println("sessId", sessId, err)

	sess, err := sessManager.Check(ctx, &session.SessionID{
		ID: sessId.ID,
	})
	fmt.Println("sess", sess, err)

	_, err = sessManager.Delete(ctx, &session.SessionID{
		ID: sessId.ID,
	})

	sess, err = sessManager.Check(ctx, &session.SessionID{
		ID: sessId.ID,
	})
	fmt.Println("sess", sess, err)

}
