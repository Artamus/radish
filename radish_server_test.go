package radish_test

import (
	"strconv"
	"testing"

	"github.com/Artamus/radish"
	"github.com/go-redis/redis"
)

func TestRadishServer(t *testing.T) {
	t.Run("it responds PONG to PING", func(t *testing.T) {
		server, _ := radish.NewRadishServer(6379)
		go func() {
			server.Listen()
		}()

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		got := rdb.Ping().Val()
		want := "PONG"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

func makeRedisClient(port int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
