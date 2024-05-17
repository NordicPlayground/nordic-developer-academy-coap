package main

import (
	"coap-server/launch"
	"flag"
	"log"
	"os"

	"coap-server/internal"

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
	network := flag.String("network", "udp4",
		"The network to use, `udp4` or `udp6`.")
	flag.Parse()

	var flagValues = launch.ServerFlags{
		Address:  *address,
		Password: *password,
		Dtls:     *dtls,
		Network:  *network,
	}

	var ports = launch.ServerPorts{
		OldUdpPort:  5688,
		NewUdpPort:  5683,
		OldDtlsPort: 5689,
		NewDtlsPort: 5684,
	}

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

	// old ports
	go launch.Server(flagValues, ports.OldUdpPort, ports.OldDtlsPort, r)

	// new ports
	launch.Server(flagValues, ports.NewUdpPort, ports.NewDtlsPort, r)

}
