package libs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func VerifySignature(message, signatureHex, publicKeyHex string) (bool, error) {
	// 将公钥从十六进制字符串转换为字节
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %v", err)
	}

	// 解析公钥
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %v", err)
	}

	// 将签名从十六进制字符串转换为字节
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %v", err)
	}

	// 解析签名
	signature, err := ecdsa.ParseDERSignature(signatureBytes)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %v", err)
	}

	// 计算消息的哈希
	messageHash := sha256.Sum256([]byte(message))

	// 验证签名
	return signature.Verify(messageHash[:], publicKey), nil
}

func GenerateSignature(message string, privateKey *btcec.PrivateKey) (string, error) {
	// 计算消息的哈希
	messageHash := sha256.Sum256([]byte(message))

	// 使用私钥签名消息哈希
	signature := ecdsa.Sign(privateKey, messageHash[:])

	// 将签名序列化为 DER 格式
	derSignature := signature.Serialize()

	// 将 DER 格式的签名转换为十六进制字符串
	return hex.EncodeToString(derSignature), nil
}
