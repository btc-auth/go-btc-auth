package main

// 从命令行获取url 对 url 进行签名
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"auth_client/libs"

	"github.com/btcsuite/btcd/btcec/v2"
)

var private_key string = "b56eba6102cd43a6daaf91832edee78a10df84faeb646fad587f2d361ec564e6"

func fillUrlPlaceholders(url string, placeholders map[string]string) string {
	for key, value := range placeholders {
		url = strings.Replace(url, fmt.Sprintf(`{%s}`, key), value, -1)
	}
	return url
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL as a command-line argument")
		return
	}

	// privateKey, err := btcec.NewPrivateKey()
	// 从 private_key 生成私钥
	privateKeyBytes, err := hex.DecodeString(private_key)
	if err != nil {
		fmt.Printf("Error decoding private key: %v\n", err)
		return
	}

	// Create private key from bytes
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)

	// privateKey, err := btcec.PrivKeyFromBytes(btcec.S256(), []byte(private_key))

	// 打印私钥的 hex 格式
	fmt.Printf("Private Key (hex): %x\n", privateKey.Serialize())
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())

	url := os.Args[1]
	sign, err := libs.GenerateSignature(url, privateKey)
	if err != nil {
		fmt.Printf("Error generating signature: %v", err)
		return
	}
	fmt.Println("Signature (hex):", sign)

	new_url := fillUrlPlaceholders(url, map[string]string{
		"sign":           sign,
		"public_key_hex": publicKeyHex,
	})
	fmt.Println("New URL:", new_url)

	// 访问 new_url,把结果 json 格式化输出
	response, err := http.Get(new_url)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Error decoding JSON response: %v", err)
		return
	}

	formattedJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v", err)
		return
	}

	fmt.Println("Formatted JSON response:")
	fmt.Println(string(formattedJSON))

}
