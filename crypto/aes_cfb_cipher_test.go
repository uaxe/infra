package crypto

import (
	. "gopkg.in/check.v1"
)

func (s *CryptoSuite) TestCfbContentEncryptCipherError(c *C) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, rsaPublicKey, rsaPrivateKey)
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)
	cc, err := contentProvider.ContentCipher()
	c.Assert(err, IsNil)

	var cipherData CipherData
	cipherData.RandomKeyIv(31, 15)

	_, err = cc.Clone(cipherData)
	c.Assert(err, NotNil)
}

func (s *CryptoSuite) TestCfbCreateCipherDataError(c *C) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.createCipherData()
	c.Assert(err, NotNil)
}

func (s *CryptoSuite) TestCfbContentCipherCDError(c *C) {
	var cd CipherData

	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.contentCipherCD(cd)
	c.Assert(err, NotNil)
}
