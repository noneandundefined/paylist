package analytics

import (
	"math"
	"strings"
	"time"
)

const (
	crowdMinSample     = 5
	crowdOverpayRatio  = 1.20
	expensiveShareMin  = 0.25
	concentrationShare = 0.40
	smallShareMax      = 0.05
	maxRecommendations = 6
)

type Service struct {
	Slug     string
	Name     string
	Category string
	Aliases  []string
}

type Subscription struct {
	ID                 uint64
	Name               string
	Tariff             string
	Price              float64
	SharePrice         float64
	SharePercent       float64
	Currency           string
	Period             string
	DatePay            time.Time
	IncludeInAnalytics bool
	IsOwner            bool
}

type CrowdPrice struct {
	Name     string
	Tariff   string
	Price    float64
	Currency string
	Period   string
	Country  string
}

type Recommendation struct {
	ID             string         `json:"id"`
	Type           string         `json:"type"`
	TitleKey       string         `json:"title_key"`
	DescKey        string         `json:"desc_key"`
	DescValues     map[string]any `json:"desc_values,omitempty"`
	SubscriptionID *uint64        `json:"subscription_id,omitempty"`
	Priority       float64        `json:"-"`
}

type Input struct {
	Subscriptions   []Subscription
	Categories      map[uint64][]string
	Services        []Service
	Crowd           []CrowdPrice
	DisplayCurrency string
	Country         string
	Now             time.Time
	Convert         func(amount float64, from string) float64
}

func ptrID(id uint64) *uint64 {
	value := id
	return &value
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func monthlyAmount(price float64, period string) float64 {
	if strings.EqualFold(period, "yearly") {
		return price / 12
	}

	return price
}

func tariffOrNone(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "none"
	}

	return normalized
}

func countryCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}
