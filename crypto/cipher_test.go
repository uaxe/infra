package crypto

import (
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var (
	letters            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	timeoutInOperation = 3 * time.Second
)

func RandStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func RandLowStr(n int) string {
	return strings.ToLower(RandStr(n))
}

func TestAesCtr(t *testing.T) {
	var cipherData CipherData
	cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCtr(cipherData)

	byteReader := strings.NewReader(RandLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	encrypter.Close()
	buff := make([]byte, 10)
	n, err := encrypter.Read(buff)
	if n != 0 {
		t.Fatal("not read empty")
		return
	}
	if err != io.EOF {
		t.Fatal("not read empty")
		return
	}

	deReader := cipher.Encrypt(byteReader)
	decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	decrypter.Close()
	buff = make([]byte, 10)
	n, err = decrypter.Read(buff)
	if n != 0 {
		t.Fatal("not read empty")
		return
	}
	if err != io.EOF {
		t.Fatal("not read empty")
		return
	}
}

func TestAesCfb(t *testing.T) {
	var cipherData CipherData
	cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCfb(cipherData)

	byteReader := strings.NewReader(RandLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	encrypter.Close()
	buff := make([]byte, 10)
	n, err := encrypter.Read(buff)
	if n != 0 {
		t.Fatal("not read empty")
		return
	}
	if err != io.EOF {
		t.Fatal("not read empty")
		return
	}

	deReader := cipher.Encrypt(byteReader)
	decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	decrypter.Close()
	buff = make([]byte, 10)
	n, err = decrypter.Read(buff)
	if n != 0 {
		t.Fatal("not read empty")
		return
	}
	if err != io.EOF {
		t.Fatal("not read empty")
		return
	}
}
