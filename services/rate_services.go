package services

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/MJarmolkiewicz/kryptonim/models"
	"github.com/avast/retry-go"
)

type openExchangeFetcher struct{}

func NewRateService() *openExchangeFetcher {
	return &openExchangeFetcher{}
}

// func (e openExchangeFetcher) FetchRatesUSD() (map[string]*big.Float, error) {
// 	appID := os.Getenv("OXR_API_KEY")
// 	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("non-200 response from OpenExchangeRates: %d", resp.StatusCode)
// 	}

// 	var data models.OpenRatesResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %w", err)
// 	}

// 	// Convert to big.Float for precision
// 	rates := make(map[string]*big.Float)
// 	for currency, value := range data.Rates {
// 		rates[currency] = new(big.Float).SetPrec(128).SetFloat64(value)
// 	}

// 	return rates, nil
// }

func (e openExchangeFetcher) FetchRatesUSD() (map[string]*big.Float, error) {
	appID := os.Getenv("OXR_API_KEY")
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)
	var data models.OpenRatesResponse

	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("non-200 response: %d", resp.StatusCode)
			}

			return json.NewDecoder(resp.Body).Decode(&data)
		},
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch or decode: %w", err)
	}

	rates := make(map[string]*big.Float)
	for currency, value := range data.Rates {
		rates[currency] = new(big.Float).SetPrec(128).SetFloat64(value)
	}

	return rates, nil
}
