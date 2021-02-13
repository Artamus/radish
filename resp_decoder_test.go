package radish

import "testing"

func TestRESPDecoder(t *testing.T) {
	t.Run("it decodes simple strings", func(t *testing.T) {
		assertEqual(t, Decode("+OK\r\n"), "OK")
		assertEqual(t, Decode("+HEY\r\n"), "HEY")
	})
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
