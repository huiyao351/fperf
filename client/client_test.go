package client

import (
	"testing"
)

type testcli struct{}

func (c *testcli) Dial(addr string) error {
	return nil
}

func TestNewClient(t *testing.T) {
	if c := NewClient("test"); c != nil {
		t.Fail()
	}

	clients["test"] = func(flag *FlagSet) Client { return &testcli{} }
	t.Log(clients)
	if c := NewClient("test"); c == nil {
		t.Fail()
	}
}
