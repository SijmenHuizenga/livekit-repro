package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	lksdk "github.com/livekit/server-sdk-go"
)

func main() {
	room, err := lksdk.ConnectToRoom(os.Getenv("LIVEKIT_HOST"), lksdk.ConnectInfo{
		APIKey:              os.Getenv("LIVEKIT_API_KEY"),
		APISecret:           os.Getenv("LIVEKIT_SECRET_KEY"),
		RoomName:            "robot-blokje",
		ParticipantIdentity: "robot",
	}, &lksdk.RoomCallback{
		OnReconnected: func() {
			println("reconnected to room")
		},
		OnReconnecting: func() {
			println("reconnecting to room")
		},
		OnDisconnected: func() {
			println("disconnected from room")
		},
	})
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	track, err := lksdk.NewLocalFileTrack("recording.h264",
		// 15 fps
		lksdk.ReaderTrackWithFrameDuration(time.Second/15),
		lksdk.ReaderTrackWithOnWriteComplete(func() { fmt.Println("track finished") }),
	)
	if err != nil {
		panic(err)
	}
	if _, err = room.LocalParticipant.PublishTrack(track, &lksdk.TrackPublicationOptions{
		Name:        "camera",
		VideoWidth:  1280,
		VideoHeight: 720,
	}); err != nil {
		panic(err)
	}

	<-sigChan
	room.Disconnect()

}
