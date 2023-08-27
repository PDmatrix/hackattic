package password_hashing

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"sync"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

type PasswordHashing struct{}

type Data struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Pbkdf2   struct {
		Rounds int    `json:"rounds"`
		Hash   string `json:"hash"`
	} `json:"pbkdf2"`
	Scrypt struct {
		N       int    `json:"N"`
		R       int    `json:"r"`
		P       int    `json:"p"`
		Buflen  int    `json:"buflen"`
		Control string `json:"_control"`
	} `json:"scrypt"`
}

type Output struct {
	Sha256 string `json:"sha256"`
	Hmac   string `json:"hmac"`
	Pbkdf2 string `json:"pbkdf2"`
	Scrypt string `json:"scrypt"`
}

func (d PasswordHashing) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	wg := sync.WaitGroup{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	wg.Add(4)

	go getSha256(&data.Password, output, &wg)
	go getHmac(&data.Password, &data.Salt, output, &wg)
	go getPbkdf2(&data.Password, &data.Salt, data.Pbkdf2.Rounds, 32, data.Pbkdf2.Hash, output, &wg)
	go getScrypt(&data.Password, &data.Salt, data.Scrypt.N, data.Scrypt.R, data.Scrypt.P, data.Scrypt.Buflen, output, &wg)

	wg.Wait()
	return output, nil
}

func getSha256(password *string, output *Output, wg *sync.WaitGroup) {
	output.Sha256 = fmt.Sprintf("%x", sha256.Sum256([]byte(*password)))
	wg.Done()
}

func getPbkdf2(password *string, salt *string, rounds int, keylen int, hashName string, output *Output, wg *sync.WaitGroup) {
	saltDecoded, _ := b64.StdEncoding.DecodeString(*salt)

	var hashFunc func() hash.Hash
	switch hashName {
	case "sha256":
		hashFunc = sha256.New
	case "sha1":
		hashFunc = sha1.New
	}

	output.Pbkdf2 = fmt.Sprintf("%x", pbkdf2.Key([]byte(*password), saltDecoded, rounds, keylen, hashFunc))
	wg.Done()
}

func getScrypt(password *string, salt *string, n int, r int, p int, buflen int, output *Output, wg *sync.WaitGroup) {
	saltDecoded, _ := b64.StdEncoding.DecodeString(*salt)

	scryptEncrypted, _ := scrypt.Key([]byte(*password), saltDecoded, n, r, p, buflen)

	output.Scrypt = fmt.Sprintf("%x", scryptEncrypted)
	wg.Done()
}

func getHmac(password *string, salt *string, output *Output, wg *sync.WaitGroup) {
	saltDecoded, _ := b64.StdEncoding.DecodeString(*salt)

	h := hmac.New(sha256.New, []byte(saltDecoded))

	h.Write([]byte(*password))

	output.Hmac = fmt.Sprintf("%x", h.Sum(nil))
	wg.Done()
}
