package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type CurrencyItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var extraCurrencies = map[string]string{
	"RUB": "Russian Ruble",
}

var (
	currenciesCache    []CurrencyItem
	currenciesCacheExp time.Time
	currenciesCacheMu  sync.RWMutex
	currenciesCacheTTL = 24 * time.Hour
)

func GetCurrencies(ctx context.Context) ([]CurrencyItem, error) {
	currenciesCacheMu.RLock()
	if len(currenciesCache) > 0 && time.Now().Before(currenciesCacheExp) {
		items := make([]CurrencyItem, len(currenciesCache))
		copy(items, currenciesCache)
		currenciesCacheMu.RUnlock()
		return items, nil
	}
	currenciesCacheMu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultAPIBase+"/currencies", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currency list request failed with status %d", resp.StatusCode)
	}

	var payload map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	byCode := make(map[string]string, len(payload)+len(extraCurrencies))
	for code, name := range payload {
		byCode[strings.ToUpper(code)] = name
	}

	for code, name := range extraCurrencies {
		byCode[strings.ToUpper(code)] = name
	}

	items := make([]CurrencyItem, 0, len(byCode))
	for code, name := range byCode {
		items = append(items, CurrencyItem{Code: code, Name: name})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})

	currenciesCacheMu.Lock()
	currenciesCache = items
	currenciesCacheExp = time.Now().Add(currenciesCacheTTL)
	currenciesCacheMu.Unlock()

	result := make([]CurrencyItem, len(items))
	copy(result, items)

	return result, nil
}
