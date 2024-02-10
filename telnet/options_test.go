package telnet_test

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/uaxe/infra/telnet"
)

func TestServerNAWS(t *testing.T) {
	client, server := net.Pipe()
	text := []byte("hello\n")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer client.Close()
		b := make([]byte, 3)
		_, err := client.Read(b)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(b, []byte{255, 253, 31}) {
			t.Errorf("Expected IAC DO NAWS, received %v", b)
		}
		// IAC WILL NAWS IAC SB NAWS W[1] W[0] H[1] H[0] IAC SE
		payload := []byte{255, 251, 31, 255, 250, 31, 0, 80, 0, 20, 255, 240}
		payload = append(payload, text...)
		_, err = client.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()
	conn := telnet.NewConnection(server, []telnet.Option{telnet.NAWSOption})
	b := make([]byte, 32)
	n, err := conn.Read(b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(text, b[:n]) {
		t.Errorf("Expected %q, got %q", text, b[:n])
	}
	wg.Wait()
	conn.Close()

	nw := conn.OptionHandlers[31].(*telnet.NAWSHandler)
	if nw.Width != 80 || nw.Height != 20 {
		t.Logf("%#v", conn)
		t.Errorf("Expected w %d, h %d, got w %d, h %d", 80, 20, nw.Width, nw.Height)
	}
}

func TestClientNAWS(t *testing.T) {
	client, server := net.Pipe()
	go func() {
		conn := telnet.NewConnection(client, []telnet.Option{telnet.ExposeNAWS})
		b := make([]byte, 32)
		_, _ = conn.Read(b)
		_ = conn.Close()
	}()

	_, err := server.Write([]byte{255, 253, 31})
	if err != nil {
		t.Error(err)
	}
	expected := []byte{255, 251, 31, 255, 250, 31, 0, 0, 0, 0, 255, 240}
	buf := bytes.NewBuffer(nil)
	_ = server.SetReadDeadline(time.Now().Add(time.Second))
	n, err := io.CopyN(buf, server, int64(len(expected)))
	if err != nil {
		t.Error(err)
	}
	b := buf.Bytes()
	_ = server.Close()
	// IAC WILL NAWS IAC SB NAWS W[1] W[0] H[1] H[0] IAC SE
	if !bytes.Equal(b[:n], expected) {
		t.Errorf("Expected %v, received %v", expected, b[:n])
	}
}

func TestOnlyServerSupportsNAWS(t *testing.T) {
	client, server := net.Pipe()
	text := []byte("hello\n")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		conn := telnet.NewConnection(client, []telnet.Option{})
		b := make([]byte, 32)
		n, err := conn.Read(b)
		if err != nil {
			t.Error(err)
		}
		_ = conn.Close()

		if !bytes.Equal(text, b[:n]) {
			t.Errorf("Expected %q, got %q", text, b[:n])
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		conn := telnet.NewConnection(server, []telnet.Option{telnet.NAWSOption})
		go func() { _, _ = io.Copy(io.Discard, conn) }()
		_, err := conn.Write(text)
		if err != nil {
			t.Error(err)
		}
		_ = conn.Close()
		wg.Done()
	}()
	wg.Wait()
}
