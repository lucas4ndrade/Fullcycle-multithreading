package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type apiResponse struct {
	apiName      string
	responseJSON string
}

func main() {
	fmt.Println("Enter cep:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputCEP := scanner.Text()

	responseCh := make(chan apiResponse)

	go getAddressFromBrasilAPI(inputCEP, responseCh)
	go getAddressFromViaCEP(inputCEP, responseCh)

	select {
	case r := <-responseCh:
		fmt.Println("")
		fmt.Printf("Received a response from %s\n", r.apiName)
		fmt.Println("===================================")
		fmt.Println(r.responseJSON)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out :(")
	}
}

func getAddressFromViaCEP(cep string, ch chan<- apiResponse) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get address from viacep API!! %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request to viacep API failed with status code %d!!", res.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to decode viacep API response!! %v", err)
		return
	}

	ch <- apiResponse{
		apiName:      "Viacep API",
		responseJSON: string(bodyBytes),
	}
}

func getAddressFromBrasilAPI(cep string, ch chan<- apiResponse) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get address from brasil API!! %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request to brasil API failed with status code %d!!", res.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to decode brasil API response!! %v", err)
		return
	}

	ch <- apiResponse{
		apiName:      "Brasil API",
		responseJSON: string(bodyBytes),
	}
}
