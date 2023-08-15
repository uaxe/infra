package crypto

import (
	"fmt"
	"io"
)

const (
	ivSize = 16

	// AES-128, AES-192, or AES-256.
	// either 16, 24, or 32 bytes to select
	AES_TYPE_128 = AesType(16)
	AES_TYPE_192 = AesType(24)
	AES_TYPE_256 = AesType(32)
)

type AesType int

type aesCfbCipherBuilder struct {
	MasterCipher MasterCipher
	aesKeySize   int
}

type aesCfbCipher struct {
	CipherData CipherData
	Cipher     Cipher
}

func CreateAesCfbCipher(cipher MasterCipher, aesType AesType) ContentCipherBuilder {
	switch aesType {
	default:
		panic(fmt.Sprintf("aes type unknown %d", aesType))
	case AES_TYPE_128, AES_TYPE_192, AES_TYPE_256:
		break
	}
	return aesCfbCipherBuilder{MasterCipher: cipher, aesKeySize: int(aesType)}
}

// createCipherData create CipherData for encrypt object data
func (builder aesCfbCipherBuilder) createCipherData() (CipherData, error) {
	var cd CipherData
	var err error
	err = cd.RandomKeyIv(builder.aesKeySize, builder.aesKeySize)
	if err != nil {
		return cd, err
	}

	cd.WrapAlgorithm = builder.MasterCipher.GetWrapAlgorithm()
	cd.CEKAlgorithm = "AES/CFB/NoPadding"
	cd.MatDesc = builder.MasterCipher.GetMatDesc()

	// EncryptedKey
	cd.EncryptedKey, err = builder.MasterCipher.Encrypt(cd.Key)
	if err != nil {
		return cd, err
	}

	// EncryptedIV
	cd.EncryptedIV, err = builder.MasterCipher.Encrypt(cd.IV)
	if err != nil {
		return cd, err
	}

	return cd, nil
}

// contentCipherCD is used to create ContentCipher with CipherData
func (builder aesCfbCipherBuilder) contentCipherCD(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCfb(cd)
	if err != nil {
		return nil, err
	}
	return &aesCfbCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}

// ContentCipher is used to create ContentCipher interface
func (builder aesCfbCipherBuilder) ContentCipher() (ContentCipher, error) {
	cd, err := builder.createCipherData()
	if err != nil {
		return nil, err
	}
	return builder.contentCipherCD(cd)
}

// ContentCipherEnv is used to create a decrption ContentCipher from Envelope
func (builder aesCfbCipherBuilder) ContentCipherEnv(envelope Envelope) (ContentCipher, error) {
	var cd CipherData
	cd.EncryptedKey = make([]byte, len(envelope.CipherKey))
	copy(cd.EncryptedKey, []byte(envelope.CipherKey))

	plainKey, err := builder.MasterCipher.Decrypt([]byte(envelope.CipherKey))
	if err != nil {
		return nil, err
	}

	cd.Key = make([]byte, len(plainKey))
	copy(cd.Key, plainKey)

	cd.EncryptedIV = make([]byte, len(envelope.IV))
	copy(cd.EncryptedIV, []byte(envelope.IV))

	plainIV, err := builder.MasterCipher.Decrypt([]byte(envelope.IV))
	if err != nil {
		return nil, err
	}

	cd.IV = make([]byte, len(plainIV))
	copy(cd.IV, plainIV)

	cd.MatDesc = envelope.MatDesc
	cd.WrapAlgorithm = envelope.WrapAlg
	cd.CEKAlgorithm = envelope.CEKAlg

	if len(envelope.WrapAlg) <= 0 {
		cd.WrapAlgorithm = builder.MasterCipher.GetWrapAlgorithm()
	}
	if len(envelope.CEKAlg) <= 0 {
		cd.CEKAlgorithm = "AES/CFB/NoPadding"
	}
	return builder.contentCipherCD(cd)
}

// GetMatDesc is used to get MasterCipher's MatDesc
func (builder aesCfbCipherBuilder) GetMatDesc() string {
	return builder.MasterCipher.GetMatDesc()
}

// EncryptContents will generate a random key and iv and encrypt the data using ctr
func (cc *aesCfbCipher) EncryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Encrypt(src)
	return &CryptoEncrypter{Body: src, Encrypter: reader}, nil
}

// DecryptContent is used to decrypt object using ctr
func (cc *aesCfbCipher) DecryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Decrypt(src)
	return &CryptoDecrypter{Body: src, Decrypter: reader}, nil
}

// GetCipherData is used to get cipher data information
func (cc *aesCfbCipher) GetCipherData() *CipherData {
	return &(cc.CipherData)
}

// GetCipherData returns cipher data
func (cc *aesCfbCipher) GetEncryptedLen(plainTextLen int64) int64 {
	// AES CTR encryption mode does not change content length
	return plainTextLen
}

// GetAlignLen is used to get align length
func (cc *aesCfbCipher) GetAlignLen() int {
	return len(cc.CipherData.IV)
}

// Clone is used to create a new aesCtrCipher from itself
func (cc *aesCfbCipher) Clone(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCfb(cd)
	if err != nil {
		return nil, err
	}

	return &aesCfbCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}
