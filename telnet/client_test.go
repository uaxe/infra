package telnet_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uaxe/infra/telnet"
)

func TestTelnetClient(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	cli := telnet.NewClient(ctx,
		telnet.WithIP("192.168.1.1"),
		telnet.WithPort("23"),
		telnet.WithAuth(true),
		telnet.WithUserName("admin"),
		telnet.WithPassWord("Fh@0217A0"),
	)

	defer func() { _ = cli.Close() }()

	if err := cli.Run(); err != nil {
		t.Fatal(err)
	}

	if err := cli.Write([]byte(`ls`)); err != nil {
		t.Fatal(err)
	}

	if b, err := cli.Read(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(b))
	}

	for {
		select {
		case <-ctx.Done():
			return
		}
	}

}
