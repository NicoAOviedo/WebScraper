package main

import (
	"net/http"
	"fmt"
)

func main() {
	statusCode,err := getStatusCode("https://x.com/home")
	if err != nil {
		fmt.Println("Error: ",err)
		return
	}
	fmt.Println("Status Code: ", statusCode)

}

func getStatusCode(url string) (int, error) {

	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode,nil
}

