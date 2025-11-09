package price

import (
	"time"

	"github.com/angad363/stocky-assignment/pkg/logger"
	"github.com/jmoiron/sqlx"
)

func StartPriceUpdater(service *PriceService, conn *sqlx.DB) {
	logger.Log.Info("💹 Price updater started")
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			updateAllHoldings(service, conn)
		}
	}()
}

func updateAllHoldings(service *PriceService, conn *sqlx.DB) {
	logger.Log.Info("🔄 Starting hourly stock price update...")

	rows, err := conn.Query("SELECT DISTINCT stock_symbol FROM ledger_entries")
	if err != nil {
		logger.Log.Errorf("Error fetching symbols: %v", err)
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
			logger.Log.WithField("symbol", symbol).Errorf("Failed to update price: %v", err)
			continue
		}
		logger.Log.WithFields(map[string]interface{}{
			"symbol": symbol,
			"price":  priceResp.Price,
		}).Info("Updated stock price")
	}

	logger.Log.Info("✅ Hourly stock price update completed")
}
