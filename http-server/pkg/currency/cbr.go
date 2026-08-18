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

const cbrDailyURL = "https://www.cbr-xml-daily.ru/daily_json.js"

type cbrValute struct {
	Nominal float64 `json:"Nominal"`
	Value   float64 `json:"Value"`
}

type cbrResponse struct {
	Valute map[string]cbrValute `json:"Valute"`
}

var (
	cbrCache    map[string]float64
	cbrCacheExp time.Time
	cbrCacheMu  sync.RWMutex
	cbrCacheTTL = 6 * time.Hour
)

func InvalidateCBRCache() {
	cbrCacheMu.Lock()
	cbrCache = nil
	cbrCacheExp = time.Time{}
	cbrCacheMu.Unlock()
}

func rubPerUnit(code string) (float64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "RUB" {
		return 1, nil
	}

	cbrCacheMu.RLock()
	if len(cbrCache) > 0 && time.Now().Before(cbrCacheExp) {
		rate, ok := cbrCache[code]
		cbrCacheMu.RUnlock()
		if ok && rate > 0 {
			return rate, nil
		}

		return 0, fmt.Errorf("CBR rate not found for %s", code)
	}
	cbrCacheMu.RUnlock()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, cbrDailyURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("CBR request failed with status %d", resp.StatusCode)
	}

	var payload cbrResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}

	rates := make(map[string]float64, len(payload.Valute)+1)
	rates["RUB"] = 1

	for code, valute := range payload.Valute {
		if valute.Nominal <= 0 || valute.Value <= 0 {
			continue
		}

		rates[strings.ToUpper(code)] = valute.Value / valute.Nominal
	}

	cbrCacheMu.Lock()
	cbrCache = rates
	cbrCacheExp = time.Now().Add(cbrCacheTTL)
	cbrCacheMu.Unlock()

	rate, ok := rates[code]
	if !ok || rate <= 0 {
		return 0, fmt.Errorf("CBR rate not found for %s", code)
	}

	return rate, nil
}

func getCBRRate(_ context.Context, from, to string) (float64, error) {
	from = strings.ToUpper(strings.TrimSpace(from))
	to = strings.ToUpper(strings.TrimSpace(to))

	fromRub, err := rubPerUnit(from)
	if err != nil {
		return 0, err
	}

	toRub, err := rubPerUnit(to)
	if err != nil {
		return 0, err
	}

	return fromRub / toRub, nil
}

func needsCBR(from, to string) bool {
	return from == "RUB" || to == "RUB"
}
