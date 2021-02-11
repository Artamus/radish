package radish_test

import (
	"strconv"
	"testing"

	"github.com/Artamus/radish"
	"github.com/go-redis/redis"
)

func TestRadishServer(t *testing.T) {
	t.Run("it responds PONG to PING", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		rdb := makeRedisClient(6379)

		got := rdb.Ping().Val()
		assertResponse(t, got, "PONG")
	})

	t.Run("it responds to multiple commands from the same client", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		rdb := makeRedisClient(6379)

		want := "PONG"
		got := rdb.Ping().Val()
		assertResponse(t, got, want)

		got = rdb.Ping().Val()
		assertResponse(t, got, want)
	})
}

func mustMakeRadishServer(t testing.TB) *radish.RadishServer {
	t.Helper()

	server, err := radish.NewRadishServer(6379)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	return server
}

func makeRedisClient(port int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func assertResponse(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
