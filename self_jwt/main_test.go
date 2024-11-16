package main

import (
	"fmt"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	header := DefaultHeader
	payload := JwtPayload{
		ID:         "asdfgfhmj",
		Issue:      "yaya.com",
		Audience:   "myself",
		Subject:    "購買點數",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(2 * time.Hour).Unix(),
		UserDefined: map[string]any{
			"user":   "小陳",
			"item":   "icash",
			"amount": 1000,
		},
	}

	token, err := GenJwt(header, payload)
	if err != nil {
		t.Fatal("failed to verify Jwt toekn:", err)
	}
	fmt.Println("token:", token)

	_, p, err := VerifyJwt(token)
	if err != nil {
		t.Fatal("failed to verify Jwt toekn:", err)
	}

	fmt.Printf("Jwt verify success: %+v\n", p)
}
