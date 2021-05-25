package certutil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/twmb/algoimpl/go/graph"
	"golang.org/x/crypto/pkcs12"
	"strings"
)

// Takes Base64 Encoded PKCS12 as String and produces
func PemPrivateKeyFromPkcs12(b64pkcs12 string) string {
	p12, _ := base64.StdEncoding.DecodeString(b64pkcs12)

	// Get the PEM Blocks
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		panic(err)
	}

	return findPrivateKeyInPemBlocks(blocks)
}

func PemPrivateKeyFromPem(data string) string {
	// Convert string to Pem Blocks
	blocks := stringToPemBlocks(data)
	// Find the Private Key from these blocks
	return findPrivateKeyInPemBlocks(blocks)
}

func PemCertFromPkcs12(b64pkcs12 string) string {
	p12, _ := base64.StdEncoding.DecodeString(b64pkcs12)

	// Get the PEM Blocks
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		panic(err)
	}

	return findLeafCertInPemBlocks(blocks)
}

func PemCertFromPem(data string) string {
	//Convert string to pem blocks
	blocks := stringToPemBlocks(data)
	//Find leaf cert and return
	return findLeafCertInPemBlocks(blocks)
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

	return findChainInPemBlocks(blocks, justIssuers)
}

func PemChainFromPem(data string, justIssuers bool) string {
	// Get the PEM blocks from the string
	blocks := stringToPemBlocks(data)

	// Sort and return the chain
	return findChainInPemBlocks(blocks, justIssuers)
}

func SortedChain(certs []*x509.Certificate, justIssuers bool) []x509.Certificate {
	g := graph.New(graph.Directed)

	// Make a graph where each node represents a certificate and the key is its subject key identifier
	certGraph := make(map[string]graph.Node, 0)

	// Construct each certificate in the chain into a full certificate object
	for _, cert := range certs {
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
	var sortedCerts []x509.Certificate

	for i := range sorted {
		cert := (*sorted[i].Value).(x509.Certificate)
		sortedCerts = append(sortedCerts, cert)
	}

	if justIssuers {
		// If we only have the leaf cert there are no issuers to return
		if len(sortedCerts) == 1 {
			return nil
		} else {
			return sortedCerts[1:]
		}
	}

	return sortedCerts
}

func stringToPemBlocks(data string) []*pem.Block {
	// Build an array of pem.Block
	var blocks []*pem.Block
	rest := []byte(data)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
	}
	return blocks
}

func findPrivateKeyInPemBlocks(blocks []*pem.Block ) string {
	var keyBuffer bytes.Buffer
	//Find the private key from all the blocks
	for _, block := range blocks {
		// Private Key?
		if block.Type == "PRIVATE KEY" || strings.HasSuffix(block.Type, " PRIVATE KEY") {
			key, err := parsePrivateKey(block.Bytes)
			if err != nil {
				panic(err)
			}

			// Force it to pkcs8 for consistency
			privBytes, err := x509.MarshalPKCS8PrivateKey(key)
			if err != nil {
				panic(err)
			}

			if err := pem.Encode(&keyBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
				panic(fmt.Sprintf("Failed to write data: %s", err))
			}
			break
		}
	}
	return keyBuffer.String()
}

// https://golang.org/src/crypto/tls/tls.go?#L370
func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
		}
	}

	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("tls: failed to parse private key")
}

func findLeafCertInPemBlocks(blocks []*pem.Block) string {
	var certs []*x509.Certificate
	//Find all the Certificate blocks
	for _, block := range blocks {
		// Private Key?
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)

			if err != nil {
				panic(err)
			}

			certs = append(certs, cert)
		}
	}

	// Sort the certs
	sortedCerts := SortedChain(certs, false)

	// PEM Encode first cert in sortedCerts
	var certBuffer bytes.Buffer
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: sortedCerts[0].Raw}); err != nil {
		panic(fmt.Sprintf("Failed to write data: %s", err))
	}

	return certBuffer.String()
}

func findChainInPemBlocks(blocks []*pem.Block, justIssuers bool) string {
	var certs []*x509.Certificate
	//Find all the Certificate blocks
	for _, block := range blocks {
		// Private Key?
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)

			if err != nil {
				panic(err)
			}

			certs = append(certs, cert)
		}
	}

	// Sort the certs
	sortedCerts := SortedChain(certs, justIssuers)

	// PEM Encode all the certs
	var certBuffer bytes.Buffer
	for i := range sortedCerts {
		if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: sortedCerts[i].Raw}); err != nil {
			panic(fmt.Sprintf("Failed to write data: %s", err))
		}
	}

	return certBuffer.String()
}