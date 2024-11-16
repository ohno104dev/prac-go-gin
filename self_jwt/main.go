package selfjwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// use std library to implement jwt encoding/decoding feature

const (
	JWT_SECRET = "abc1234567890" // just for test. use env or configuration file to keep in server side
)

var (
	DefaultHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"`
	Type string `json:"typ"`
}

type JwtPayload struct {
	ID          string         `json:"jti"` // JWT ID標示該JWT
	Issue       string         `json:"iss"` // 發行者
	Audience    string         `json:"aud"` // 受眾
	Subject     string         `json:"sub"` // 主題
	IssueAt     int64          `json:"iat"` // 發布時間 sec
	NotBefore   int64          `json:"nbf"` // 在此之前之前不可用 sec
	Expiration  int64          `json:"exp"` // 到期時間 sec
	UserDefined map[string]any `json:"ud"`  // 用戶自訂
}

func GenJwt(header JwtHeader, payload JwtPayload) (string, error) {
	var part1, part2, signature string

	// header to JSON and base64 encoder
	if bs1, err := json.Marshal(header); err != nil {
		return "", err
	} else {
		part1 = base64.RawURLEncoding.EncodeToString(bs1) // does not include special characters in the URL
	}

	// payload to JSON and base64 encoder
	if bs2, err := json.Marshal(payload); err != nil {
		return "", err
	} else {
		part2 = base64.RawURLEncoding.EncodeToString(bs2)
	}

	// authentication using sha256
	h := hmac.New(sha256.New, []byte(JWT_SECRET))
	h.Write([]byte(part1 + "." + part2))
	// signature
	signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return part1 + "." + part2 + "." + signature, nil
}

func VerifyJwt(token string) (*JwtHeader, *JwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("token verify failed, format error")
	}

	h := hmac.New(sha256.New, []byte(JWT_SECRET))
	h.Write([]byte(parts[0] + "." + parts[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != parts[2] {
		return nil, nil, fmt.Errorf("token verify failed, signature error")
	}

	var part1, part2 []byte
	var err error
	if part1, err = base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return nil, nil, fmt.Errorf("jwt header decode failed: %v", err)
	}

	if part2, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return nil, nil, fmt.Errorf("jwt payload decode failed: %v", err)
	}

	var header JwtHeader
	var payload JwtPayload
	if err = json.Unmarshal(part1, &header); err != nil {
		return nil, nil, fmt.Errorf("header to JSON failed")
	}
	if err = json.Unmarshal(part2, &payload); err != nil {
		return nil, nil, fmt.Errorf("payload to JSON failed")
	}

	return &header, &payload, nil
}
