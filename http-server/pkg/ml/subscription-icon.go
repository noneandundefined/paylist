package ml

import (
	"regexp"
	"strings"
)

type subscriptionBrand struct {
	Image    string
	Keywords []string
}

var subscriptionBrands = []subscriptionBrand{
	{Image: "yplus.png", Keywords: []string{
		"yandex plus", "yandex premium", "yandex music", "yandex pl", "yandexplus", "яндекс плюс", "яндекс+", "яндексплюс",
		"яндекс премиум", "яндекс музыка", "яндекс подписка", "yplus",
	}},
	{Image: "sber.jpg", Keywords: []string{
		"sber prime", "sberprime", "сбер прайм", "сберпрайм", "сбербанк прайм", "sberbank prime",
		"sber", "сбер", "сбербанк",
	}},
	{Image: "tbank.jpg", Keywords: []string{
		"tinkoff pro", "tinkoff black", "tinkoff premium", "t-bank pro", "t bank pro", "tbank pro",
		"tbank", "tinkoff", "t-bank", "t bank", "тинькофф", "тинькоф",
	}},
	{Image: "prime.png", Keywords: []string{
		"amazon prime", "prime video", "amazon video", "амазон прайм", "prime amazon",
	}},
	{Image: "netflix.jpg", Keywords: []string{
		"netflix", "нетфликс",
	}},
	{Image: "spotify.jpg", Keywords: []string{
		"spotify", "спотифай", "spotify premium", "spotify family",
	}},
	{Image: "youtube.png", Keywords: []string{
		"youtube premium", "youtube music", "youtube family", "yt premium", "ютуб премиум", "ютуб музыка",
		"youtube",
	}},
	{Image: "apple.png", Keywords: []string{
		"apple one", "apple music", "apple tv+", "apple tv plus", "apple arcade", "apple fitness",
		"icloud+", "icloud plus", "эпл one", "эпл музыка", "apple",
	}},
	{Image: "googleone.jpg", Keywords: []string{
		"google one", "google storage", "google drive", "google ai pro", "google workspace",
		"гугл one", "гугл диск", "google play pass",
	}},
	{Image: "gemini.png", Keywords: []string{
		"google gemini", "gemini advanced", "gemini ai", "google ai ultra",
	}},
	{Image: "office365.jpg", Keywords: []string{
		"microsoft 365", "office 365", "microsoft office", "onedrive", "офис 365",
	}},
	{Image: "adobe.jpg", Keywords: []string{
		"adobe creative cloud", "adobe cc", "creative cloud", "photoshop", "lightroom", "adobe",
	}},
	{Image: "chatgpt.png", Keywords: []string{
		"chatgpt plus", "chatgpt pro", "chat gpt plus", "openai plus", "openai pro", "chatgpt", "openai", "чатgpt", "чат gpt",
	}},
	{Image: "claude.png", Keywords: []string{
		"claude pro", "claude ai", "anthropic claude", "claude", "anthropic",
	}},
	{Image: "midjourney.png", Keywords: []string{
		"midjourney", "mid journey",
	}},
	{Image: "copilot.png", Keywords: []string{
		"github copilot", "copilot pro", "microsoft copilot", "github pro",
	}},
	{Image: "cursor.png", Keywords: []string{
		"cursor pro", "cursor ai", "cursor",
	}},
	{Image: "perplexity.png", Keywords: []string{
		"perplexity pro", "perplexity ai", "perplexity",
	}},
	{Image: "discord.png", Keywords: []string{
		"discord nitro", "discord", "nitro", "дискорд",
	}},
	{Image: "tg.png", Keywords: []string{
		"telegram premium", "telegram", "телеграм премиум", "телеграм",
	}},
	{Image: "steam.png", Keywords: []string{
		"steam", "стим",
	}},
	{Image: "xbox.png", Keywords: []string{
		"xbox game pass", "xbox live", "game pass ultimate", "xbox", "геймпас",
	}},
	{Image: "psplus.jpg", Keywords: []string{
		"playstation plus", "ps plus", "psn plus", "ps plus extra", "пс плюс", "playstation",
	}},
	{Image: "nintendo.png", Keywords: []string{
		"nintendo switch online", "switch online", "nintendo online", "nintendo",
	}},
	{Image: "twitch.png", Keywords: []string{
		"twitch turbo", "twitch", "твич",
	}},
	{Image: "kinopoisk.jpg", Keywords: []string{
		"kinopoisk", "кинопоиск", "kino poisk",
	}},
	{Image: "ivi.jpg", Keywords: []string{
		"ivi", "иви",
	}},
	{Image: "okko.png", Keywords: []string{
		"okko", "окко",
	}},
	{Image: "wink.png", Keywords: []string{
		"wink", "wink music", "винк",
	}},
	{Image: "start.png", Keywords: []string{
		"start ru", "start.ru", "start", "старту",
	}},
	{Image: "premier.png", Keywords: []string{
		"premier", "premier one", "премьер", "руспремьер",
	}},
	{Image: "rutube.png", Keywords: []string{
		"rutube premium", "rutube", "рутуб",
	}},
	{Image: "dzen.png", Keywords: []string{
		"dzen premium", "yandex dzen", "яндекс дзен", "dzen", "дзен",
	}},
	{Image: "disney.png", Keywords: []string{
		"disney+", "disney plus", "disney", "дисней",
	}},
	{Image: "max.png", Keywords: []string{
		"hbo max", "max streaming", "hbo", "hbomax",
	}},
	{Image: "crunchyroll.png", Keywords: []string{
		"crunchyroll", "crunchy roll", "кранчиролл",
	}},
	{Image: "deezer.png", Keywords: []string{
		"deezer", "deezer premium", "дизер",
	}},
	{Image: "tidal.png", Keywords: []string{
		"tidal", "tidal hifi",
	}},
	{Image: "notion.png", Keywords: []string{
		"notion plus", "notion ai", "notion",
	}},
	{Image: "figma.png", Keywords: []string{
		"figma professional", "figma organization", "figma",
	}},
	{Image: "canva.png", Keywords: []string{
		"canva pro", "canva teams", "canva",
	}},
	{Image: "dropbox.png", Keywords: []string{
		"dropbox plus", "dropbox professional", "dropbox",
	}},
	{Image: "icloud.png", Keywords: []string{
		"icloud", "icloud storage", "айклауд",
	}},
	{Image: "github.png", Keywords: []string{
		"github", "git hub",
	}},
	{Image: "slack.png", Keywords: []string{
		"slack pro", "slack business", "slack",
	}},
	{Image: "zoom.png", Keywords: []string{
		"zoom pro", "zoom workplace", "zoom one", "zoom",
	}},
	{Image: "linkedin.png", Keywords: []string{
		"linkedin premium", "linkedin sales", "linkedin",
	}},
	{Image: "coursera.png", Keywords: []string{
		"coursera plus", "coursera",
	}},
	{Image: "udemy.png", Keywords: []string{
		"udemy pro", "udemy personal", "udemy",
	}},
	{Image: "duolingo.png", Keywords: []string{
		"duolingo max", "duolingo plus", "duolingo", "дуолинго",
	}},
	{Image: "stepik.jpg", Keywords: []string{
		"stepik", "степик",
	}},
	{Image: "geekbrains.png", Keywords: []string{
		"geekbrains", "гикбрейнс", "skillbox", "скиллбокс", "netology", "нетология",
	}},
	{Image: "patreon.png", Keywords: []string{
		"patreon", "пatreon",
	}},
	{Image: "vpn.png", Keywords: []string{
		"nordvpn", "surfshark", "expressvpn", "proton vpn", "windscribe", "mullvad", "vpn",
		"впн",
	}},
	{Image: "cloudflare.png", Keywords: []string{
		"cloudflare", "cloudflare zero trust",
	}},
	{Image: "hetzner.png", Keywords: []string{
		"hetzner", "hetzner cloud",
	}},
	{Image: "digitalocean.png", Keywords: []string{
		"digitalocean", "digital ocean", "do droplet",
	}},
	{Image: "vds.jpg", Keywords: []string{
		"timeweb cloud", "timeweb", "selectel", "vdsina", "vds", "dedicated server", "дедик",
	}},
	{Image: "vps.jpg", Keywords: []string{
		"vps", "virtual server", "reg.ru", "regru",
	}},
	{Image: "mts.jpg", Keywords: []string{
		"mts premium", "mts music", "mts", "мтс",
	}},
	{Image: "beeline.png", Keywords: []string{
		"beeline", "билайн",
	}},
	{Image: "megafon.png", Keywords: []string{
		"megafon", "мегафон",
	}},
	{Image: "t2.png", Keywords: []string{
		"tele2", "t2 russia", "t2",
	}},
	{Image: "domru.png", Keywords: []string{
		"dom.ru", "dom ru", "дом ру", "дом.ру",
	}},
	{Image: "ozon.png", Keywords: []string{
		"ozon premium", "ozon", "озон",
	}},
	{Image: "wb.png", Keywords: []string{
		"wildberries plus", "wildberries", "wb plus", "вайлдберриз", "wildberries",
	}},
	{Image: "samokat.png", Keywords: []string{
		"samokat", "самокат",
	}},
	{Image: "vk.png", Keywords: []string{
		"vk music", "vk video", "vk combo", "vkontakte", "vk музыка", "vk видео",
		"вконтакте", "вк музыка", "вк видео", "вк комбо",
	}},
	{Image: "subscriptionplus.png", Keywords: []string{
		"subscription plus", "subscriptionplus", "sub plus",
	}},
	{Image: "buscard.png", Keywords: []string{
		"transport card", "транспортная карта", "troika", "тройка",
	}},
	{Image: "money.png", Keywords: []string{
		"loan payment", "credit payment", "перевод", "кредит", "займ",
	}},

	{Image: "paramount.png", Keywords: []string{
		"paramount+", "paramount plus", "paramount", "парамаунт",
	}},
	{Image: "peacock.png", Keywords: []string{
		"peacock", "peacock premium", "peacock tv",
	}},
	{Image: "audible.png", Keywords: []string{
		"audible", "audible plus", "аудible",
	}},
	{Image: "soundcloud.png", Keywords: []string{
		"soundcloud", "soundcloud go", "soundcloud go+",
	}},
	{Image: "pandora.png", Keywords: []string{
		"pandora", "pandora plus", "pandora premium",
	}},
	{Image: "ea.png", Keywords: []string{
		"ea play", "origin access", "electronic arts",
	}},
	{Image: "ubisoft.png", Keywords: []string{
		"ubisoft+", "ubisoft plus", "ubisoft connect",
	}},
	{Image: "nvidia.png", Keywords: []string{
		"geforce now", "nvidia geforce now",
	}},
	{Image: "roblox.png", Keywords: []string{
		"roblox premium", "roblox",
	}},
	{Image: "evernote.png", Keywords: []string{
		"evernote", "evernote personal",
	}},
	{Image: "todoist.png", Keywords: []string{
		"todoist", "todoist pro",
	}},
	{Image: "bitwarden.png", Keywords: []string{
		"bitwarden", "bitwarden premium", "bitwarden family",
	}},
	{Image: "onepassword.png", Keywords: []string{
		"1password", "one password",
	}},
	{Image: "grammarly.png", Keywords: []string{
		"grammarly", "grammarly premium",
	}},
	{Image: "jetbrains.png", Keywords: []string{
		"jetbrains", "intellij idea", "pycharm", "webstorm", "datagrip",
	}},
	{Image: "linear.png", Keywords: []string{
		"linear app", "linear team",
	}},
	{Image: "trello.png", Keywords: []string{
		"trello", "trello premium", "atlassian",
	}},
	{Image: "miro.png", Keywords: []string{
		"miro", "miro starter", "miro business",
	}},
	{Image: "airtable.png", Keywords: []string{
		"airtable", "airtable plus",
	}},
	{Image: "asana.png", Keywords: []string{
		"asana", "asana premium",
	}},
	{Image: "clickup.png", Keywords: []string{
		"clickup", "click up",
	}},
	{Image: "monday.png", Keywords: []string{
		"monday.com", "monday com", "monday work",
	}},
	{Image: "aws.png", Keywords: []string{
		"amazon web services", "aws", "amazon aws",
	}},
	{Image: "azure.png", Keywords: []string{
		"microsoft azure", "azure",
	}},
	{Image: "proton.png", Keywords: []string{
		"proton mail", "proton drive", "proton unlimited", "proton pass", "proton",
	}},
	{Image: "headspace.png", Keywords: []string{
		"headspace", "headspace plus",
	}},
	{Image: "calm.png", Keywords: []string{
		"calm", "calm premium",
	}},
	{Image: "strava.png", Keywords: []string{
		"strava", "strava summit",
	}},
	{Image: "peloton.png", Keywords: []string{
		"peloton", "peloton app",
	}},
	{Image: "medium.png", Keywords: []string{
		"medium membership", "medium",
	}},
	{Image: "substack.png", Keywords: []string{
		"substack", "sub stack",
	}},
	{Image: "nytimes.png", Keywords: []string{
		"new york times", "nytimes", "ny times",
	}},
	{Image: "revolut.png", Keywords: []string{
		"revolut premium", "revolut metal", "revolut",
	}},
	{Image: "wise.png", Keywords: []string{
		"wise premium", "wise", "transferwise",
	}},
	{Image: "runway.png", Keywords: []string{
		"runway ml", "runway gen", "runway",
	}},
	{Image: "elevenlabs.png", Keywords: []string{
		"elevenlabs", "eleven labs",
	}},
	{Image: "uber.png", Keywords: []string{
		"uber one", "uber pass", "uber eats pass",
	}},
	{Image: "plex.png", Keywords: []string{
		"plex pass", "plex",
	}},
	{Image: "boosty.png", Keywords: []string{
		"boosty", "boosty.to", "бусти",
	}},
	{Image: "yandexeda.png", Keywords: []string{
		"yandex eda", "yandex food", "яндекс еда", "eda yandex",
	}},
	{Image: "deliveryclub.png", Keywords: []string{
		"delivery club", "deliveryclub", "деливери",
	}},
	{Image: "kuper.png", Keywords: []string{
		"kuper", "sbermarket", "сбермаркет", "купер",
	}},
	{Image: "vkusvill.png", Keywords: []string{
		"vkusvill", "vkus vill", "вкусвилл",
	}},
	{Image: "lamoda.png", Keywords: []string{
		"lamoda club", "lamoda premium", "lamoda", "ламода",
	}},
	{Image: "shopify.png", Keywords: []string{
		"shopify", "shopify plus",
	}},
	{Image: "wix.png", Keywords: []string{
		"wix premium", "wix",
	}},
	{Image: "godaddy.png", Keywords: []string{
		"godaddy", "go daddy",
	}},
	{Image: "namecheap.png", Keywords: []string{
		"namecheap", "name cheap",
	}},
	{Image: "vultr.png", Keywords: []string{
		"vultr", "vultr cloud",
	}},
	{Image: "linode.png", Keywords: []string{
		"linode", "akamai cloud",
	}},
	{Image: "ovh.png", Keywords: []string{
		"ovh cloud", "ovh",
	}},
	{Image: "mailchimp.png", Keywords: []string{
		"mailchimp", "mail chimp",
	}},
	{Image: "salesforce.png", Keywords: []string{
		"salesforce", "sales force",
	}},
}

var (
	subscriptionNoiseSuffixes = []string{
		" subscription", " subscriptions", " premium", " pro", " plus", " family", " plan",
		" подписка", " подписки", " премиум", " про", " плюс", " семейная", " тариф",
	}
	subscriptionNoiseReplacer = strings.NewReplacer(
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
	wordBoundaryPattern = regexp.MustCompile(`[^a-z0-9а-я]+`)
)

func normalizeSubscriptionQuery(input string) string {
	normalized := strings.TrimSpace(strings.ToLower(subscriptionNoiseReplacer.Replace(input)))
	normalized = wordBoundaryPattern.ReplaceAllString(normalized, " ")
	normalized = strings.Join(strings.Fields(normalized), " ")

	for _, suffix := range subscriptionNoiseSuffixes {
		normalized = strings.TrimSuffix(normalized, suffix)
	}

	return strings.TrimSpace(normalized)
}

func keywordMatchScore(normalized, keyword string) int {
	kw := normalizeSubscriptionQuery(keyword)
	if kw == "" || normalized == "" {
		return 0
	}

	if normalized == kw {
		return 1000 + len(kw)
	}

	if !strings.Contains(normalized, kw) {
		return 0
	}

	if len(kw) <= 3 && !hasWordBoundary(normalized, kw) {
		return 0
	}

	if strings.HasPrefix(normalized, kw+" ") || strings.HasSuffix(normalized, " "+kw) || strings.Contains(normalized, " "+kw+" ") {
		return 800 + len(kw)
	}

	return 500 + len(kw)
}

func hasWordBoundary(normalized, keyword string) bool {
	if keyword == "" {
		return false
	}

	padded := " " + normalized + " "
	needle := " " + keyword + " "

	return strings.Contains(padded, needle)
}

func (nlp *NLPBuilder) GetSubscriptionImage(input string) string {
	if strings.TrimSpace(input) == "" {
		return ""
	}

	normalized := normalizeSubscriptionQuery(input)
	if normalized == "" {
		return ""
	}

	bestImage := ""
	bestScore := 0

	for _, brand := range subscriptionBrands {
		for _, keyword := range brand.Keywords {
			score := keywordMatchScore(normalized, keyword)
			if score > bestScore {
				bestScore = score
				bestImage = brand.Image
			}
		}
	}

	return bestImage
}
