package zcli

import "testing"

func TestCli(t *testing.T) {
	c := NewCli("test", "description", "0")
	t.Run("Run AddCommand()", func(t *testing.T) {
		c.AddCommand(&Command{name: "test"})
	})

	t.Run("Run PrintBanner()", func(t *testing.T) {
		c.PrintBanner()
	})

	t.Run("Run IntFlag()", func(t *testing.T) {
		var variable int
		c.IntFlag("int", "description", &variable)
	})
}

type testStruct struct {
	Mode  string `name:"mode" description:"The mode of build"`
	Count int
}

func TestCli_CommandAddFlags(t *testing.T) {
	c := NewCli("test", "description", "0")
	sub := c.NewSubCommand("sub", "sub description")

	ts := &testStruct{}
	sub.AddFlags(ts)

	sub.Action(func() error {
		if ts.Mode != "123" {
			t.Errorf("expected flag value to be set")
		}
		return nil
	})
	e := c.Run("sub", "-mode", "123")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}
