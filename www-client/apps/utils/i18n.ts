import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import en from '@/locale/languages/en-us.json';
import ru from '@/locale/languages/ru-ru.json';
import de from '@/locale/languages/de-de.json';
import es from '@/locale/languages/es-es.json';

import { DEFAULT_LANGUAGE, getAppLanguage } from '@/constants/Language.constant';

const resources = {
	en: { translation: en },
	ru: { translation: ru },
	de: { translation: de },
	es: { translation: es },
};

const savedLang = localStorage.getItem('lang');
const browserLang = navigator.language.slice(0, 2);
i18n.use(initReactI18next).init({
	resources,
	lng: savedLang ? getAppLanguage(savedLang) : getAppLanguage(browserLang),
	fallbackLng: DEFAULT_LANGUAGE,
	interpolation: {
		escapeValue: false,
	},
});

export default i18n;
