package radish

import "testing"

func TestRESPSimpleStrings(t *testing.T) {
	t.Run("it decodes simple strings", func(t *testing.T) {
		got, _ := Decode("+OK\r\n")
		assertEqual(t, got, "OK")

		got, _ = Decode("+HEY\r\n")
		assertEqual(t, got, "HEY")
	})

	t.Run("it fails on incomplete strings", func(t *testing.T) {
		_, err := Decode("+")

		if err == nil {
			t.Error("wanted error, but got none")
		}
	})
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
