package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestSanitizeURI(t *testing.T) {
	uri, err := sanitizeURI(" http://example.com ")
	if err != nil || uri != "http://example.com" {
		t.Fatalf("unexpected sanitize result %s %v", uri, err)
	}
	if _, err := sanitizeURI("::bad::"); err == nil {
		t.Error("expected error for bad uri")
	}
}

func TestMCPClientSendAndRead(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	c := NewMCPClientWithIO(in, out, io.Discard)
	req := &Request{Method: "ping"}
	if err := c.SendRequest(req); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected output to be written")
	}
	respJSON := `{"result":"pong"}`
	in.WriteString(respJSON + "\n")
	resp, err := c.ReadResponse()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if resp.Result != "pong" {
		t.Errorf("unexpected result %v", resp.Result)
	}
}

func TestMCPClientCallHelpers(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	c := NewMCPClientWithIO(in, out, errBuf)

	in.WriteString(`{"result":"ok"}` + "\n")
	resp, err := c.Call("test/method", nil, 1)
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if resp.Result != "ok" {
		t.Errorf("unexpected result %v", resp.Result)
	}
	if !bytes.Contains(out.Bytes(), []byte("test/method")) {
		t.Errorf("request not sent correctly: %s", out.String())
	}

	out.Reset()
	in.Reset()
	in.WriteString(`{"result":null}` + "\n")
	if _, err := c.ListResources(2); err != nil {
		t.Fatalf("list resources failed: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("resources/list")) {
		t.Errorf("list resources wrong request: %s", out.String())
	}

	out.Reset()
	in.Reset()
	in.WriteString(`{"result":null}` + "\n")
	if _, err := c.ListTools(3); err != nil {
		t.Fatalf("list tools failed: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("tools/list")) {
		t.Errorf("list tools wrong request: %s", out.String())
	}

	out.Reset()
	in.Reset()
	in.WriteString(`{"result":null}` + "\n")
	if _, err := c.ReadResource("http://x", 4); err != nil {
		t.Fatalf("read resource failed: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("resources/read")) {
		t.Errorf("read resource wrong request: %s", out.String())
	}

	out.Reset()
	in.Reset()
	in.WriteString(`{"result":null}` + "\n")
	if _, err := c.CallTool("tool", nil, 5); err != nil {
		t.Fatalf("call tool failed: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("tools/call")) {
		t.Errorf("call tool wrong request: %s", out.String())
	}

	c.PrintError("err %d", 1)
	c.PrintInfo("info %d", 2)
	if !bytes.Contains(errBuf.Bytes(), []byte("err 1")) || !bytes.Contains(errBuf.Bytes(), []byte("info 2")) {
		t.Errorf("stderr not written correctly: %s", errBuf.String())
	}
}
func TestSanitizeURIEmpty(t *testing.T) {
	if _, err := sanitizeURI("   "); err == nil {
		t.Error("expected error for empty uri")
	}
}

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) { return 0, fmt.Errorf("write fail") }

type failReader struct{}

func (failReader) Read(p []byte) (int, error) { return 0, fmt.Errorf("read fail") }

func TestSendRequestError(t *testing.T) {
	c := NewMCPClientWithIO(bytes.NewBuffer(nil), failWriter{}, io.Discard)
	if err := c.SendRequest(&Request{Method: "x"}); err == nil {
		t.Error("expected write error")
	}
}

func TestReadResponseErrors(t *testing.T) {
	c := NewMCPClientWithIO(failReader{}, io.Discard, io.Discard)
	if _, err := c.ReadResponse(); err == nil || !strings.Contains(err.Error(), "failed to read response") {
		t.Fatalf("unexpected error %v", err)
	}

	c = NewMCPClientWithIO(bytes.NewBufferString("bad\n"), io.Discard, io.Discard)
	if _, err := c.ReadResponse(); err == nil || !strings.Contains(err.Error(), "failed to unmarshal response") {
		t.Fatalf("expected json error, got %v", err)
	}

	c = NewMCPClientWithIO(bytes.NewBuffer(nil), io.Discard, io.Discard)
	if _, err := c.ReadResponse(); err == nil || !strings.Contains(err.Error(), "no response") {
		t.Fatalf("expected no response error, got %v", err)
	}
}

func TestCallError(t *testing.T) {
	c := NewMCPClientWithIO(bytes.NewBuffer(nil), failWriter{}, io.Discard)
	if _, err := c.Call("m", nil, 1); err == nil || !strings.Contains(err.Error(), "write fail") {
		t.Fatalf("expected send error, got %v", err)
	}
}

func TestNewMCPConfigAndSetTransport(t *testing.T) {
	cfg := NewMCPConfig("app", "0.1.0", "desc", nil, nil)
	if cfg.License != "MIT" || cfg.Transport.Type != "stdio" {
		t.Fatalf("unexpected defaults %+v", cfg)
	}
	cfg.SetTransport("rest", map[string]interface{}{"port": 80})
	if cfg.Transport.Type != "rest" || cfg.Transport.Options["port"] != 80 {
		t.Fatalf("set transport failed %+v", cfg.Transport)
	}
}

func TestLineAndColumnValid(t *testing.T) {
	data := []byte("a\nbc\nd")
	l, c, err := lineAndColumn(data, 5)
	if err != nil || l != 3 || c != 1 {
		t.Fatalf("expected line 3 col 1 got %d %d err %v", l, c, err)
	}
}

func TestFormatJSONErrorDefault(t *testing.T) {
	err := FormatJSONError([]byte("{}"), errors.New("boom"), "context")
	if err == nil || !strings.Contains(err.Error(), "context") {
		t.Fatalf("unexpected error %v", err)
	}
}
