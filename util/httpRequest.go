package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// RequestData http获取数据
func RequestData(url string) ([]byte, error) {
	c := http.Client{
		Timeout: 180 * time.Second,
	}
	fmt.Printf("get URL :%+v\n", url)
	resp, err := c.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("request res:success\n")
	return body, err
}
