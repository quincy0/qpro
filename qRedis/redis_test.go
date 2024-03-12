package qRedis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type MockStruct struct {
	RoomID string `redis:"room_id"`
	Play   string `redis:"play"`
	Check  bool   `redis:"check"`
	Switch int    `redis:"switch"`
	PlayID int64  `redis:"PlayID"`
}

func TestHSet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	ctx := context.Background()
	c := &AeClient{client: redis.NewClient(&redis.Options{Addr: s.Addr()})}

	err = c.HSet(ctx, "test", "f", "v")
	if err != nil {
		t.Fatal()
	}
	if got := s.HGet("test", "f"); got != "v" {
		t.Fatal()
	}

	data := &MockStruct{
		RoomID: "23456789",
		Play:   "playinfo",
		Check:  false,
		Switch: 0,
		PlayID: 0,
	}

	err = c.HSet(ctx, "mock", data)
	if err != nil {
		t.Fatal()
	}

	if got := s.HGet("mock", "room_id"); got != "23456789" {
		t.Fatal()
	}

}
