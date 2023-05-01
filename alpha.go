package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type question struct{
	text string
}

type answer struct{
	text string 
}

func main() {

	fmt.Printf("Enter your question : ")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')

	input = strings.ReplaceAll(input, " ", "+")
	input = strings.ReplaceAll(input, "\n", "")

	resp, getErr := http.Get("http://api.wolframalpha.com/v1/result?appid=H4RUYV-2WP4YG72XJ&i=" + input)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println(string(body))

}
