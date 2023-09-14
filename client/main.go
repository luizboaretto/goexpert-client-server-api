package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	api_url := "http://localhost:8080/cotacao"

	req, err := http.NewRequestWithContext(ctx, "GET", api_url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		var exchangerate ExchangeRate
		err = json.Unmarshal(body, &exchangerate)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create("cotacao.txt")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		file_size, err := fmt.Fprintln(f, "Dólar: ", exchangerate.Bid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Exchange Rate saved. Size: %d bytes\n", file_size)

	}
}
