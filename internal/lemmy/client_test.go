package lemmy

import "testing"

func TestNewClient(t *testing.T) {
  client := NewClient("url", "user", "pass")
  if client.User != "user" {
		t.Error("user not set")
  }
}
