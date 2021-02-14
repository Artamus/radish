package radish_test

import (
	"strconv"
	"testing"

	"github.com/Artamus/radish"
	"github.com/go-redis/redis"
)

var dummyStorage = make(map[string]string)

func TestRadishServer(t *testing.T) {
	t.Run("it responds PONG to the PING command", func(t *testing.T) {
		server := mustMakeRadishServer(t, dummyStorage)
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
		server := mustMakeRadishServer(t, dummyStorage)
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
		server := mustMakeRadishServer(t, dummyStorage)
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
		server := mustMakeRadishServer(t, dummyStorage)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Echo("hey").Val()
		assertResponse(t, got, "hey")
	})

	t.Run("it gets nil value with GET when data does not exist", func(t *testing.T) {
		server := mustMakeRadishServer(t, dummyStorage)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Get("somekey")
		if got.Err() != redis.Nil {
			t.Errorf("got '%v' but wanted nil", got)
		}
	})

	t.Run("it fetches value with GET", func(t *testing.T) {
		mockStorage := make(map[string]string)
		mockStorage["somekey"] = "somevalue"
		server := mustMakeRadishServer(t, mockStorage)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		got := client.Get("somekey").Val()
		assertResponse(t, got, "somevalue")
	})

	t.Run("it saves value with SET", func(t *testing.T) {
		spyStorage := make(map[string]string)
		server := mustMakeRadishServer(t, spyStorage)
		defer server.Close()
		go func() {
			server.Listen()
		}()
		client := makeRedisClient(6379)
		defer client.Close()

		client.Set("otherkey", "othervalue", 0)

		if spyStorage["otherkey"] != "othervalue" {
			t.Errorf("expected new value to be in storage, but it wasn't")
		}
	})
}

func mustMakeRadishServer(t testing.TB, storage map[string]string) *radish.RadishServer {
	t.Helper()

	server, err := radish.NewRadishServer(6379, storage)
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
