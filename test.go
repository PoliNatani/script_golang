package main

import (

	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
	"bufio"
	"os"
	"sync"
)

func Read_url(ch_url chan<- string){
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ch_url <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
        log.Println(err)
    }
	close(ch_url)
}
	

func Get_requests(ch_url <-chan string, res chan int){
	wg := sync.WaitGroup{}
	for url := range ch_url{
		wg.Add(1)
		go func(url string, res chan int){
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			pattern := `Go[\W\s]+`
			re := regexp.MustCompile(pattern)
			matches := re.FindAllString(string(body), -1)
			url += ":"
			fmt.Println("Count for", url, len(matches))
			res <- len(matches)
		}(url, res)
	}
	wg.Wait()
	close(res)
}

func Total(res <-chan int) (total int){
	total = 0
	for count := range res{
		total += count
	}
	return
	
}

func main() { 
	chan_url := make(chan string, 5)
	res := make(chan int)
	go Read_url(chan_url)
	go Get_requests(chan_url, res)
	fmt.Println("Total:", Total(res))	
}