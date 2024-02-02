package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"coap-server/internal"
	piondtls "github.com/pion/dtls/v2"
	coap "github.com/plgd-dev/go-coap/v3"


	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	address := flag.String("address", "localhost",
		"The UDP Server listen address, e.g. `localhost` or `0.0.0.0`.")
	password := flag.String("password", "connect:anything",
		"The password to use for the PSK in dTLS.")
	dtls := flag.Bool("dTLS", false, "Start a dTLS server")
	flag.Parse()

	storageConnectionString, ok := os.LookupEnv("STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'STORAGE_CONNECTION_STRING' could not be found")
	}
	storageClient, err := azblob.NewClientFromConnectionString(storageConnectionString, nil)
	check(err)
	log.Printf("Azure Storage client URL: %s\n", storageClient.URL())

	containerName, ok := os.LookupEnv("STORAGE_CONTAINER_NAME")
	if !ok {
		log.Fatal("the environment variable 'STORAGE_CONTAINER_NAME' could not be found")
	}
	log.Printf("Container name: %s\n", containerName)

	r := internal.NewServer(storageClient, containerName)

	if *dtls {
		dtlsAddr := fmt.Sprintf("%s:%d", *address, 5689)
		log.Printf("dTLS UDP Server listening on: %s\n", dtlsAddr)
		log.Printf("dTLS PSK: %s\n", *password)
		log.Fatal(coap.ListenAndServeDTLS("udp", dtlsAddr, &piondtls.Config{
			PSK: func(hint []byte) ([]byte, error) {
				log.Printf("Client's hint: %s \n", hint)
				return []byte(*password), nil
			},
			PSKIdentityHint: []byte("Pion DTLS Client"),
			CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
		}, r))
	} else {
		udpAddr := fmt.Sprintf("%s:%d", *address, 5688)
		log.Printf("UDP Server listening on: %s\n", udpAddr)
		log.Fatal(coap.ListenAndServe("udp", udpAddr, r))
	}
}
