package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	if resp, err := http.Get("http://localhost:8000/login"); err != nil {
		t.Error(err)
	} else {
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n")
		loginCookies := resp.Cookies()
		resp.Body.Close()
		if req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/home", nil); err != nil {
			t.Error(err)
		} else {
			for _, cookie := range loginCookies {
				fmt.Printf("receive cookie %s = %s \n", cookie.Name, cookie.Value)
				// cookie.Value += "1"	// simulate a request with a wrong cookies
				req.AddCookie(cookie) // simulate a request with a cookie
			}

			client := &http.Client{Timeout: 3 * time.Second}
			if resp, err := client.Do(req); err != nil {
				t.Error(err)
			} else {
				defer resp.Body.Close()
				fmt.Println("response body")
				io.Copy(os.Stdout, resp.Body)
				os.Stdout.WriteString("\n")
			}
		}
	}
}
