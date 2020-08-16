package certutil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/twmb/algoimpl/go/graph"
	"golang.org/x/crypto/pkcs12"
	log "github.com/sirupsen/logrus"
)

func PemPrivateKeyFromPkcs12(b64pkcs12 string) string {
	p12, _ := base64.StdEncoding.DecodeString(b64pkcs12)

	// Get the PEM Blocks
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		panic(err)
	}

	// Append all PEM Blocks together
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	return PemPrivateKeyFromPem(string(pemData))
}

func PemPrivateKeyFromPem(data string) string {
	pemBytes := []byte(data)

	// Use tls lib to construct tls certificate and key object from PEM data
	// The tls.X509KeyPair function is smart enough to parse combined cert and key pem data
	certAndKey, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		panic(err)
	}

	// Get parsed private key as PKCS8 data
	privBytes, err := x509.MarshalPKCS8PrivateKey(certAndKey.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("Unable to marshal private key: %v", err))
	}

	// Encode just the private key back to PEM and return it
	var privPem bytes.Buffer
	if err := pem.Encode(&privPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		panic(fmt.Sprintf("Failed to write data: %s", err))
	}

	return privPem.String()
}

func PemCertFromPkcs12(b64pkcs12 string) string {
	p12, _ := base64.StdEncoding.DecodeString(b64pkcs12)

	// Get the PEM Blocks
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		panic(err)
	}

	// Append all PEM Blocks together
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	return PemCertFromPem(string(pemData))
}

func PemCertFromPem(data string) string {
	pemBytes := []byte(data)

	// Use tls lib to construct tls certificate and key object from PEM data
	// The tls.X509KeyPair function is smart enough to parse combined cert and key pem data
	certAndKey, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		panic(fmt.Sprintf("Error generating X509KeyPair: %v", err))
	}

	leaf, err := x509.ParseCertificate(certAndKey.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Encode just the leaf cert as pem
	var certPem bytes.Buffer
	if err := pem.Encode(&certPem, &pem.Block{Type: "CERTIFICATE", Bytes: leaf.Raw}); err != nil {
		panic(fmt.Sprintf("Failed to write data: %s", err))
	}

	return certPem.String()
}

func PemCertFromBytes(derBytes []byte) string {
	// Encode just the leaf cert as pem
	var certPem bytes.Buffer
	if err := pem.Encode(&certPem, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		panic(fmt.Sprintf("Failed to write data: %s", err))
	}

	return certPem.String()
}

func PemChainFromPkcs12(b64pkcs12 string, justIssuers bool) string {
	p12, _ := base64.StdEncoding.DecodeString(b64pkcs12)

	// Get the PEM Blocks
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		panic(err)
	}

	// Append all PEM Blocks together
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	return PemChainFromPem(string(pemData), justIssuers)
}

func PemChainFromPem(data string, justIssuers bool) string {
	pemBytes := []byte(data)

	// Use tls lib to construct tls certificate and key object from PEM data
	// The tls.X509KeyPair function is smart enough to parse combined cert and key pem data
	certAndKey, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		panic(fmt.Sprintf("Error generating X509KeyPair: %v", err))
	}

	return SortedChain(certAndKey.Certificate, justIssuers)
}

func PemChainFromBytes(derBytes []byte, justIssuers bool) string {
	certs, err := x509.ParseCertificates(derBytes)
	if err != nil {
		panic(fmt.Sprintf("Error parsing Certificate: %v", err))
	}

	var rawCerts [][]byte
	for _, c := range certs {
		rawCerts = append(rawCerts, c.Raw)
	}

	return SortedChain(rawCerts, justIssuers)
}

func SortedChain(rawChain [][]byte, justIssuers bool) string {
	g := graph.New(graph.Directed)

	// Make a graph where each node represents a certificate and the key is its subject key identifier
	certGraph := make(map[string]graph.Node, 0)

	// Construct each certificate in the chain into a full certificate object
	for _, certBytes := range rawChain {
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			panic("Unable to parse certificate chain")
		}
		certGraph[string(cert.SubjectKeyId)] = g.MakeNode()
		*certGraph[string(cert.SubjectKeyId)].Value = *cert
	}

	// Make the edges of the graph from child cert to issuer
	for _, node := range certGraph {
		cert := (*node.Value).(x509.Certificate)
		g.MakeEdge(certGraph[string(cert.SubjectKeyId)], certGraph[string(cert.AuthorityKeyId)])
	}

	// Sort the graph
	sorted := g.TopologicalSort()

	// If sorted only has one element that must be the leaf and we have no chain to return
	if len(sorted) == 1 {
		log.Print("No chain detected in input")
		return ""
	}

	// Construct the sorted chain PEM block
	var chainPem bytes.Buffer

	// If sorted len is greater than 1 we have a chain to parse
	// Check if we want just the issuers or the full chain
	issuers := sorted
	if justIssuers {
		issuers = sorted[1:]
	}

	for i := range issuers {
		if err := pem.Encode(&chainPem, &pem.Block{Type: "CERTIFICATE", Bytes: (*issuers[i].Value).(x509.Certificate).Raw}); err != nil {
			panic(fmt.Sprintf("Failed to write data: %s", err))
		}
	}

	return chainPem.String()
}