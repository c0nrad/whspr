package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
)

var TrustedCertPath = "./certs/trusted"
var NameKeyPairs map[string]*rsa.PublicKey

func parseTrustedCerts() {
	NameKeyPairs = make(map[string]*rsa.PublicKey)

	directory, err := ioutil.ReadDir(TrustedCertPath)
	if err != nil {
		log.Fatalln("Error reading directory: ", TrustedCertPath, err)
	}

	for _, file := range directory {

		name := file.Name()
		commonName, publicKey := parseNameAndPublicKey(TrustedCertPath + "/" + name)

		log.Println("Found entry for", commonName)
		NameKeyPairs[commonName] = publicKey
	}
}

func getPublicKey(name []byte) *rsa.PublicKey {
	return NameKeyPairs[string(name)]
}

func parseNameAndPublicKey(filename string) (string, *rsa.PublicKey) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error reading file: ", filename, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatalln("Error parsing pem")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalln("Failed to parse certificates: ", err)
	}

	if cert.PublicKeyAlgorithm != x509.RSA {
		log.Fatalln("Certificate public key is not rsa", filename)
	}

	return cert.Subject.CommonName, cert.PublicKey.(*rsa.PublicKey)
}

func parsePrivateKey(filename string) *rsa.PrivateKey {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}

	return priv
}

// // First, create the set of root certificates. For this example we only
//   // have one. It's also possible to omit this in order to use the
//   // default root set of the current operating system.
//   roots := x509.NewCertPool()
//   ok := roots.AppendCertsFromPEM([]byte(rootPEM))
//   if !ok {
//     panic("failed to parse root certificate")
//   }

//   block, _ := pem.Decode([]byte(certPEM))
//   if block == nil {
//     panic("failed to parse certificate PEM")
//   }
//   cert, err := x509.ParseCertificate(block.Bytes)
//   if err != nil {
//     panic("failed to parse certificate: " + err.Error())
//   }

//   opts := x509.VerifyOptions{
//     DNSName: "mail.google.com",
//     Roots:   roots,
//   }

//   if _, err := cert.Verify(opts); err != nil {
//     panic("failed to verify certificate: " + err.Error())
//   }
