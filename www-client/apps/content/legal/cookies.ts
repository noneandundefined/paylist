import type { LegalDocumentsByLang } from './types';

export const cookiesDocuments: LegalDocumentsByLang = {
	ru: {
		title: 'Политика использования cookie',
		updated: 'Последнее обновление: 16 августа 2026 г.',
		intro: 'Настоящая Политика описывает, какие cookie и аналогичные технологии использует веб-сайт https://paylist.site/ (далее — Сервис, Paylist), для каких целей они применяются и как Пользователь может управлять своим выбором.',
		callout:
			'Аналитические cookie (Яндекс Метрика) включаются только после согласия Пользователя. Строго необходимые cookie нужны для работы Сервиса и сохраняются независимо от выбора в баннере. Контакт: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Что такое cookie',
				paragraphs: [
					'1.1. Cookie — небольшие текстовые файлы, которые сайт сохраняет в браузере Пользователя. Аналогичные технологии (локальное хранилище браузера и пиксели) используются для тех же целей.',
					'1.2. Cookie помогают запомнить настройки, обеспечить вход в аккаунт, защитить Сервис и — при согласии Пользователя — понять, как им пользуются.',
				],
			},
			{
				heading: '2. Какие cookie мы используем',
				paragraphs: [
					'2.1. Строго необходимые cookie. Они нужны, чтобы Сервис работал: тема оформления, язык интерфейса, сессия авторизации и сам выбор согласия на cookie. Без них нельзя войти в аккаунт, сохранить тему и запомнить отказ или согласие на аналитику.',
					'2.2. Аналитические cookie. Мы используем Яндекс Метрику, чтобы считать посещения, просмотры страниц и понимать, какие разделы Сервиса востребованы. Эти cookie устанавливаются только после нажатия «Принять все».',
					'2.3. Мы не используем cookie рекламных сетей и не продаём данные Пользователя.',
				],
			},
			{
				heading: '3. Яндекс Метрика',
				paragraphs: [
					'3.1. Поставщик аналитики — ООО «ЯНДЕКС» (Яндекс Метрика). Сервис может передавать Метрике технические данные: IP-адрес в сокращённом виде, данные браузера, адрес посещённой страницы и действия на сайте.',
					'3.2. Обработка данных Яндексом осуществляется в соответствии с политикой конфиденциальности Яндекса: https://yandex.ru/legal/confidential/',
					'3.3. Запись сессий (вебвизор) в Paylist отключена.',
				],
			},
			{
				heading: '4. Как управлять согласием',
				paragraphs: [
					'4.1. При первом посещении Сервиса показывается баннер cookie. «Принять все» включает аналитику. «Отклонить все» оставляет только строго необходимые cookie.',
					'4.2. Выбор можно изменить в любой момент по ссылке «Настройки cookie» внизу страницы или в разделе профиля.',
					'4.3. Пользователь также может удалить cookie в настройках браузера. Отключение строго необходимых cookie может привести к некорректной работе Сервиса.',
				],
			},
			{
				heading: '5. Срок хранения',
				paragraphs: [
					'5.1. Выбор согласия на cookie хранится в локальном хранилище браузера, пока Пользователь его не изменит или не очистит данные сайта.',
					'5.2. Срок хранения cookie Яндекс Метрики определяется правилами Яндекса и обычно не превышает 13 месяцев.',
				],
			},
			{
				heading: '6. Контакты',
				paragraphs: ['По вопросам обработки данных и cookie можно обратиться к Оператору по адресу paylist.info@gmail.com. Политика обработки персональных данных: https://paylist.site/legal/privacy'],
			},
		],
	},
	en: {
		title: 'Cookie Policy',
		updated: 'Last updated: August 16, 2026',
		intro: 'This Policy describes which cookies and similar technologies the website https://paylist.site/ (the Service, Paylist) uses, why they are used, and how the User can manage their choice.',
		callout:
			'Analytics cookies (Yandex Metrica) are enabled only after the User’s consent. Strictly necessary cookies are required for the Service to work and are stored regardless of the banner choice. Contact: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. What cookies are',
				paragraphs: [
					'1.1. Cookies are small text files that a website stores in the User’s browser. Similar technologies (browser local storage and pixels) are used for the same purposes.',
					'1.2. Cookies help remember settings, keep the User signed in, protect the Service and — with the User’s consent — understand how the Service is used.',
				],
			},
			{
				heading: '2. Cookies we use',
				paragraphs: [
					'2.1. Strictly necessary cookies. They are required for the Service to work: appearance theme, interface language, auth session, and the cookie consent choice itself. Without them the User cannot sign in, keep a theme, or remember an analytics opt-in or opt-out.',
					'2.2. Analytics cookies. We use Yandex Metrica to count visits, page views, and understand which parts of the Service are used. These cookies are set only after the User clicks “Accept all”.',
					'2.3. We do not use advertising-network cookies and do not sell User data.',
				],
			},
			{
				heading: '3. Yandex Metrica',
				paragraphs: [
					'3.1. The analytics provider is YANDEX LLC (Yandex Metrica). The Service may send Metrica technical data: a truncated IP address, browser data, the visited page URL, and on-site actions.',
					'3.2. Yandex processes data in accordance with its privacy policy: https://yandex.com/legal/confidential/',
					'3.3. Session replay (Webvisor) is disabled in Paylist.',
				],
			},
			{
				heading: '4. How to manage consent',
				paragraphs: [
					'4.1. On the first visit the Service shows a cookie banner. “Accept all” enables analytics. “Reject all” keeps only strictly necessary cookies.',
					'4.2. The choice can be changed at any time via the “Cookie preferences” link at the bottom of the page or in the profile section.',
					'4.3. The User may also delete cookies in the browser settings. Disabling strictly necessary cookies may prevent the Service from working correctly.',
				],
			},
			{
				heading: '5. Retention',
				paragraphs: [
					'5.1. The cookie consent choice is stored in the browser’s local storage until the User changes it or clears site data.',
					'5.2. Retention of Yandex Metrica cookies is determined by Yandex rules and usually does not exceed 13 months.',
				],
			},
			{
				heading: '6. Contact',
				paragraphs: ['Questions about data processing and cookies can be sent to the Operator at paylist.info@gmail.com. Personal data processing policy: https://paylist.site/legal/privacy'],
			},
		],
	},
	de: {
		title: 'Cookie-Richtlinie',
		updated: 'Letzte Aktualisierung: 16. August 2026',
		intro: 'Diese Richtlinie beschreibt, welche Cookies und ähnlichen Technologien die Website https://paylist.site/ (der Dienst, Paylist) verwendet, zu welchen Zwecken sie eingesetzt werden und wie der Nutzer seine Wahl verwalten kann.',
		callout:
			'Analyse-Cookies (Yandex Metrica) werden nur nach Zustimmung des Nutzers aktiviert. Streng notwendige Cookies sind für den Betrieb des Dienstes erforderlich und werden unabhängig von der Banner-Auswahl gespeichert. Kontakt: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Was Cookies sind',
				paragraphs: [
					'1.1. Cookies sind kleine Textdateien, die eine Website im Browser des Nutzers speichert. Ähnliche Technologien (lokaler Browserspeicher und Pixel) werden für dieselben Zwecke verwendet.',
					'1.2. Cookies helfen, Einstellungen zu merken, die Anmeldung aufrechtzuerhalten, den Dienst zu schützen und — mit Zustimmung des Nutzers — die Nutzung des Dienstes zu verstehen.',
				],
			},
			{
				heading: '2. Welche Cookies wir verwenden',
				paragraphs: [
					'2.1. Streng notwendige Cookies. Sie sind für den Betrieb des Dienstes erforderlich: Erscheinungsbild, Sprache, Anmeldesitzung und die Cookie-Einwilligung selbst. Ohne sie kann sich der Nutzer nicht anmelden, ein Theme behalten oder eine Analyse-Zustimmung bzw. -Ablehnung speichern.',
					'2.2. Analyse-Cookies. Wir verwenden Yandex Metrica, um Besuche und Seitenaufrufe zu zählen und zu verstehen, welche Bereiche des Dienstes genutzt werden. Diese Cookies werden nur nach Klick auf „Alle akzeptieren“ gesetzt.',
					'2.3. Wir verwenden keine Cookies von Werbenetzwerken und verkaufen keine Nutzerdaten.',
				],
			},
			{
				heading: '3. Yandex Metrica',
				paragraphs: [
					'3.1. Anbieter der Analyse ist YANDEX LLC (Yandex Metrica). Der Dienst kann Metrica technische Daten übermitteln: eine gekürzte IP-Adresse, Browserdaten, die URL der besuchten Seite und Aktionen auf der Website.',
					'3.2. Yandex verarbeitet Daten gemäß seiner Datenschutzrichtlinie: https://yandex.com/legal/confidential/',
					'3.3. Sitzungsaufzeichnung (Webvisor) ist in Paylist deaktiviert.',
				],
			},
			{
				heading: '4. Einwilligung verwalten',
				paragraphs: [
					'4.1. Beim ersten Besuch zeigt der Dienst ein Cookie-Banner. „Alle akzeptieren“ aktiviert die Analyse. „Alle ablehnen“ belässt nur streng notwendige Cookies.',
					'4.2. Die Auswahl kann jederzeit über den Link „Cookie-Einstellungen“ am Seitenende oder im Profil geändert werden.',
					'4.3. Der Nutzer kann Cookies auch in den Browsereinstellungen löschen. Das Deaktivieren streng notwendiger Cookies kann den Dienst beeinträchtigen.',
				],
			},
			{
				heading: '5. Speicherdauer',
				paragraphs: [
					'5.1. Die Cookie-Einwilligung wird im lokalen Speicher des Browsers gespeichert, bis der Nutzer sie ändert oder die Website-Daten löscht.',
					'5.2. Die Speicherdauer der Cookies von Yandex Metrica richtet sich nach den Regeln von Yandex und beträgt in der Regel nicht mehr als 13 Monate.',
				],
			},
			{
				heading: '6. Kontakt',
				paragraphs: [
					'Fragen zur Datenverarbeitung und zu Cookies können an den Betreiber unter paylist.info@gmail.com gerichtet werden. Richtlinie zur Verarbeitung personenbezogener Daten: https://paylist.site/legal/privacy',
				],
			},
		],
	},
	es: {
		title: 'Política de cookies',
		updated: 'Última actualización: 16 de agosto de 2026',
		intro: 'Esta Política describe qué cookies y tecnologías similares utiliza el sitio web https://paylist.site/ (el Servicio, Paylist), con qué fines se emplean y cómo el Usuario puede gestionar su elección.',
		callout:
			'Las cookies de analítica (Yandex Metrica) se activan solo tras el consentimiento del Usuario. Las cookies estrictamente necesarias son precisas para el funcionamiento del Servicio y se guardan con independencia de la elección en el banner. Contacto: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Qué son las cookies',
				paragraphs: [
					'1.1. Las cookies son pequeños archivos de texto que un sitio web guarda en el navegador del Usuario. Tecnologías similares (almacenamiento local del navegador y píxeles) se usan con los mismos fines.',
					'1.2. Las cookies ayudan a recordar ajustes, mantener la sesión, proteger el Servicio y — con el consentimiento del Usuario — entender cómo se usa el Servicio.',
				],
			},
			{
				heading: '2. Cookies que utilizamos',
				paragraphs: [
					'2.1. Cookies estrictamente necesarias. Son precisas para que el Servicio funcione: tema de apariencia, idioma, sesión de autenticación y la propia elección de consentimiento de cookies. Sin ellas el Usuario no puede iniciar sesión, conservar el tema ni recordar la aceptación o el rechazo de la analítica.',
					'2.2. Cookies de analítica. Usamos Yandex Metrica para contar visitas, páginas vistas y entender qué partes del Servicio se utilizan. Estas cookies se establecen solo después de pulsar «Aceptar todo».',
					'2.3. No usamos cookies de redes publicitarias ni vendemos datos del Usuario.',
				],
			},
			{
				heading: '3. Yandex Metrica',
				paragraphs: [
					'3.1. El proveedor de analítica es YANDEX LLC (Yandex Metrica). El Servicio puede enviar a Metrica datos técnicos: una dirección IP truncada, datos del navegador, la URL de la página visitada y acciones en el sitio.',
					'3.2. Yandex trata los datos de conformidad con su política de privacidad: https://yandex.com/legal/confidential/',
					'3.3. La grabación de sesiones (Webvisor) está desactivada en Paylist.',
				],
			},
			{
				heading: '4. Cómo gestionar el consentimiento',
				paragraphs: [
					'4.1. En la primera visita el Servicio muestra un banner de cookies. «Aceptar todo» activa la analítica. «Rechazar todo» deja solo las cookies estrictamente necesarias.',
					'4.2. La elección se puede cambiar en cualquier momento mediante el enlace «Preferencias de cookies» al final de la página o en el perfil.',
					'4.3. El Usuario también puede eliminar las cookies en la configuración del navegador. Desactivar las cookies estrictamente necesarias puede impedir que el Servicio funcione correctamente.',
				],
			},
			{
				heading: '5. Conservación',
				paragraphs: [
					'5.1. La elección de consentimiento de cookies se guarda en el almacenamiento local del navegador hasta que el Usuario la cambie o borre los datos del sitio.',
					'5.2. El plazo de conservación de las cookies de Yandex Metrica lo determinan las normas de Yandex y normalmente no supera los 13 meses.',
				],
			},
			{
				heading: '6. Contacto',
				paragraphs: ['Las consultas sobre el tratamiento de datos y las cookies pueden enviarse al Operador a paylist.info@gmail.com. Política de tratamiento de datos personales: https://paylist.site/legal/privacy'],
			},
		],
	},
};
