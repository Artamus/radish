package radish_test

import (
	"reflect"
	"testing"

	"github.com/Artamus/radish"
)

func TestRESPSimpleStrings(t *testing.T) {
	t.Run("it decodes simple strings", func(t *testing.T) {
		cases := []respTestData{
			{"+OK\r\n", []string{"OK"}},
			{"+HEY\r\n", []string{"HEY"}},
		}

		for _, c := range cases {
			got, _ := radish.Decode(c.encoded)
			assertDecoded(t, got, c.decoded)
		}
	})

	t.Run("it fails on incomplete simple strings", func(t *testing.T) {
		cases := []string{
			"", "+", "+OK", "+OK\r",
		}

		for _, c := range cases {
			_, err := radish.Decode(c)
			assertIncompleteRESPError(t, err)
		}
	})
}

func TestRESPBulkStrings(t *testing.T) {
	t.Run("it decodes bulk strings", func(t *testing.T) {
		cases := []respTestData{
			{"$2\r\nOK\r\n", []string{"OK"}}, {"$3\r\nHEY\r\n", []string{"HEY"}},
		}

		for _, c := range cases {
			got, _ := radish.Decode(c.encoded)
			assertDecoded(t, got, c.decoded)
		}
	})

	t.Run("it fails on incomplete bulk strings", func(t *testing.T) {
		cases := []string{"$", "$2", "$2\r", "$2\r\n", "$2\r\nOK", "$2\r\nOK\r"}

		for _, c := range cases {
			_, err := radish.Decode(c)
			assertIncompleteRESPError(t, err)
		}
	})
}

func TestRESPArrays(t *testing.T) {
	t.Run("it decodes arrays", func(t *testing.T) {
		cases := []respTestData{
			{"*1\r\n$4\r\nPING\r\n", []string{"PING"}}, {"*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n", []string{"ECHO", "hey"}},
		}

		for _, c := range cases {
			got, _ := radish.Decode(c.encoded)
			assertDecoded(t, got, c.decoded)
		}
	})

	t.Run("it fails on incomplete arrays", func(t *testing.T) {
		cases := []string{"*", "*1", "*1\r\n", "*1\r\n$4", "*2\r\n$4\r\nECHO\r\n"}

		for _, c := range cases {
			_, err := radish.Decode(c)
			assertIncompleteRESPError(t, err)
		}
	})
}

func assertDecoded(t testing.TB, got, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v', want '%v'", got, want)
	}
}

func assertIncompleteRESPError(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("want error, but got none")
	}

	if err != radish.ErrIncompleteRESP {
		t.Errorf("want incomplete resp error, got %v", err)
	}
}

type respTestData struct {
	encoded string
	decoded []string
}
