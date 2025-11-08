package price

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func StartPriceUpdater(service *PriceService, conn *sqlx.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	// For testing, you can use: time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				updateAllHoldings(service, conn)
			}
		}
	}()
}

func updateAllHoldings(service *PriceService, conn *sqlx.DB) {
	log.Println("🔄 Starting hourly stock price update...")

	rows, err := conn.Query("SELECT DISTINCT stock_symbol FROM ledger_entries")
	if err != nil {
		log.Printf("Error fetching symbols: %v", err)
		return
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err == nil {
			symbols = append(symbols, symbol)
		}
	}

	for _, symbol := range symbols {
		priceResp, err := service.GetStockPrice(symbol)
		if err != nil {
			log.Printf("Error updating price for %s: %v", symbol, err)
			continue
		}
		log.Printf("✅ Updated %s -> %.2f INR", symbol, priceResp.Price)
	}

	log.Println("✅ Hourly stock price update completed")
}
