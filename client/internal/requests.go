package internal

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp/client"
)

// Test Hello resource
func TestHello(co *client.Conn, ctx context.Context) {
	resp, err := co.Get(ctx, "/static/hello")
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	data, err := io.ReadAll(resp.Body())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", data)
}

// Test writing to custom resource
func TestPutCustom(co *client.Conn, ctx context.Context, path string, value []byte) {
	_, err := co.Put(ctx, path, message.TextPlain, bytes.NewReader(value))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Written to %s\n", path)
}

// Test getting custom resource
func TestGetCustom(co *client.Conn, ctx context.Context, path string) ([]byte, error) {
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return nil, err
	}
	data, err := io.ReadAll(resp.Body())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("%s\n", data)
	return data, nil
}
