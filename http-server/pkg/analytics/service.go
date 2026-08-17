package analytics

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	nameNoiseSuffixes = []string{
		" subscription", " subscriptions", " premium", " pro", " plus", " family", " plan", " basic", " standard",
		" individual", " duo", " max", " lite", " mini", " student", " business",
		" подписка", " подписки", " премиум", " про", " плюс", " семейная", " тариф", " базовый",
	}
	nameNoiseReplacer = strings.NewReplacer(
		"ё", "е",
		"+", " plus ",
		"&", " and ",
		"™", " ",
		"®", " ",
		"©", " ",
		".", " ",
		"_", " ",
		"-", " ",
		"/", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"{", " ",
		"}", " ",
		",", " ",
		";", " ",
		":", " ",
		"!", " ",
		"?", " ",
	)
	nameBoundaryPattern = regexp.MustCompile(`[^a-z0-9а-я]+`)
)

type serviceIndex struct {
	bySlug map[string]Service
	exact  map[string]string
	keys   []string
}

func newServiceIndex(services []Service) *serviceIndex {
	index := &serviceIndex{
		bySlug: make(map[string]Service, len(services)),
		exact:  make(map[string]string, len(services)*4),
	}

	seen := make(map[string]struct{})

	for _, service := range services {
		index.bySlug[service.Slug] = service
		for _, raw := range append([]string{service.Name, service.Slug}, service.Aliases...) {
			key := normalizeServiceName(raw)
			if key == "" || utf8.RuneCountInString(key) < 2 {
				continue
			}

			if _, exists := index.exact[key]; !exists {
				index.exact[key] = service.Slug
			}

			if _, exists := seen[key]; exists {
				continue
			}

			seen[key] = struct{}{}
			index.keys = append(index.keys, key)
		}
	}

	for i := 0; i < len(index.keys); i++ {
		for j := i + 1; j < len(index.keys); j++ {
			if len(index.keys[j]) > len(index.keys[i]) {
				index.keys[i], index.keys[j] = index.keys[j], index.keys[i]
			}
		}
	}

	return index
}

func MatchService(services []Service, name string) (Service, bool) {
	return newServiceIndex(services).match(name)
}

func (index *serviceIndex) match(name string) (Service, bool) {
	normalized := normalizeServiceName(name)
	if normalized == "" {
		return Service{}, false
	}

	if slug, ok := index.exact[normalized]; ok {
		return index.bySlug[slug], true
	}

	padded := " " + normalized + " "
	for _, key := range index.keys {
		if strings.Contains(padded, " "+key+" ") {
			if slug, ok := index.exact[key]; ok {
				return index.bySlug[slug], true
			}
		}
	}

	return Service{}, false
}

func normalizeServiceName(input string) string {
	normalized := strings.TrimSpace(strings.ToLower(nameNoiseReplacer.Replace(input)))
	normalized = nameBoundaryPattern.ReplaceAllString(normalized, " ")
	normalized = strings.Join(strings.Fields(normalized), " ")

	for {
		next := normalized
		for _, suffix := range nameNoiseSuffixes {
			next = strings.TrimSuffix(next, suffix)
		}

		next = strings.TrimSpace(next)
		if next == normalized {
			break
		}

		normalized = next
	}

	return normalized
}

func displayName(sub Subscription, matched Service) string {
	if matched.Name != "" {
		return matched.Name
	}

	return sub.Name
}
