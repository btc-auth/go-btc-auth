# Bitcoin Private Key Authentication Protocol

## Table of Contents

1. Introduction
2. Protocol Overview
3. Technical Architecture
   3.1. Server-Side Components
   3.2. Client-Side Requirements
4. Authentication Flow
   4.1. Token Generation
   4.2. QR Code Creation
   4.3. Client Signature
   4.4. Server Verification
5. Security Considerations
6. Implementation Details
7. Future Enhancements
8. Conclusion

## 1. Introduction

This white paper presents a novel authentication protocol that leverages Bitcoin private keys for secure user authentication. By utilizing the cryptographic principles underlying Bitcoin, this protocol offers a robust, decentralized authentication mechanism suitable for various applications.

## 2. Protocol Overview

The Bitcoin Private Key Authentication Protocol enables users to authenticate themselves to a service using their Bitcoin private keys, without revealing the keys themselves. This is achieved through a challenge-response mechanism involving digital signatures.

Key features:
- Utilizes existing Bitcoin infrastructure
- No password storage on servers
- Resistant to phishing and man-in-the-middle attacks
- Provides a seamless user experience through QR code scanning

## 3. Technical Architecture

### 3.1. Server-Side Components

The server architecture consists of the following key components:

1. **Token Manager**: Generates and manages authentication tokens.
2. **QR Code Generator**: Creates QR codes containing authentication challenges.
3. **Signature Verifier**: Validates client-signed messages using public keys.
4. **Session Manager**: Maintains authenticated user sessions.

### 3.2. Client-Side Requirements

Clients need a Bitcoin wallet capable of:
- Scanning QR codes
- Signing arbitrary messages with their private key
- Exporting the corresponding public key

## 4. Authentication Flow

### 4.1. Token Generation

1. The server generates a unique token for each authentication attempt using a cryptographically secure random number generator.
2. Tokens are associated with a creation timestamp and have configurable expiration times, typically set to a short duration (e.g., 60 seconds) to minimize the risk of replay attacks.

### 4.2. QR Code Creation

1. The server constructs an authentication URL with the following format:
   ```
   http://<server_host>:<port>/auth?token=<generated_token>&sign={sign}&public_key_hex={public_key_hex}
   ```
   Where:
   - `<server_host>` is the domain or IP address of the authentication server
   - `<port>` is the port number the server is listening on
   - `<generated_token>` is the unique token generated in step 4.1
   - `{sign}` is a placeholder for the client's signature
   - `{public_key_hex}` is a placeholder for the client's public key in hexadecimal format

2. This URL is encoded into a QR code using a library such as go-qrcode.

### 4.3. Client Signature

1. The client scans the QR code using their Bitcoin wallet application.
2. The wallet extracts the authentication URL from the QR code.
3. The wallet uses the user's private key to create an ECDSA (Elliptic Curve Digital Signature Algorithm) signature of the entire authentication URL. The process is as follows:
   a. The wallet hashes the URL using SHA-256.
   b. The resulting hash is signed using the user's private key and the secp256k1 elliptic curve (the same curve used in Bitcoin).
   c. The signature is typically produced in DER (Distinguished Encoding Rules) format.
4. The wallet converts the DER signature to a hexadecimal string.
5. The wallet also exports the user's public key and converts it to a hexadecimal string.
6. The wallet replaces the `{sign}` and `{public_key_hex}` placeholders in the URL with the actual signature and public key hexadecimal strings.
7. The completed URL is sent back to the server as an HTTP GET request.

### 4.4. Server Verification

1. The server receives the GET request containing the completed authentication URL.
2. It extracts the token, signature, and public key from the URL parameters.
3. The server verifies the token hasn't expired by checking its creation timestamp against the current time.
4. If the token is valid, the server proceeds to verify the signature:
   a. It reconstructs the original authentication URL using the received token and the `{sign}` and `{public_key_hex}` placeholders.
   b. The server hashes this reconstructed URL using SHA-256.
   c. It then uses the provided public key to verify that the signature matches the hash of the reconstructed URL.
5. The signature verification process uses the following steps:
   a. The public key is parsed from its hexadecimal representation.
   b. The signature is parsed from its hexadecimal representation and converted from DER format.
   c. The ECDSA verification algorithm is used to check if the signature is valid for the hash of the reconstructed URL using the provided public key.
6. If the signature is valid, the server associates the public key with the user's session, effectively authenticating the user.

This process ensures that only a user in possession of the private key corresponding to the provided public key could have produced a valid signature for the given authentication URL. The use of a unique token for each authentication attempt prevents replay attacks, as a previously used URL will contain an expired token and thus be rejected by the server.

## 5. Security Considerations

- **Token Expiration**: Tokens have short lifespans to minimize the risk of replay attacks.
- **Signature Verification**: Robust signature verification prevents forgery attempts.
- **No Private Key Transmission**: Private keys never leave the client device.
- **Unique Challenges**: Each authentication attempt uses a unique token, preventing replay attacks.

## 6. Implementation Details

The protocol is implemented using the following technologies:

- **Server**: Go (Golang) with the Gin web framework
- **Cryptography**: btcd/btcec library for Bitcoin-compatible elliptic curve operations
- **QR Code Generation**: go-qrcode library
- **Frontend**: HTML templates with embedded QR codes

Key code snippets:

```go
func VerifySignature(message, signatureHex, publicKeyHex string) (bool, error) {
    // Implementation details...
}

func GenerateSignature(message string, privateKey *btcec.PrivateKey) (string, error) {
    // Implementation details...
}
```

## 7. Future Enhancements

- Multi-factor authentication options
- Integration with hardware wallets
- Support for other cryptocurrency key pairs
- Distributed server architecture for improved scalability

## 8. Conclusion

The Bitcoin Private Key Authentication Protocol offers a secure, user-friendly authentication mechanism that leverages the robust cryptographic foundations of Bitcoin. By eliminating the need for password storage and utilizing existing wallet infrastructure, this protocol provides a promising alternative to traditional authentication methods.