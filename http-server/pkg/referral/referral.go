package referral

import (
	"crypto/rand"
	"strings"
	"unicode"
)

const (
	StartPrefix    = "ref_"
	codeLength     = 8
	codeAlphabet   = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	PremiumPlan    = "Premium"
	RankNovice     = 1
	RankPartner    = 2
	RankAmbassador = 3
	RankLeader     = 4
)

type Rank struct {
	Level      int
	Key        string
	MinCount   int
	MaxCount   *int
	RewardDays int
}

func intPtr(value int) *int {
	return &value
}

var Ranks = []Rank{
	{Level: RankNovice, Key: "novice", MinCount: 0, MaxCount: intPtr(1), RewardDays: 0},
	{Level: RankPartner, Key: "partner", MinCount: 2, MaxCount: intPtr(25), RewardDays: 30},
	{Level: RankAmbassador, Key: "ambassador", MinCount: 26, MaxCount: intPtr(100), RewardDays: 90},
	{Level: RankLeader, Key: "leader", MinCount: 101, MaxCount: nil, RewardDays: 365},
}

func RankForCount(count int) Rank {
	if count < 0 {
		count = 0
	}

	current := Ranks[0]
	for _, rank := range Ranks {
		if count >= rank.MinCount {
			current = rank
		}
	}

	return current
}

func SanitizeCode(code string) string {
	var builder strings.Builder

	for _, char := range strings.TrimSpace(code) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func GenerateCode() (string, error) {
	buf := make([]byte, codeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.Grow(codeLength)

	for _, b := range buf {
		builder.WriteByte(codeAlphabet[int(b)%len(codeAlphabet)])
	}

	return builder.String(), nil
}

func StartToken(code string) string {
	return StartPrefix + code
}

func CodeFromStartToken(token string) string {
	token = strings.TrimSpace(token)
	if !strings.HasPrefix(strings.ToLower(token), StartPrefix) {
		return ""
	}

	return SanitizeCode(token[len(StartPrefix):])
}

func SiteURL(appURL, code string) string {
	appURL = strings.TrimRight(strings.TrimSpace(appURL), "/")
	if appURL == "" {
		appURL = "https://paylist.site"
	}

	return appURL + "/sign-up?ref=" + code
}

func BotURL(botName, code string) string {
	botName = strings.TrimPrefix(strings.TrimSpace(botName), "@")
	if botName == "" || code == "" {
		return ""
	}

	return "https://t.me/" + botName + "?start=" + StartToken(code)
}
