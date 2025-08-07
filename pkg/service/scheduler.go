package service

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "crypto-watcher/pkg/db"
)

func StartScheduler() {
    ticker := time.NewTicker(10 * time.Second)

    for {
        <-ticker.C
        coins := db.GetTrackedCoins()
        for _, coin := range coins {
            go fetchAndStore(coin)
        }
    }
}

func fetchAndStore(coin string) {
    url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin)

    resp, err := http.Get(url)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var data map[string]map[string]float64
    if err := json.Unmarshal(body, &data); err != nil {
        return
    }

    price := data[coin]["usd"]
    ts := time.Now().Unix()
    db.SavePrice(coin, price, ts)
}