package tales_of_ssl

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"strconv"
	"time"
)

type TalesOfSsl struct{}

type Data struct {
	PrivateKey   string `json:"private_key"`
	RequiredData struct {
		Domain       string `json:"domain"`
		SerialNumber string `json:"serial_number"`
		Country      string `json:"country"`
	} `json:"required_data"`
}

type Output struct {
	Certificate string `json:"certificate"`
}

func (d TalesOfSsl) Solve(input string) (interface{}, error) {
	data := new(Data)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	value, err := strconv.ParseInt(data.RequiredData.SerialNumber[2:], 16, 64)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(value),
		Subject: pkix.Name{
			Country:    []string{"CC"},
			CommonName: "test",
		},

		DNSNames: []string{data.RequiredData.Domain},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	block, _ := pem.Decode([]byte("-----BEGIN RSA PRIVATE KEY-----\n" + data.PrivateKey + "\n-----END RSA PRIVATE KEY-----"))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	return &Output{
		Certificate: b64.StdEncoding.EncodeToString(derBytes),
	}, nil
}
