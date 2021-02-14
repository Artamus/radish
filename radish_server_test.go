package radish_test

import (
	"strconv"
	"testing"

	"github.com/Artamus/radish"
	"github.com/go-redis/redis"
)

func TestRadishServer(t *testing.T) {
	t.Run("it responds PONG to the PING command", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Ping().Val()

		assertResponse(t, got, "PONG")
	})

	t.Run("it responds to multiple commands from the same client", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Ping().Val()
		assertResponse(t, got, "PONG")

		got = client.Ping().Val()
		assertResponse(t, got, "PONG")
	})

	t.Run("it allows multiple clients to send commands", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()

		client1 := makeRedisClient(6379)
		defer client1.Close()
		client2 := makeRedisClient(6379)
		defer client2.Close()

		assertResponse(t, client1.Ping().Val(), "PONG")
		assertResponse(t, client2.Ping().Val(), "PONG")
	})

	t.Run("it responds to ECHO", func(t *testing.T) {
		server := mustMakeRadishServer(t)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Echo("hey").Val()
		assertResponse(t, got, "hey")
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
		Addr:     "0.0.0.0:" + strconv.Itoa(port),
		Password: "",
		DB:       0,
	})
}

func assertResponse(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
