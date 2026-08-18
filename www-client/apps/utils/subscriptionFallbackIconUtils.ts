export type SubscriptionFallbackKind = 'phone' | 'internet' | 'transport-car' | 'transport-train' | 'home' | 'utilities' | 'music' | 'streaming' | 'gaming' | 'cloud' | 'productivity' | 'fitness' | 'news' | 'education' | 'wallet';

type FallbackRule = {
	kind: SubscriptionFallbackKind;
	keywords: string[];
};

const normalize = (value: string): string =>
	value
		.toLowerCase()
		.replace(/ё/g, 'е')
		.replace(/[^\p{L}\p{N}]+/gu, ' ')
		.replace(/\s+/g, ' ')
		.trim();

const includesAny = (haystack: string, keywords: string[]): boolean => keywords.some((keyword) => haystack.includes(keyword));

const NAME_RULES: FallbackRule[] = [
	{
		kind: 'phone',
		keywords: [
			'мобильн',
			'сотов',
			'мегафон',
			'билайн',
			'теле2',
			'tele2',
			'yota',
			'йота',
			'сбермобайл',
			'тинькофф мобайл',
			'tinkoff mobile',
			'мотива',
			'kcell',
			'lifecell',
			'kyivstar',
			'beeline',
			'megafon',
			't mobile',
			'vodafone',
		],
	},
	{ kind: 'internet', keywords: ['интернет', 'wifi', 'wi fi', 'провайдер', 'роутер', 'broadband', 'дом ру', 'домру', 'ростелеком', 'мтс домашний', 'домашний интернет', 'vpn', 'впн'] },
	{ kind: 'cloud', keywords: ['облак', 'cloud', 'icloud', 'google one', 'яндекс диск', 'dropbox', 'vps', 'vds', 'впс', 'вдс', 'хостинг', 'hosting'] },
	{ kind: 'transport-train', keywords: ['метро', 'поезд', 'ржд', 'трамвай', 'автобус', 'троллейбус', 'проездн', 'транспортн карт', 'тройка', 'подорожник'] },
	{ kind: 'transport-car', keywords: ['транспорт', 'такси', 'uber', 'яндекс го', 'каршеринг', 'ситимобил', 'парковк', 'бензин', 'топлив', 'азс', 'parking', 'carsharing', 'citymobil', 'платон', 'осаго', 'каско'] },
	{ kind: 'home', keywords: ['квартир', 'аренд', 'жилье', 'ипотека', 'apartment', 'mortgage', 'rent', 'капремонт'] },
	{ kind: 'utilities', keywords: ['жкх', 'коммунал', 'электричеств', 'электроэнерг', 'водоканал', 'отоплен', 'газоснаб', 'мусор', 'тко', 'домофон', 'intercom', 'мосэнерго', 'теплоэнерго', 'квартплат', 'utility', 'utilities'] },
	{ kind: 'music', keywords: ['музык', 'spotify', 'яндекс музык', 'apple music', 'youtube music'] },
	{ kind: 'streaming', keywords: ['стриминг', 'кино', 'сериал', 'онлайн кино', 'youtube', 'netflix', 'ivi', 'okko', 'wink', 'кион', 'media station', 'медиа стейшн'] },
	{ kind: 'gaming', keywords: ['игр', 'gaming', 'xbox', 'playstation', 'steam', 'game pass'] },
	{ kind: 'fitness', keywords: ['фитнес', 'спортзал', 'тренажер', 'yoga', 'йога', 'gym', 'workout'] },
	{ kind: 'news', keywords: ['новост', 'газет', 'журнал', 'news'] },
	{ kind: 'education', keywords: ['образован', 'курс', 'школ', 'универ', 'skillbox', 'geekbrains', 'education'] },
	{ kind: 'productivity', keywords: ['офис', 'notion', 'chatgpt', 'openai', 'adobe', 'figma', 'productivity'] },
	{ kind: 'phone', keywords: ['связь', 'мтс'] },
];

const CATEGORY_KIND: Record<string, SubscriptionFallbackKind> = {
	subscriptions: 'wallet',
	music: 'music',
	streaming: 'streaming',
	gaming: 'gaming',
	cloud: 'cloud',
	productivity: 'productivity',
	fitness: 'fitness',
	news: 'news',
	education: 'education',
	utilities: 'utilities',
};

export const getSubscriptionFallbackKind = (name: string, categories: string[] = []): SubscriptionFallbackKind => {
	const haystack = normalize(name);

	if (haystack) {
		for (const rule of NAME_RULES) {
			if (includesAny(haystack, rule.keywords)) {
				return rule.kind;
			}
		}
	}

	for (const category of categories) {
		const kind = CATEGORY_KIND[normalize(category)];
		if (kind) {
			return kind;
		}
	}

	return 'wallet';
};
