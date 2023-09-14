package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type ExchangeRateResponse struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite3", "./exchange_rate.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	checkTable(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getExchangeRate(db))
	http.ListenAndServe(":8080", mux)

}

func checkTable(db *sql.DB) {
	_, err := db.Query("select * from exchange_rate")
	if err != nil {
		db.Exec("create table exchange_rate(id integer primary key autoincrement, code varchar(10), codein varchar(10), name varchar(50), high decimal(10,5), low decimal(10,5), varBid decimal(10,5), pctChange decimal(10,5), bid decimal(10,5), ask decimal(10,5), timestamp integer, create_date datetime);")
	}
}

func getExchangeRate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchange_rate, err := getExchangeRateFromApi()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = insertExchangeRate(db, exchange_rate)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		bid := ExchangeRateResponse{Bid: exchange_rate.USDBRL.Bid}
		json.NewEncoder(w).Encode(bid)
	}
}

func getExchangeRateFromApi() (ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return ExchangeRate{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ExchangeRate{}, err
	}
	defer resp.Body.Close()
	var er ExchangeRate
	err = json.NewDecoder(resp.Body).Decode(&er)
	if err != nil {
		return ExchangeRate{}, err
	}
	return er, nil
}

func insertExchangeRate(db *sql.DB, er ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("insert into exchange_rate (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, er.USDBRL.Code, er.USDBRL.Codein, er.USDBRL.Name, er.USDBRL.High, er.USDBRL.Low, er.USDBRL.VarBid, er.USDBRL.PctChange, er.USDBRL.Bid, er.USDBRL.Ask, er.USDBRL.Timestamp, er.USDBRL.CreateDate)
	return err
}
