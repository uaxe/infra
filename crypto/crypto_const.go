package crypto

const (
	EncryptionKey                      string = "encryption-key"
	EncryptionIv                              = "encryption-iv"
	EncryptionCekAlg                          = "encryption-cek-alg"
	EncryptionWrapAlg                         = "encryption-wrap-alg"
	EncryptionVersion                         = "encryption-version"
	EncryptionMatDesc                         = "encryption-matdesc"
	EncryptionUnencryptedContentLength        = "encryption-unencrypted-content-length"
	EncryptionUnencryptedContentMD5           = "encryption-unencrypted-content-md5"
	EncryptionDataSize                        = "encryption-data-size"
	EncryptionPartSize                        = "encryption-part-size"
)

// encryption Algorithm
const (
	RsaCryptoWrap   string = "RSA/NONE/PKCS1Padding"
	AesCtrAlgorithm string = "AES/CTR/NoPadding"
	AesCfbAlgorithm string = "AES/CFB/NoPadding"
)
