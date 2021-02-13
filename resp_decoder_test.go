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
		_, err := Decode("")
		assertIncompleteRESPError(t, err)

		_, err = Decode("+")
		assertIncompleteRESPError(t, err)

		_, err = Decode("'OK")
		assertIncompleteRESPError(t, err)

		_, err = Decode("+OK\r")
		assertIncompleteRESPError(t, err)
	})
}

func TestRESPBulkStrings(t *testing.T) {
	t.Run("it decodes bulk strings", func(t *testing.T) {
		got, _ := Decode("$2\r\nOK\r\n")
		assertEqual(t, got, "OK")

		got, _ = Decode("$3\r\nHEY\r\n")
		assertEqual(t, got, "HEY")
	})
}

func assertEqual(t testing.TB, got, want string) {
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
