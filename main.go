package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var api_key string = os.Getenv("API_KEY")
var api_passphrase string = os.Getenv("API_PASSPHRASE")
var api_secret_key string = os.Getenv("API_SECRET_KEY")

func main() {
	endpoint := "/api/v1/accounts"
	now := time.Now().Local().UnixNano() / int64(time.Millisecond)
	signature := create_header("kucoin", "GET", endpoint, now)
	passphrase := sign_string(api_passphrase, api_secret_key)
	headers := http.Header{
		"KC-API-SIGN":        []string{signature},
		"KC-API-TIMESTAMP":   []string{strconv.FormatInt(now, 10)},
		"KC-API-KEY":         []string{api_key},
		"KC-API-PASSPHRASE":  []string{passphrase},
		"KC-API-KEY-VERSION": []string{strconv.Itoa(2)},
	}
	get_result(headers, endpoint, now)
}

func create_header(client_type, method, endpoint string, now int64) string {
	to_sign := strconv.FormatInt(now, 10) + method + endpoint
	signature := sign_string(to_sign, api_secret_key)
	return signature
}

func sign_string(sign string, secret_key string) string {
	hmac := hmac.New(sha256.New, []byte(secret_key))
	hmac.Write([]byte(sign))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func get_result(headers http.Header, endpoint string, now int64) {
	client := &http.Client{}
	base_url := "https://api.kucoin.com"
	url := base_url + endpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}
