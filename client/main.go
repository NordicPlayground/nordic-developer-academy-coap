package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"
	"flag"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/go-coap/v3/udp/client"
)

func main() {
	address := flag.String("address", "localhost",
		"The UDP Server listen address, e.g. `localhost` or `0.0.0.0`.")
	flag.Parse()
	udpAddr := fmt.Sprintf("%s:%d", *address, 5688)
	dtlsAddr := fmt.Sprintf("%s:%d", *address, 5689)
	fmt.Printf("UDP Server listening on: %s\n", udpAddr)
	fmt.Printf("dTLS Server listening on: %s\n", dtlsAddr)

	co, err := udp.Dial(udpAddr)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	testHello(co, ctx)
	testPutCustom(co, ctx, "/my-random-id")
	testGetCustom(co, ctx, "/my-random-id")
}

// Test Hello resource
func testHello(co *client.Conn, ctx context.Context) {
	resp, err := co.Get(ctx, "/static/hello")
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	data, err := io.ReadAll(resp.Body())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

// Test writing to custom resource
func testPutCustom(co *client.Conn, ctx context.Context, path string) {
	_, err := co.Put(ctx, path, message.TextPlain, bytes.NewReader([]byte("Some random value.")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Written to %s\n", path)
}

// Test getting custom resource
func testGetCustom(co *client.Conn, ctx context.Context, path string) []byte {
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	data, err := io.ReadAll(resp.Body())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
	return data
}