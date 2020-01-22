package templaterenderer

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/Masterminds/sprig"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"log"
	"text/template"
	"github.com/twmb/algoimpl/go/graph"
)

func RenderFile(path string, resourceMap resource.ResourceMap) string {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		log.Panicf("Error reading template %v: %v", path, err)
	}

	return RenderInline(string(contents), resourceMap)
}

func RenderInline(templateContents string, resourceMap resource.ResourceMap) string {
	helpers := template.FuncMap{
		"privateKey": func(name string) interface{} {
			value, ok := resourceMap.Secrets[name]
			privateKey := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					privateKey = pemPrivateKeyFromPem(*value.Value)
				case "application/x-pkcs12":
					privateKey = pemPrivateKeyFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("privateKey lookup failed: Expected a Secret with name %v\n", name)
			}
			return privateKey
		},
		"cert": func(name string) interface{} {
			// TODO: If the cert can be found on either a Cert or a Secret, we should handle discovering it from both
			value, ok := resourceMap.Secrets[name]
			cert := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					cert = pemCertFromPem(*value.Value)
				case "application/x-pkcs12":
					cert = pemCertFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("cert lookup failed: Expected a Secret with name %v\n", name)
			}
			return cert
		},
		"chain": func(name string) interface{} {
			value, ok := resourceMap.Secrets[name]
			chain := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					chain = pemChainFromPem(*value.Value)
				case "application/x-pkcs12":
					chain = pemChainFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("cert lookup failed: Expected a Secret with name %v\n", name)
			}
			return chain
		},
	}

	// Init the template
	t, err := template.New("template").Funcs(helpers).Funcs(sprig.TxtFuncMap()).Parse(templateContents)
	if err != nil {
		//log.Panicf("Error parsing template:\n%v\nError:\n%v\n", templateContents, err)
	}

	// Execute the template
	var buf bytes.Buffer
	err = t.Execute(&buf, resourceMap)
	if err != nil {
		log.Panicf("Error executing template:\n%v\nResources:\n%v\nError:\n%v\n", templateContents, resourceMap, err)
	}

	result := buf.String()

	return result
}

func pemPrivateKeyFromPkcs12(b64pkcs12 string) string {
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

	return pemPrivateKeyFromPem(string(pemData))
}

func pemPrivateKeyFromPem(data string) string {
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
		log.Panicf("Unable to marshal private key: %v", err)
	}

	// Encode just the private key back to PEM and return it
	var privPem bytes.Buffer
	if err := pem.Encode(&privPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Panicf("Failed to write data: %s", err)
	}

	return privPem.String()
}

func pemCertFromPkcs12(b64pkcs12 string) string {
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

	return pemCertFromPem(string(pemData))
}

func pemCertFromPem(data string) string {
	pemBytes := []byte(data)

	// Use tls lib to construct tls certificate and key object from PEM data
	// The tls.X509KeyPair function is smart enough to parse combined cert and key pem data
	certAndKey, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		log.Panicf("Error generating X509KeyPair: %v", err)
	}

	leaf, err := x509.ParseCertificate(certAndKey.Certificate[0])
	if err != nil {
		log.Panic(err)
	}

	// Encode just the leaf cert as pem
	var certPem bytes.Buffer
	if err := pem.Encode(&certPem, &pem.Block{Type: "CERTIFICATE", Bytes: leaf.Raw}); err != nil {
		log.Panicf("Failed to write data: %s", err)
	}

	return certPem.String()
}

func pemChainFromPkcs12(b64pkcs12 string) string {
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

	return pemChainFromPem(string(pemData))
}

func pemChainFromPem(data string) string {
	pemBytes := []byte(data)

	// Use tls lib to construct tls certificate and key object from PEM data
	// The tls.X509KeyPair function is smart enough to parse combined cert and key pem data
	certAndKey, err := tls.X509KeyPair(pemBytes, pemBytes)
	if err != nil {
		log.Panicf("Error generating X509KeyPair: %v", err)
	}

	// The chain is the rest of the certs
	/*
	var chain bytes.Buffer
	for _, issuerBytes := range certAndKey.Certificate[1:] {
		if err := pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: issuerBytes}); err != nil {
			log.Panicf("Failed to write data: %s", err)
		}
	}
	 */

	//TODO make sure the chain is in the right order - each cert certifies the one preceding it
	return sortedChain(certAndKey.Certificate)
}

func sortedChain(rawChain [][]byte) string {
	g := graph.New(graph.Directed)

	// Make a graph where each node represents a certificate and the key is its subject key identifier
	certDict :=  make(map[string]x509.Certificate)
	certGraph := make(map[string]graph.Node, 0)

	// Construct each certificate in the chain into a full certificate object
	for _, certBytes := range rawChain {
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			log.Panic("Unable to parse certificate chain")
		}
		certGraph[string(cert.SubjectKeyId)] = g.MakeNode()
		*certGraph[string(cert.SubjectKeyId)].Value = cert.SubjectKeyId
		certDict[string(cert.SubjectKeyId)] = *cert
	}

	// Make the edges of the graph from child cert to issuer
	for _, cert := range certDict {
		g.MakeEdge(certGraph[string(cert.SubjectKeyId)], certGraph[string(cert.AuthorityKeyId)])
	}

	sorted := g.TopologicalSort()

	// Construct the sorted chain PEM block
	var chain bytes.Buffer
	for i := range sorted {
		// Convert the value to a string so we can use it for lookups
		//subjectKeyId := fmt.Sprintf("%v", *sorted[i].Value)
		subjectKeyIdBytes := (*sorted[i].Value).([]byte)
		subjectKeyId := string(subjectKeyIdBytes)
		if err := pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: certDict[subjectKeyId].Raw}); err != nil {
			log.Panicf("Failed to write data: %s", err)
		}
	}

	log.Print(chain.String())
	return chain.String()
}