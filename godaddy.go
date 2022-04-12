package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type record struct {
	Data string `form:"data" json:"data" binding:"required"`
}

func httpGet(url, key, secret string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 120}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("new request err: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client do err:: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call server failed: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func httpPost(url, body, key, secret string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 120}

	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		log.Println("new request err: ", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client do err:: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("call godaddy server failed: %d", resp.StatusCode)
	}

	log.Println("call godaddy server success. record has changed.")

	response, err := ioutil.ReadAll(resp.Body)
	log.Println(string(response))
	return err
}

func updateGodaddy(address string, key, secret string) error {
	url := "https://api.godaddy.com/v1/domains/ieevee.com/records/A/e"

	records := make([]record, 0, 1)
	records = append(records, record{Data: address})
	body, err := json.Marshal(records)
	if err != nil {
		return err
	}
	log.Println(string(body))

	return httpPost(url, string(body), key, secret)
}

func getGodaddy(key, secret string) (string, error) {
	url := "https://api.godaddy.com/v1/domains/ieevee.com/records/A/e"

	reponse, err := httpGet(url, key, secret)
	if err != nil {
		return "", err
	}

	records := make([]record, 1)

	err = json.Unmarshal(reponse, &records)
	if err != nil {
		return "", err
	}

	log.Printf("get records: %v", records[0].Data)
	return records[0].Data, nil
}

func getExternalIP() (addr string, err error) {
	resp, err := http.Get("http://ip4only.me/api/")
	if err != nil {
		log.Printf("call externalip failed, %v", err)
		return "", err
	}

	response, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("response", string(response))
	s := strings.Split(string(response), ",")
	if len(s) <= 2 {
		return "", fmt.Errorf("response error")
	}
	addr = s[1]
	return addr, nil
}

func main() {
	ticker := time.NewTicker(time.Hour)

	if len(os.Args) != 3 {
		log.Println("needs key and secret")
		return
	}

	key := os.Args[1]
	secret := os.Args[2]

	stopCh := make(chan interface{})
	lastAddr, err := getGodaddy(key, secret)
	if err != nil {
		panic(err)
	}

	go func(ch chan interface{}, key, secret string) {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			addr, err := getExternalIP()
			if err != nil {
				log.Printf("call externalip failed, %v", err)
				continue
			}
			if addr != lastAddr {
				log.Printf("external ip address changed from %s to %s", lastAddr, addr)
				lastAddr = addr
				err := updateGodaddy(addr, key, secret)
				if err != nil {
					log.Printf("call godaddy failed: %v", err)
				}
			} else {
				log.Println("external ip not change, continue")
			}
		}
		ch <- 0
	}(stopCh, key, secret)

	<-stopCh
}
