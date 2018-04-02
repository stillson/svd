package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	spiffe "github.com/spiffe/go-spiffe/spiffe"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"flag"
)

// Constants
const (
	// file containing spiffe URI to check against
	SVIDFILE = "./URI"
	// cert and key for the svd service
	CERTFILE = "./svd.crt"
	KEYFILE  = "./svd.key"
	// root and intermediate of the CA for verification
	ROOTCA = "./ROOT.crt"
	INTCA  = "./INT.crt"
	// adress to run the service on
	ADDRESS = ":8443"
)

// global variables
var (
	confSan string
	rootCP  *x509.CertPool
	intCP   *x509.CertPool
)

// Give a path to a file, verify the path and return
// the contents of the file
func loadCert(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	rv, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

// separate cert checking into it's own function for testing
func checkSvid(sans []string, cert *x509.Certificate) error {
	fmt.Printf("|%v ", cert.Subject)
	err := spiffe.MatchID(sans, cert)

	if err != nil {
		// check the easy part (the SAN)
		fmt.Printf("|no matching SPIFFE id: %v ", err)
		return err
	} else {
		// now, check the hard part (the cert)
		err = spiffe.VerifyCertificate(cert, intCP, rootCP)

		if err != nil {
			fmt.Printf("|bad cert: %v ", err)
			return err
		}
	}

	return nil
}

// http request handler
// return 200 if the client cert of request is signed by the proper ca
// and the SPIFFE URI in the client cert is the same as the one
// configured for this server. Otherwise, returns 401
func svidHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sans := []string{confSan}

	fmt.Printf("%v ", r.RemoteAddr)

	// match svid
	if len(r.TLS.PeerCertificates) != 0 {
		err := checkSvid(sans, r.TLS.PeerCertificates[0])
		if err != nil {
			fmt.Printf("|SVID invalid\n")
			http.Error(w, "Invalid", http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Printf("|Wrong number of peer certs\n")
		http.Error(w, "Invalid", http.StatusUnauthorized)
	}

	fmt.Printf("|SVID valid\n")
	fmt.Fprintf(w, "OK\n")
}

// Main
func main() {
	// command line arguments
	address := flag.String("addr", ADDRESS, "adress to run on")
	flag.Parse()

	// load config (i.e. SPIFFE ID)
	rvString, err := ioutil.ReadFile(SVIDFILE)
	if err != nil {
		log.Fatal("Unable to open spiffe id file: ", err)
	}

	confSan = strings.TrimSpace(string(rvString))

	// load server ssl certs
	absPathServerCrt, err := filepath.Abs(CERTFILE)
	if err != nil {
		log.Fatal(err, "Unable to open cert file: ")
	}

	absPathServerKey, err := filepath.Abs(KEYFILE)
	if err != nil {
		log.Fatal(err, "Unable to open key file: ")
	}

	// load CA certs

	rootCACert, err := loadCert(ROOTCA)
	if err != nil {
		log.Fatal(err, "Unable to load root cert")
	}
	rootCP = x509.NewCertPool()
	rootCP.AppendCertsFromPEM(rootCACert)

	intCACert, err := loadCert(INTCA)
	if err != nil {
		log.Fatal(err, "Unable to load intemdiate cert cert")
	}
	intCP = x509.NewCertPool()
	intCP.AppendCertsFromPEM(intCACert)

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAnyClientCert,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	http.HandleFunc("/", svidHandler)

	httpServer := &http.Server{
		Addr:      *address,
		TLSConfig: tlsConfig,
	}

	err = httpServer.ListenAndServeTLS(absPathServerCrt, absPathServerKey) // set listen port
	if err != nil {
		log.Fatal(err, "ListenAndServe: ")
	}

}
