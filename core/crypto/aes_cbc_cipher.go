package crypto

import (
	"fmt"
	"io"
)

// aesCbcCipherBuilder for building ContentCipher
type aesCbcCipherBuilder struct {
	MasterCipher MasterCipher
	aesKeySize   int
}

// aesCbcCipher will use aes ctr algorithm
type aesCbcCipher struct {
	CipherData CipherData
	Cipher     Cipher
}

// CreateAesCbcCipher creates ContentCipherBuilder
func CreateAesCbcCipher(cipher MasterCipher, aesType AesType) ContentCipherBuilder {
	switch aesType {
	default:
		panic(fmt.Sprintf("aes type unknown %d", aesType))
	case AES_TYPE_128, AES_TYPE_192, AES_TYPE_256:
		break
	}
	return aesCbcCipherBuilder{MasterCipher: cipher, aesKeySize: int(aesType)}
}

// createCipherData create CipherData for encrypt object data
func (builder aesCbcCipherBuilder) createCipherData() (CipherData, error) {
	var cd CipherData
	var err error
	err = cd.RandomKeyIv(builder.aesKeySize, builder.aesKeySize)
	if err != nil {
		return cd, err
	}

	cd.WrapAlgorithm = builder.MasterCipher.GetWrapAlgorithm()
	cd.CEKAlgorithm = "AES/CBC/NoPadding"
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
func (builder aesCbcCipherBuilder) contentCipherCD(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCtr(cd)
	if err != nil {
		return nil, err
	}

	return &aesCbcCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}

// ContentCipher is used to create ContentCipher interface
func (builder aesCbcCipherBuilder) ContentCipher() (ContentCipher, error) {
	cd, err := builder.createCipherData()
	if err != nil {
		return nil, err
	}
	return builder.contentCipherCD(cd)
}

// ContentCipherEnv is used to create a decrption ContentCipher from Envelope
func (builder aesCbcCipherBuilder) ContentCipherEnv(envelope Envelope) (ContentCipher, error) {
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

	return builder.contentCipherCD(cd)
}

// GetMatDesc is used to get MasterCipher's MatDesc
func (builder aesCbcCipherBuilder) GetMatDesc() string {
	return builder.MasterCipher.GetMatDesc()
}

// EncryptContents will generate a random key and iv and encrypt the data using ctr
func (cc *aesCbcCipher) EncryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Encrypt(src)
	return &CryptoEncrypter{Body: src, Encrypter: reader}, nil
}

// DecryptContent is used to decrypt object using ctr
func (cc *aesCbcCipher) DecryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Decrypt(src)
	return &CryptoDecrypter{Body: src, Decrypter: reader}, nil
}

// GetCipherData is used to get cipher data information
func (cc *aesCbcCipher) GetCipherData() *CipherData {
	return &(cc.CipherData)
}

// GetCipherData returns cipher data
func (cc *aesCbcCipher) GetEncryptedLen(plainTextLen int64) int64 {
	// AES CTR encryption mode does not change content length
	return plainTextLen
}

// GetAlignLen is used to get align length
func (cc *aesCbcCipher) GetAlignLen() int {
	return len(cc.CipherData.IV)
}

// Clone is used to create a new aesCbcCipher from itself
func (cc *aesCbcCipher) Clone(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCbc(cd)
	if err != nil {
		return nil, err
	}

	return &aesCtrCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}
