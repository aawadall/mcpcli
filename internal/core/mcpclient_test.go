package core

import (
	"bytes"
	"io"
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
