package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type JwtPayloadType string

var (
	JwtPayloadTypeAccessToken  JwtPayloadType = "access-token"
	JwtPayloadTypeRefreshToken JwtPayloadType = "refresh-token"
)

type JwtPayload struct {
	Sub string `json:"sub"`
}

func DecryptString(key, ciphertext string) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.Wrap(err, "failed decrypt string")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "failed decrypt string")
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "failed decrypt string")
	}

	nonceSize := aesGcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext is too short")
	}

	nonce, encryptedData := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	plaintext, err := aesGcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed decrypt string")
	}

	return string(plaintext), nil
}

func EncryptString(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

func HashPassword(password string, cost int) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes)
}

func GenerateToken(secret string, expirationInHour int, jwtPayload JwtPayload) (string, map[string]any, error) {
	claims := jwt.MapClaims{
		"iss": "marketplace",
		"sub": jwtPayload.Sub,
		"exp": time.Now().Add(time.Hour * time.Duration(expirationInHour)).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, claims, err
}
