package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"auth_server/libs"

	"github.com/btcsuite/btcd/btcec/v2"
)

func TestGenerateSignature(t *testing.T) {
	// 设置随机数种子
	// rand.Seed(uint64(time.Now().UnixNano()))

	// 生成一个新的私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Error generating private key: %v", err)
	}

	// 获取对应的公钥
	publicKey := privateKey.PubKey()

	// 要签名的消息
	message := "Hello, World!"

	// 生成签名
	signatureHex, err := libs.GenerateSignature(message, privateKey)
	if err != nil {
		t.Fatalf("Error generating signature: %v", err)
	}

	// 打印结果
	fmt.Println("Message:", message)
	fmt.Println("Signature (hex):", signatureHex)
	fmt.Println("Public Key (hex):", hex.EncodeToString(publicKey.SerializeCompressed()))

	// 验证签名
	valid, err := libs.VerifySignature(message, signatureHex, hex.EncodeToString(publicKey.SerializeCompressed()))
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}

	if !valid {
		t.Fatalf("Signature verification failed")
	}
}
