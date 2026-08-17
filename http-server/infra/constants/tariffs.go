package constants

import "strings"

const (
	TariffNone       = "none"
	TariffBasic      = "basic"
	TariffStandard   = "standard"
	TariffPlus       = "plus"
	TariffPro        = "pro"
	TariffPremium    = "premium"
	TariffMax        = "max"
	TariffLite       = "lite"
	TariffMini       = "mini"
	TariffStudent    = "student"
	TariffDuo        = "duo"
	TariffFamily     = "family"
	TariffIndividual = "individual"
	TariffBusiness   = "business"
)

func NormalizeTariff(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return TariffNone
	}

	return normalized
}
