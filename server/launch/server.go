package launch

import (
	"fmt"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/mux"
	"log"
)

type ServerPorts struct {
	OldUdpPort  int
	NewUdpPort  int
	OldDtlsPort int
	NewDtlsPort int
}

type ServerFlags struct {
	Address  string
	Password string
	Dtls     bool
	Network  string
}

func Server(flagValues ServerFlags, udpPort int, dtlsPort int, r *mux.Router) {
	udpAddr := fmt.Sprintf("%s:%d", flagValues.Address, udpPort)
	dtlsAddr := fmt.Sprintf("%s:%d", flagValues.Address, dtlsPort)

	if flagValues.Dtls {
		log.Printf("dTLS %s Server listening on: %s\n", flagValues.Network, dtlsAddr)
		log.Printf("dTLS PSK: %s\n", flagValues.Password)
		log.Fatal(coap.ListenAndServeDTLS(flagValues.Network, dtlsAddr, &piondtls.Config{
			PSK: func(hint []byte) ([]byte, error) {
				log.Printf("Client's hint: %s \n", hint)
				return []byte(flagValues.Password), nil
			},
			PSKIdentityHint: []byte("Pion DTLS Client"),
			CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
		}, r))
	} else {
		log.Printf("%s Server listening on: %s\n", flagValues.Network, udpAddr)
		log.Fatal(coap.ListenAndServe(flagValues.Network, udpAddr, r))
	}
}
