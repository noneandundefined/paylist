package country

import (
	"sort"
	"strings"
)

type CountryItem struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	InflationRate float64 `json:"inflation_rate"`
}

// InflationRate is an estimated annual subscription price growth:
// historical CPI + typical digital-services premium (~2.5%).
var countries = []CountryItem{
	{Code: "US", Name: "United States", InflationRate: 5.5},
	{Code: "GB", Name: "United Kingdom", InflationRate: 5.2},
	{Code: "DE", Name: "Germany", InflationRate: 4.8},
	{Code: "FR", Name: "France", InflationRate: 4.6},
	{Code: "IT", Name: "Italy", InflationRate: 4.5},
	{Code: "ES", Name: "Spain", InflationRate: 4.4},
	{Code: "NL", Name: "Netherlands", InflationRate: 4.7},
	{Code: "BE", Name: "Belgium", InflationRate: 4.6},
	{Code: "AT", Name: "Austria", InflationRate: 4.5},
	{Code: "CH", Name: "Switzerland", InflationRate: 3.2},
	{Code: "SE", Name: "Sweden", InflationRate: 4.3},
	{Code: "NO", Name: "Norway", InflationRate: 4.4},
	{Code: "DK", Name: "Denmark", InflationRate: 4.2},
	{Code: "FI", Name: "Finland", InflationRate: 4.1},
	{Code: "PL", Name: "Poland", InflationRate: 6.8},
	{Code: "CZ", Name: "Czech Republic", InflationRate: 5.9},
	{Code: "HU", Name: "Hungary", InflationRate: 7.5},
	{Code: "RO", Name: "Romania", InflationRate: 7.2},
	{Code: "BG", Name: "Bulgaria", InflationRate: 6.5},
	{Code: "GR", Name: "Greece", InflationRate: 5.1},
	{Code: "PT", Name: "Portugal", InflationRate: 4.8},
	{Code: "IE", Name: "Ireland", InflationRate: 5.0},
	{Code: "UA", Name: "Ukraine", InflationRate: 9.8},
	{Code: "RU", Name: "Russia", InflationRate: 9.5},
	{Code: "BY", Name: "Belarus", InflationRate: 8.8},
	{Code: "KZ", Name: "Kazakhstan", InflationRate: 8.2},
	{Code: "TR", Name: "Turkey", InflationRate: 12.5},
	{Code: "IL", Name: "Israel", InflationRate: 5.8},
	{Code: "AE", Name: "United Arab Emirates", InflationRate: 4.0},
	{Code: "SA", Name: "Saudi Arabia", InflationRate: 4.2},
	{Code: "IN", Name: "India", InflationRate: 6.2},
	{Code: "CN", Name: "China", InflationRate: 3.8},
	{Code: "JP", Name: "Japan", InflationRate: 3.5},
	{Code: "KR", Name: "South Korea", InflationRate: 4.1},
	{Code: "SG", Name: "Singapore", InflationRate: 4.0},
	{Code: "HK", Name: "Hong Kong", InflationRate: 3.6},
	{Code: "TW", Name: "Taiwan", InflationRate: 3.7},
	{Code: "TH", Name: "Thailand", InflationRate: 4.5},
	{Code: "VN", Name: "Vietnam", InflationRate: 5.4},
	{Code: "ID", Name: "Indonesia", InflationRate: 5.6},
	{Code: "MY", Name: "Malaysia", InflationRate: 4.9},
	{Code: "PH", Name: "Philippines", InflationRate: 6.0},
	{Code: "AU", Name: "Australia", InflationRate: 5.3},
	{Code: "NZ", Name: "New Zealand", InflationRate: 5.1},
	{Code: "CA", Name: "Canada", InflationRate: 5.0},
	{Code: "MX", Name: "Mexico", InflationRate: 6.5},
	{Code: "BR", Name: "Brazil", InflationRate: 7.8},
	{Code: "AR", Name: "Argentina", InflationRate: 15.0},
	{Code: "CL", Name: "Chile", InflationRate: 6.2},
	{Code: "CO", Name: "Colombia", InflationRate: 7.0},
	{Code: "ZA", Name: "South Africa", InflationRate: 7.4},
	{Code: "EG", Name: "Egypt", InflationRate: 8.5},
	{Code: "NG", Name: "Nigeria", InflationRate: 9.2},
}

const defaultInflationRate = 5.0

func GetCountries() []CountryItem {
	items := make([]CountryItem, len(countries))
	copy(items, countries)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items
}

func GetInflationRate(code string) float64 {
	normalized := strings.ToUpper(strings.TrimSpace(code))

	for _, item := range countries {
		if item.Code == normalized {
			return item.InflationRate
		}
	}

	return defaultInflationRate
}

func FindCountry(code string) (*CountryItem, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(code))

	for _, item := range countries {
		if item.Code == normalized {
			copyItem := item
			return &copyItem, true
		}
	}

	return nil, false
}
