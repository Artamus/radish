package radish

import (
	"reflect"
	"testing"
)

func TestRESPSimpleStrings(t *testing.T) {
	t.Run("it decodes simple strings", func(t *testing.T) {
		cases := []struct {
			input string
			want  string
		}{
			{"+OK\r\n", "OK"},
			{"+HEY\r\n", "HEY"},
		}

		for _, c := range cases {
			got, _ := Decode(c.input)
			gotValue, _ := got.(string)
			assertStringEqual(t, gotValue, c.want)
		}
	})

	t.Run("it fails on incomplete simple strings", func(t *testing.T) {
		cases := []string{
			"", "+", "+OK", "+OK\r",
		}

		for _, c := range cases {
			_, err := Decode(c)
			assertIncompleteRESPError(t, err)
		}
	})
}

func TestRESPBulkStrings(t *testing.T) {
	t.Run("it decodes bulk strings", func(t *testing.T) {
		cases := []struct {
			input string
			want  string
		}{
			{"$2\r\nOK\r\n", "OK"}, {"$3\r\nHEY\r\n", "HEY"},
		}

		for _, c := range cases {
			got, _ := Decode(c.input)
			gotValue, _ := got.(string)
			assertStringEqual(t, gotValue, c.want)
		}
	})

	t.Run("it fails on incomplete bulk strings", func(t *testing.T) {
		cases := []string{"$", "$2", "$2\r", "$2\r\n", "$2\r\nOK", "$2\r\nOK\r"}

		for _, c := range cases {
			_, err := Decode(c)
			assertIncompleteRESPError(t, err)
		}
	})
}

func TestRESPArrays(t *testing.T) {
	t.Run("it decodes arrays", func(t *testing.T) {
		cases := []struct {
			input string
			want  []interface{}
		}{
			{"*1\r\n$4\r\nPING\r\n", []interface{}{"PING"}}, {"*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n", []interface{}{"ECHO", "hey"}},
		}

		for _, c := range cases {
			got, _ := Decode(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		}
	})
}

func assertStringEqual(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertIncompleteRESPError(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("want error, but got none")
	}

	if err != IncompleteRESPError {
		t.Errorf("want incomplete resp error, got %v", err)
	}
}
