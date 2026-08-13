package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultAPIBase = "https://api.frankfurter.dev/v1"

type frankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

type rateCacheEntry struct {
	rate      float64
	expiresAt time.Time
}

var (
	rateCache    = map[string]rateCacheEntry{}
	rateCacheMu  sync.RWMutex
	rateCacheTTL = 30 * time.Minute
)

func cacheKey(from, to string) string {
	return strings.ToUpper(from) + "->" + strings.ToUpper(to)
}

func getCachedRate(from, to string) (float64, bool) {
	rateCacheMu.RLock()
	defer rateCacheMu.RUnlock()

	entry, ok := rateCache[cacheKey(from, to)]
	if !ok || time.Now().After(entry.expiresAt) {
		return 0, false
	}

	return entry.rate, true
}

func setCachedRate(from, to string, rate float64) {
	rateCacheMu.Lock()
	defer rateCacheMu.Unlock()

	rateCache[cacheKey(from, to)] = rateCacheEntry{
		rate:      rate,
		expiresAt: time.Now().Add(rateCacheTTL),
	}
}

func GetRate(ctx context.Context, from, to string) (float64, error) {
	from = strings.ToUpper(strings.TrimSpace(from))
	to = strings.ToUpper(strings.TrimSpace(to))

	if from == "" || to == "" {
		return 0, fmt.Errorf("currency codes are required")
	}

	if from == to {
		return 1, nil
	}

	if needsCBR(from, to) {
		rate, err := getCBRRate(ctx, from, to)
		if err != nil {
			return 0, err
		}

		setCachedRate(from, to, rate)
		return rate, nil
	}

	if rate, ok := getCachedRate(from, to); ok {
		return rate, nil
	}

	url := fmt.Sprintf("%s/latest?from=%s&to=%s", defaultAPIBase, from, to)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("currency rate request failed with status %d", resp.StatusCode)
	}

	var payload frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}

	rate, ok := payload.Rates[to]
	if !ok || rate <= 0 {
		return 0, fmt.Errorf("rate not found for %s -> %s", from, to)
	}

	setCachedRate(from, to, rate)

	return rate, nil
}

func Convert(ctx context.Context, from, to string, amount float64) (float64, error) {
	rate, err := GetRate(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}

func GetMonthlyAmount(price float64, period string) float64 {
	if period == "yearly" {
		return price / 12
	}

	return price
}

func GetRates(ctx context.Context, base string, symbols []string) (map[string]float64, error) {
	base = strings.ToUpper(strings.TrimSpace(base))
	if base == "" {
		return nil, fmt.Errorf("base currency is required")
	}

	filtered := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.ToUpper(strings.TrimSpace(symbol))
		if symbol == "" || symbol == base {
			continue
		}

		filtered = append(filtered, symbol)
	}

	if len(filtered) == 0 {
		return map[string]float64{}, nil
	}

	if needsCBR(base, filtered[0]) || containsSymbol(filtered, "RUB") {
		rates := make(map[string]float64, len(filtered))
		for _, symbol := range filtered {
			rate, err := GetRate(ctx, base, symbol)
			if err != nil {
				return nil, err
			}

			rates[symbol] = rate
		}

		return rates, nil
	}

	url := fmt.Sprintf("%s/latest?from=%s&to=%s", defaultAPIBase, base, strings.Join(filtered, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currency rates request failed with status %d", resp.StatusCode)
	}

	var payload frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	for _, symbol := range filtered {
		if rate, ok := payload.Rates[symbol]; ok {
			setCachedRate(base, symbol, rate)
		}
	}

	return payload.Rates, nil
}

func containsSymbol(symbols []string, target string) bool {
	for _, symbol := range symbols {
		if symbol == target {
			return true
		}
	}

	return false
}
