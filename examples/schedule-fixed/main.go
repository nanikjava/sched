package main

import (
	"encoding/json"
	"fmt"
	"github.com/sherifabdlnaby/sched"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fixedTimer30second, err := sched.NewFixed(3 * time.Second)
	if err != nil {
		panic(fmt.Sprintf("invalid interval: %s", err.Error()))
	}

	job := func() {
		log.Println("Doing some work...")
		//time.Sleep(1 * time.Second)
		getBitcoinData()
		log.Println("Finished Work.")
	}

	// Create Schedule
	schedule := sched.NewSchedule("every30s", fixedTimer30second, job, sched.WithLogger(sched.DefaultLogger()))

	// Start Schedule
	schedule.Start()

	// Listen to CTRL + C And indefintly wait shutdown.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	_ = <-signalChan

	// Stop before shutting down.
	schedule.Stop()

	return
}

type CoinDeskCurrentPriceInAud struct {
	Bpi struct {
		Aud struct {
			Code        string  `json:"code"`
			Description string  `json:"description"`
			Rate        string  `json:"rate"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"AUD"`
		Usd struct {
			Code        string  `json:"code"`
			Description string  `json:"description"`
			Rate        string  `json:"rate"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
	} `json:"bpi"`
	Disclaimer string `json:"disclaimer"`
	Time       struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
}

func getBitcoinData() {
	// https://api.coindesk.com/v1/bpi/currentprice/aud.json
	/**
	{
	    "bpi": {
	        "AUD": {
	            "code": "AUD",
	            "description": "Australian Dollar",
	            "rate": "47,630.0410",
	            "rate_float": 47630.041
	        },
	        "USD": {
	            "code": "USD",
	            "description": "United States Dollar",
	            "rate": "35,617.7433",
	            "rate_float": 35617.7433
	        }
	    },
	    "disclaimer": "This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org",
	    "time": {
	        "updated": "Jun 20, 2021 03:32:00 UTC",
	        "updatedISO": "2021-06-20T03:32:00+00:00",
	        "updateduk": "Jun 20, 2021 at 04:32 BST"
	    }
	}
	*/
	response, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice/aud.json")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	var currentPrice CoinDeskCurrentPriceInAud
	json.Unmarshal([]byte(responseData), &currentPrice)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentPrice.Bpi.Aud.RateFloat)
}
