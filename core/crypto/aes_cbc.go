package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type aesCbc struct {
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
}

func newAesCbc(cd CipherData) (Cipher, error) {
	block, err := aes.NewCipher(cd.Key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(block, cd.IV)
	decrypter := cipher.NewCBCDecrypter(block, cd.IV)
	return &aesCbc{encrypter, decrypter}, nil
}

func (c *aesCbc) Encrypt(src io.Reader) io.Reader {
	reader := &cbcEncryptReader{
		encrypter: c.encrypter,
		src:       src,
	}
	return reader
}

type cbcEncryptReader struct {
	encrypter cipher.BlockMode
	src       io.Reader
}

func (reader *cbcEncryptReader) Read(data []byte) (int, error) {
	plainText := make([]byte, len(data), len(data))
	n, err := reader.src.Read(plainText)
	if n > 0 {
		plainText = plainText[0:n]
		reader.encrypter.CryptBlocks(data, plainText)
	}
	return n, err
}

func (c *aesCbc) Decrypt(src io.Reader) io.Reader {
	return &cbcDecryptReader{
		decrypter: c.decrypter,
		src:       src,
	}
}

type cbcDecryptReader struct {
	decrypter cipher.BlockMode
	src       io.Reader
}

func (reader *cbcDecryptReader) Read(data []byte) (int, error) {
	cryptoText := make([]byte, len(data), len(data))
	n, err := reader.src.Read(cryptoText)
	if n > 0 {
		cryptoText = cryptoText[0:n]
		reader.decrypter.CryptBlocks(data, cryptoText)
	}
	return n, err
}
