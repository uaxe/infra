package crypto

import (
	"errors"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestAesCtr(t *testing.T) {
	var cipherData CipherData
	_ = cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCtr(cipherData)

	byteReader := strings.NewReader(RandLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	_ = encrypter.Close()
	buff := make([]byte, 10)
	_, err := encrypter.Read(buff)
	if !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}

	deReader := cipher.Encrypt(byteReader)
	decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	_ = decrypter.Close()
	buff = make([]byte, 10)
	_, err = decrypter.Read(buff)
	if !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
}

func RandLowStr(n int) string {
	return strings.ToLower(RandStr(n))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}
