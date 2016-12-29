package alexa

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrInvalidCertificateURL is returned when the certificate URL is invalid
var ErrInvalidCertificateURL = errors.New("invalid certificate URL")

// ErrFailedDecodingPEM is returned when PEM decoding fails
var ErrFailedDecodingPEM = errors.New("failed decoding PEM")

//ErrAmazonCertificateExpired is returned when the AWS certificate has expired
var ErrAmazonCertificateExpired = errors.New("amazon certificate expired")

// ErrNameNotFound is returned when the domain isn't found on the certificate
var ErrNameNotFound = errors.New("amazon certificate name not found")

// ErrFailedToVerifyBody is returned when verification of the request body failed
var ErrFailedToVerifyBody = errors.New("failed to verify request body")

// ErrInvalidSignature is returned when the request signature is invalid
var ErrInvalidSignature = errors.New("signature invalid")

// VerifyRequest verifies the request was from AWS
func VerifyRequest(r *http.Request) (err error) {
	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		return ErrInvalidCertificateURL
	}

	// Fetch certificate data
	cert, err := http.Get(certURL)
	if err != nil {
		return fmt.Errorf("could not fetch Amazon certificate: %s", err.Error())
	}

	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err.Error())
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		return ErrFailedDecodingPEM
	}

	x509cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("could not parse certificate: %s", err.Error())
	}

	// Check the certificate date
	if time.Now().Unix() < x509cert.NotBefore.Unix() || time.Now().Unix() > x509cert.NotAfter.Unix() {
		return ErrAmazonCertificateExpired
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range x509cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		return ErrNameNotFound
	}

	// Verify the key
	publicKey := x509cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err = io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
	if err != nil {
		return ErrFailedToVerifyBody
	}

	r.Body = ioutil.NopCloser(&bodyBuf)

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil {
		return ErrInvalidSignature
	}

	return nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}
