package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
		DualStack: true,
	}).DialContext,
}

var httpClient = &http.Client{
	Transport: transport,
}

func main() {
	var cookies string
	var threads int
	urls := make(chan string)
	var wg sync.WaitGroup
	flag.StringVar(&cookies, "c", "", "Specify Cookie Header")
	flag.IntVar(&threads, "t", 20, "Number of threads to run")
	flag.Parse()
	//Read stdin
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for input.Scan() {
			urls <- input.Text()
		}
		close(urls)
	}()
	for i := 0; i < threads; i++ {
		wg.Add(1)
		worker(urls, &wg, cookies)
	}
	wg.Done()
}

func checkclickjack(url string, cookie string) {
	headers := map[string]string{
		"Cache-Control": "no-cache",
		"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Cookie":        cookie,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for a, b := range headers {
		request.Header.Add(a, b)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if len(resp.Header.Get("X-Frame-Options")) == 0 {
		fmt.Println(url)
	}
}

func worker(cha chan string, wg *sync.WaitGroup, cookie string) {
	for i := range cha {
		checkclickjack(i, cookie)
	}
	wg.Done()
}
