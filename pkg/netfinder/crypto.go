package netfinder

import (
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"io"
	"net"

	"golang.org/x/crypto/chacha20poly1305"
)

// ECDH密钥对管理器
type keyManager struct {
	privateKey *ecdh.PrivateKey
	publicKey  string // Base64编码的公钥
}

// 生成新的ECDH密钥对
func generateKeyPair() (*keyManager, error) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成ECDH密钥对失败: %w", err)
	}

	return &keyManager{
		privateKey: privateKey,
		publicKey:  base64.StdEncoding.EncodeToString(privateKey.PublicKey().Bytes()),
	}, nil
}

// 获取公钥（Base64编码）
func (km *keyManager) publicKeyBase64() string {
	return km.publicKey
}

// 计算共享密钥
func (km *keyManager) computeSharedSecret(peerPublicKeyBase64 string) ([]byte, error) {
	peerPubKeyBytes, err := base64.StdEncoding.DecodeString(peerPublicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("解码对方公钥失败: %w", err)
	}

	peerPublicKey, err := ecdh.X25519().NewPublicKey(peerPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("解析对方公钥失败: %w", err)
	}

	sharedSecret, err := km.privateKey.ECDH(peerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("ECDH密钥交换失败: %w", err)
	}

	return sharedSecret, nil
}

// 派生加密密钥
func deriveEncryptionKey(sharedSecret []byte) []byte {
	h := fnv.New64a()
	h.Write(sharedSecret)
	hash := h.Sum(nil)

	// 使用哈希值作为ChaCha20密钥
	key := make([]byte, chacha20poly1305.KeySize)
	copy(key, hash)
	return key
}

// 流式加密器接口
type StreamCipher interface {
	// Read 从加密连接读取并解密数据
	Read(p []byte) (n int, err error)
	// Write 加密并写入数据
	Write(p []byte) (n int, err error)
	// Close 关闭连接
	Close() error
}

// 加密连接
type encryptedConn struct {
	conn   net.Conn
	cipher cipher.AEAD
	nonce  []byte
}

// 创建加密连接（发送方使用）
func newEncryptedConn(conn net.Conn, sharedSecret []byte) (*encryptedConn, error) {
	// 生成随机nonce
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("生成nonce失败: %w", err)
	}

	// 派生加密密钥
	key := deriveEncryptionKey(sharedSecret)

	// 创建加密器
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("创建加密器失败: %w", err)
	}

	return &encryptedConn{
		conn:   conn,
		cipher: cipher,
		nonce:  nonce,
	}, nil
}

// 递增nonce
func (ec *encryptedConn) incrementNonce() {
	for i := len(ec.nonce) - 1; i >= 0; i-- {
		ec.nonce[i]++
		if ec.nonce[i] != 0 {
			break
		}
	}
}

// Write 加密并写入数据
func (ec *encryptedConn) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// 加密数据
	ciphertext := ec.cipher.Seal(nil, ec.nonce, p, nil)

	// 写入加密数据长度（4字节大端序）
	length := uint32(len(ciphertext))
	lengthBytes := make([]byte, 4)
	lengthBytes[0] = byte(length >> 24)
	lengthBytes[1] = byte(length >> 16)
	lengthBytes[2] = byte(length >> 8)
	lengthBytes[3] = byte(length)

	if _, err := ec.conn.Write(lengthBytes); err != nil {
		return 0, fmt.Errorf("写入数据长度失败: %w", err)
	}

	// 写入加密数据
	if _, err := ec.conn.Write(ciphertext); err != nil {
		return 0, fmt.Errorf("写入加密数据失败: %w", err)
	}

	// 递增nonce
	ec.incrementNonce()

	return len(p), nil
}

// Read 读取并解密数据
func (ec *encryptedConn) Read(p []byte) (n int, err error) {
	// 读取加密数据长度
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(ec.conn, lengthBytes); err != nil {
		return 0, fmt.Errorf("读取数据长度失败: %w", err)
	}

	length := (uint32(lengthBytes[0]) << 24) |
		(uint32(lengthBytes[1]) << 16) |
		(uint32(lengthBytes[2]) << 8) |
		uint32(lengthBytes[3])

	// 读取加密数据
	ciphertext := make([]byte, length)
	if _, err := io.ReadFull(ec.conn, ciphertext); err != nil {
		return 0, fmt.Errorf("读取加密数据失败: %w", err)
	}

	// 解密数据
	plaintext, err := ec.cipher.Open(nil, ec.nonce, ciphertext, nil)
	if err != nil {
		return 0, fmt.Errorf("解密数据失败: %w", err)
	}

	// 递增nonce
	ec.incrementNonce()

	// 复制到目标缓冲区
	copy(p, plaintext)
	return len(plaintext), nil
}

// Close 关闭连接
func (ec *encryptedConn) Close() error {
	return ec.conn.Close()
}
