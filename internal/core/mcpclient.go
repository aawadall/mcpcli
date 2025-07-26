package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

type MCPClient struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewMCPClient() *MCPClient {
	return &MCPClient{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func NewMCPClientWithIO(stdin io.Reader, stdout io.Writer, stderr io.Writer) *MCPClient {
	return &MCPClient{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *MCPClient) SendRequest(request *Request) error {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = fmt.Fprintln(c.stdout, string(requestJSON))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	return nil
}

func (c *MCPClient) ReadResponse() (*Response, error) {
	scanner := bufio.NewScanner(c.stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		return nil, fmt.Errorf("no response received")
	}
	var response Response
	if err := json.Unmarshal([]byte(scanner.Text()), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &response, nil
}

func (c *MCPClient) Call(method string, params map[string]interface{}, id interface{}) (*Response, error) {
	req := &Request{
		Method: method,
		Params: params,
		ID:     id,
	}
	if err := c.SendRequest(req); err != nil {
		return nil, err
	}
	return c.ReadResponse()
}

func (c *MCPClient) ListResources(id interface{}) (*Response, error) {
	return c.Call("resources/list", nil, id)
}

func sanitizeURI(uri string) (string, error) {
	trimmed := strings.TrimSpace(uri)
	if trimmed == "" {
		return "", fmt.Errorf("uri is empty")
	}
	if _, err := url.ParseRequestURI(trimmed); err != nil {
		return "", fmt.Errorf("invalid uri: %w", err)
	}
	return trimmed, nil
}

func (c *MCPClient) ReadResource(uri string, id interface{}) (*Response, error) {
	sanitized, err := sanitizeURI(uri)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{"uri": sanitized}
	return c.Call("resources/read", params, id)
}

func (c *MCPClient) ListTools(id interface{}) (*Response, error) {
	return c.Call("tools/list", nil, id)
}

func (c *MCPClient) CallTool(name string, arguments map[string]interface{}, id interface{}) (*Response, error) {
	params := map[string]interface{}{"name": name, "arguments": arguments}
	return c.Call("tools/call", params, id)
}

func (c *MCPClient) PrintError(format string, args ...interface{}) {
	fmt.Fprintf(c.stderr, format+"\n", args...)
}

func (c *MCPClient) PrintInfo(format string, args ...interface{}) {
	fmt.Fprintf(c.stderr, format+"\n", args...)
}
