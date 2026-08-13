import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import PDFDocument from 'pdfkit';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const rootDir = path.resolve(__dirname, '..');
const outputDir = path.join(rootDir, 'public', 'docs', 'legal');
const fontPath = path.join(rootDir, 'node_modules', 'dejavu-fonts-ttf', 'ttf', 'DejaVuSans.ttf');

const languages = ['ru', 'en', 'de', 'es'];

const termsContent = {
	en: {
		title: 'Paylist Terms of Service',
		updated: 'Last updated: August 13, 2026',
		sections: [
			{
				heading: '1. Agreement',
				paragraphs: [
					'These Terms of Service ("Terms") govern your access to and use of the Paylist web application and related services ("Service"). By creating an account or using the Service, you agree to these Terms.',
					'If you do not agree, do not use the Service.',
				],
			},
			{
				heading: '2. Service description',
				paragraphs: [
					'Paylist helps users track subscriptions, view spending analytics, and receive optional payment reminders. Features may differ between free and premium plans.',
					'We may update, suspend, or discontinue parts of the Service at any time.',
				],
			},
			{
				heading: '3. Accounts',
				paragraphs: [
					'You must provide accurate registration information and keep your credentials secure. You are responsible for all activity under your account.',
					'You must be at least 16 years old, or the minimum age required in your jurisdiction, to use the Service.',
				],
			},
			{
				heading: '4. Subscriptions and payments',
				paragraphs: [
					'Premium features may require a paid plan. Prices, billing periods, and included features are shown before purchase.',
					'Payments are processed by third-party payment providers. Refunds are handled according to applicable law and our billing policy shown at checkout.',
				],
			},
			{
				heading: '5. Acceptable use',
				paragraphs: [
					'You may not misuse the Service, attempt unauthorized access, interfere with infrastructure, or use the Service for unlawful purposes.',
					'We may suspend or terminate accounts that violate these Terms.',
				],
			},
			{
				heading: '6. Intellectual property',
				paragraphs: [
					'Paylist, its design, software, and branding remain our property or that of our licensors. You retain ownership of data you enter into the Service.',
					'You grant us a limited license to host and process your data solely to operate the Service.',
				],
			},
			{
				heading: '7. Disclaimer and liability',
				paragraphs: [
					'The Service is provided "as is" without warranties of uninterrupted or error-free operation. Spending forecasts and reminders are informational only.',
					'To the maximum extent permitted by law, Paylist is not liable for indirect, incidental, or consequential damages arising from use of the Service.',
				],
			},
			{
				heading: '8. Changes and contact',
				paragraphs: ['We may revise these Terms. Material changes will be communicated through the Service or by email when appropriate.', 'Questions about these Terms: support@paylist.app'],
			},
		],
	},
	ru: {
		title: 'Условия использования Paylist',
		updated: 'Последнее обновление: 13 августа 2026 г.',
		sections: [
			{
				heading: '1. Соглашение',
				paragraphs: [
					'Настоящие Условия использования («Условия») регулируют доступ к веб-приложению Paylist и связанным сервисам («Сервис»). Создавая аккаунт или используя Сервис, вы принимаете эти Условия.',
					'Если вы не согласны с Условиями, не используйте Сервис.',
				],
			},
			{
				heading: '2. Описание сервиса',
				paragraphs: [
					'Paylist помогает отслеживать подписки, анализировать расходы и получать необязательные напоминания об оплате. Набор функций может отличаться для бесплатного и Premium-тарифов.',
					'Мы можем обновлять, приостанавливать или прекращать работу отдельных функций Сервиса.',
				],
			},
			{
				heading: '3. Аккаунт',
				paragraphs: [
					'Вы обязуетесь указывать достоверные данные при регистрации и обеспечивать безопасность учётных данных. Вы несёте ответственность за действия в своём аккаунте.',
					'Использовать Сервис могут пользователи не моложе 16 лет либо минимального возраста, установленного законодательством вашей страны.',
				],
			},
			{
				heading: '4. Подписки и оплата',
				paragraphs: [
					'Premium-функции могут требовать оплачиваемого тарифа. Стоимость, период списания и состав функций отображаются до оплаты.',
					'Платежи обрабатываются сторонними платёжными провайдерами. Возвраты осуществляются в соответствии с применимым законодательством и правилами, указанными при оплате.',
				],
			},
			{
				heading: '5. Допустимое использование',
				paragraphs: [
					'Запрещается злоупотреблять Сервисом, пытаться получить несанкционированный доступ, нарушать работу инфраструктуры или использовать Сервис в незаконных целях.',
					'Мы можем ограничить или удалить аккаунт при нарушении Условий.',
				],
			},
			{
				heading: '6. Интеллектуальная собственность',
				paragraphs: [
					'Paylist, его дизайн, программное обеспечение и бренд принадлежат нам или нашим лицензиарам. Вы сохраняете права на данные, которые вводите в Сервис.',
					'Вы предоставляете нам ограниченную лицензию на хранение и обработку данных исключительно для работы Сервиса.',
				],
			},
			{
				heading: '7. Отказ от гарантий и ответственность',
				paragraphs: [
					'Сервис предоставляется «как есть» без гарантий бесперебойной или безошибочной работы. Прогнозы расходов и напоминания носят информационный характер.',
					'В пределах, допустимых законом, Paylist не несёт ответственности за косвенные или сопутствующие убытки, связанные с использованием Сервиса.',
				],
			},
			{
				heading: '8. Изменения и контакты',
				paragraphs: ['Мы можем изменять Условия. О существенных изменениях сообщим через Сервис или по email, когда это уместно.', 'Вопросы по Условиям: support@paylist.app'],
			},
		],
	},
	de: {
		title: 'Paylist Nutzungsbedingungen',
		updated: 'Zuletzt aktualisiert: 13. August 2026',
		sections: [
			{
				heading: '1. Vereinbarung',
				paragraphs: [
					'Diese Nutzungsbedingungen ("Bedingungen") regeln Ihren Zugang zur Paylist-Webanwendung und zu den zugehörigen Diensten ("Dienst"). Mit der Kontoerstellung oder Nutzung stimmen Sie diesen Bedingungen zu.',
					'Wenn Sie nicht einverstanden sind, nutzen Sie den Dienst nicht.',
				],
			},
			{
				heading: '2. Leistungsbeschreibung',
				paragraphs: [
					'Paylist hilft dabei, Abonnements zu verfolgen, Ausgaben zu analysieren und optionale Zahlungserinnerungen zu erhalten. Funktionen können je nach kostenlosem oder Premium-Tarif variieren.',
					'Wir können Teile des Dienstes jederzeit aktualisieren, aussetzen oder einstellen.',
				],
			},
			{
				heading: '3. Konten',
				paragraphs: [
					'Sie müssen korrekte Registrierungsdaten angeben und Ihre Zugangsdaten schützen. Sie sind für alle Aktivitäten unter Ihrem Konto verantwortlich.',
					'Sie müssen mindestens 16 Jahre alt sein oder das in Ihrem Land geltende Mindestalter erreicht haben.',
				],
			},
			{
				heading: '4. Abonnements und Zahlungen',
				paragraphs: [
					'Premium-Funktionen können einen kostenpflichtigen Tarif erfordern. Preise, Abrechnungszeiträume und enthaltene Funktionen werden vor dem Kauf angezeigt.',
					'Zahlungen werden von Drittanbietern verarbeitet. Erstattungen erfolgen gemäß geltendem Recht und den beim Checkout angezeigten Regeln.',
				],
			},
			{
				heading: '5. Zulässige Nutzung',
				paragraphs: [
					'Sie dürfen den Dienst nicht missbrauchen, unbefugten Zugriff versuchen, die Infrastruktur stören oder den Dienst für rechtswidrige Zwecke nutzen.',
					'Wir können Konten bei Verstößen gegen diese Bedingungen sperren oder löschen.',
				],
			},
			{
				heading: '6. Geistiges Eigentum',
				paragraphs: [
					'Paylist, sein Design, seine Software und Marke bleiben unser Eigentum oder das unserer Lizenzgeber. Sie behalten das Eigentum an den von Ihnen eingegebenen Daten.',
					'Sie gewähren uns eine begrenzte Lizenz zur Speicherung und Verarbeitung Ihrer Daten ausschließlich zum Betrieb des Dienstes.',
				],
			},
			{
				heading: '7. Haftungsausschluss',
				paragraphs: [
					'Der Dienst wird "wie besehen" bereitgestellt. Prognosen und Erinnerungen dienen nur der Information.',
					'Soweit gesetzlich zulässig, haftet Paylist nicht für indirekte oder Folgeschäden aus der Nutzung des Dienstes.',
				],
			},
			{
				heading: '8. Änderungen und Kontakt',
				paragraphs: ['Wir können diese Bedingungen ändern. Wesentliche Änderungen werden über den Dienst oder per E-Mail mitgeteilt.', 'Fragen: support@paylist.app'],
			},
		],
	},
	es: {
		title: 'Términos de servicio de Paylist',
		updated: 'Última actualización: 13 de agosto de 2026',
		sections: [
			{
				heading: '1. Acuerdo',
				paragraphs: [
					'Estos Términos de servicio ("Términos") regulan el acceso y uso de la aplicación web Paylist y servicios relacionados ("Servicio"). Al crear una cuenta o usar el Servicio, acepta estos Términos.',
					'Si no está de acuerdo, no utilice el Servicio.',
				],
			},
			{
				heading: '2. Descripción del servicio',
				paragraphs: [
					'Paylist ayuda a hacer seguimiento de suscripciones, analizar gastos y recibir recordatorios de pago opcionales. Las funciones pueden variar entre planes gratuitos y premium.',
					'Podemos actualizar, suspender o discontinuar partes del Servicio en cualquier momento.',
				],
			},
			{
				heading: '3. Cuentas',
				paragraphs: [
					'Debe proporcionar información de registro precisa y mantener seguras sus credenciales. Usted es responsable de toda actividad en su cuenta.',
					'Debe tener al menos 16 años o la edad mínima exigida en su jurisdicción.',
				],
			},
			{
				heading: '4. Suscripciones y pagos',
				paragraphs: [
					'Las funciones premium pueden requerir un plan de pago. Precios, periodos de facturación y funciones incluidas se muestran antes de la compra.',
					'Los pagos son procesados por proveedores externos. Los reembolsos se gestionan conforme a la ley aplicable y la política mostrada al pagar.',
				],
			},
			{
				heading: '5. Uso aceptable',
				paragraphs: [
					'No debe hacer un uso indebido del Servicio, intentar accesos no autorizados, interferir con la infraestructura ni usar el Servicio con fines ilícitos.',
					'Podemos suspender o eliminar cuentas que infrinjan estos Términos.',
				],
			},
			{
				heading: '6. Propiedad intelectual',
				paragraphs: [
					'Paylist, su diseño, software y marca son nuestra propiedad o la de nuestros licenciantes. Usted conserva la propiedad de los datos que introduce en el Servicio.',
					'Nos otorga una licencia limitada para alojar y procesar sus datos únicamente para operar el Servicio.',
				],
			},
			{
				heading: '7. Exención de responsabilidad',
				paragraphs: [
					'El Servicio se proporciona "tal cual". Las previsiones y recordatorios son solo informativos.',
					'En la medida permitida por la ley, Paylist no será responsable de daños indirectos o consecuentes derivados del uso del Servicio.',
				],
			},
			{
				heading: '8. Cambios y contacto',
				paragraphs: ['Podemos revisar estos Términos. Los cambios importantes se comunicarán a través del Servicio o por correo electrónico cuando corresponda.', 'Consultas: support@paylist.app'],
			},
		],
	},
};

const privacyContent = {
	en: {
		title: 'Paylist Privacy Policy',
		updated: 'Last updated: August 13, 2026',
		sections: [
			{
				heading: '1. Overview',
				paragraphs: ['This Privacy Policy explains how Paylist ("we", "us") collects, uses, and protects personal data when you use our Service.', 'By using Paylist, you agree to the practices described in this Policy.'],
			},
			{
				heading: '2. Data we collect',
				paragraphs: [
					'Account data: email address, name (if provided), authentication identifiers, and subscription plan status.',
					'Usage data: subscriptions you add (names, prices, currencies, billing dates, categories, notes), app settings, and language preferences.',
					'Technical data: device/session identifiers, IP address, browser type, and logs needed for security and troubleshooting.',
					'Payment data: processed by payment providers; we store billing references, not full card numbers.',
				],
			},
			{
				heading: '3. How we use data',
				paragraphs: ['We use data to provide and improve the Service, calculate analytics, send optional reminders, process premium payments, and respond to support requests.', 'We do not sell your personal data.'],
			},
			{
				heading: '4. Notifications',
				paragraphs: [
					"If you enable reminders, we may send notifications via Telegram or email based on your account settings and each subscription's notification preference.",
					'You can disable notifications in profile settings or per subscription.',
				],
			},
			{
				heading: '5. Storage and security',
				paragraphs: [
					'Data is stored on secure servers with access controls. We apply reasonable technical and organizational measures to protect information.',
					'No method of transmission or storage is completely secure; please use a strong password.',
				],
			},
			{
				heading: '6. Third parties',
				paragraphs: [
					'We may share limited data with infrastructure, payment, email, and messaging providers strictly to operate the Service.',
					'Third parties are required to protect data according to contractual and legal obligations.',
				],
			},
			{
				heading: '7. Your rights',
				paragraphs: [
					'Depending on your location, you may request access, correction, deletion, or export of your personal data, and object to certain processing.',
					'Contact support@paylist.app to exercise these rights. You may delete your account in profile settings.',
				],
			},
			{
				heading: '8. Changes and contact',
				paragraphs: ['We may update this Policy. Material changes will be communicated through the Service when appropriate.', 'Privacy questions: support@paylist.app'],
			},
		],
	},
	ru: {
		title: 'Политика конфиденциальности Paylist',
		updated: 'Последнее обновление: 13 августа 2026 г.',
		sections: [
			{
				heading: '1. Общие положения',
				paragraphs: [
					'Настоящая Политика конфиденциальности описывает, как Paylist («мы») собирает, использует и защищает персональные данные при использовании Сервиса.',
					'Используя Paylist, вы соглашаетесь с практиками, описанными в этой Политике.',
				],
			},
			{
				heading: '2. Какие данные мы собираем',
				paragraphs: [
					'Данные аккаунта: email, имя (если указано), идентификаторы авторизации и статус тарифного плана.',
					'Данные использования: добавленные подписки (названия, цены, валюты, даты оплаты, категории, заметки), настройки приложения и язык интерфейса.',
					'Технические данные: идентификаторы сессий/устройств, IP-адрес, тип браузера и журналы для безопасности и диагностики.',
					'Платёжные данные обрабатываются платёжными провайдерами; мы храним ссылки на платежи, но не полные номера карт.',
				],
			},
			{
				heading: '3. Как мы используем данные',
				paragraphs: [
					'Данные используются для работы и улучшения Сервиса, расчёта аналитики, отправки необязательных напоминаний, обработки Premium-платежей и поддержки пользователей.',
					'Мы не продаём ваши персональные данные.',
				],
			},
			{
				heading: '4. Уведомления',
				paragraphs: [
					'При включении напоминаний мы можем отправлять уведомления в Telegram или по email согласно настройкам профиля и каждой подписки.',
					'Уведомления можно отключить в настройках профиля или для отдельной подписки.',
				],
			},
			{
				heading: '5. Хранение и безопасность',
				paragraphs: [
					'Данные хранятся на защищённых серверах с контролем доступа. Мы применяем разумные технические и организационные меры защиты.',
					'Ни один способ передачи или хранения не является абсолютно безопасным; используйте надёжный пароль.',
				],
			},
			{
				heading: '6. Третьи стороны',
				paragraphs: [
					'Мы можем передавать ограниченный набор данных инфраструктурным, платёжным, почтовым и мессенджер-провайдерам исключительно для работы Сервиса.',
					'Третьи стороны обязаны защищать данные согласно договорным и правовым требованиям.',
				],
			},
			{
				heading: '7. Ваши права',
				paragraphs: [
					'В зависимости от вашей юрисдикции вы можете запросить доступ, исправление, удаление или экспорт персональных данных, а также возразить против отдельных видов обработки.',
					'Для этого напишите на support@paylist.app. Аккаунт можно удалить в настройках профиля.',
				],
			},
			{
				heading: '8. Изменения и контакты',
				paragraphs: ['Мы можем обновлять Политику. О существенных изменениях сообщим через Сервис, когда это уместно.', 'Вопросы по конфиденциальности: support@paylist.app'],
			},
		],
	},
	de: {
		title: 'Paylist Datenschutzrichtlinie',
		updated: 'Zuletzt aktualisiert: 13. August 2026',
		sections: [
			{
				heading: '1. Überblick',
				paragraphs: [
					'Diese Datenschutzrichtlinie erklärt, wie Paylist ("wir") personenbezogene Daten beim Nutzen des Dienstes erhebt, verwendet und schützt.',
					'Mit der Nutzung von Paylist stimmen Sie den beschriebenen Praktiken zu.',
				],
			},
			{
				heading: '2. Erhobene Daten',
				paragraphs: [
					'Kontodaten: E-Mail-Adresse, Name (falls angegeben), Authentifizierungskennungen und Tarifstatus.',
					'Nutzungsdaten: von Ihnen hinzugefügte Abonnements (Namen, Preise, Währungen, Zahlungstermine, Kategorien, Notizen), App-Einstellungen und Sprache.',
					'Technische Daten: Sitzungs-/Gerätekennungen, IP-Adresse, Browsertyp und Protokolle für Sicherheit und Fehlerbehebung.',
					'Zahlungsdaten werden von Zahlungsanbietern verarbeitet; wir speichern Abrechnungsreferenzen, keine vollständigen Kartennummern.',
				],
			},
			{
				heading: '3. Verwendung der Daten',
				paragraphs: ['Wir verwenden Daten zur Bereitstellung und Verbesserung des Dienstes, für Analysen, optionale Erinnerungen, Premium-Abrechnung und Support.', 'Wir verkaufen Ihre personenbezogenen Daten nicht.'],
			},
			{
				heading: '4. Benachrichtigungen',
				paragraphs: [
					'Bei aktivierten Erinnerungen können wir Benachrichtigungen per Telegram oder E-Mail gemäß Kontoeinstellungen und Abo-Einstellungen senden.',
					'Sie können Benachrichtigungen im Profil oder pro Abonnement deaktivieren.',
				],
			},
			{
				heading: '5. Speicherung und Sicherheit',
				paragraphs: [
					'Daten werden auf geschützten Servern mit Zugriffskontrollen gespeichert. Wir setzen angemessene technische und organisatorische Schutzmaßnahmen ein.',
					'Keine Übertragungs- oder Speichermethode ist vollständig sicher; verwenden Sie ein starkes Passwort.',
				],
			},
			{
				heading: '6. Drittanbieter',
				paragraphs: [
					'Wir können begrenzte Daten an Infrastruktur-, Zahlungs-, E-Mail- und Messaging-Anbieter ausschließlich zum Betrieb des Dienstes weitergeben.',
					'Drittanbieter sind vertraglich und gesetzlich zum Schutz der Daten verpflichtet.',
				],
			},
			{
				heading: '7. Ihre Rechte',
				paragraphs: [
					'Je nach Standort können Sie Auskunft, Berichtigung, Löschung oder Export personenbezogener Daten verlangen und bestimmter Verarbeitung widersprechen.',
					'Kontakt: support@paylist.app. Sie können Ihr Konto in den Profileinstellungen löschen.',
				],
			},
			{
				heading: '8. Änderungen und Kontakt',
				paragraphs: ['Wir können diese Richtlinie aktualisieren. Wesentliche Änderungen werden über den Dienst mitgeteilt.', 'Datenschutzfragen: support@paylist.app'],
			},
		],
	},
	es: {
		title: 'Política de privacidad de Paylist',
		updated: 'Última actualización: 13 de agosto de 2026',
		sections: [
			{
				heading: '1. Resumen',
				paragraphs: [
					'Esta Política de privacidad explica cómo Paylist ("nosotros") recopila, usa y protege datos personales cuando utiliza el Servicio.',
					'Al usar Paylist, acepta las prácticas descritas en esta Política.',
				],
			},
			{
				heading: '2. Datos que recopilamos',
				paragraphs: [
					'Datos de cuenta: correo electrónico, nombre (si se proporciona), identificadores de autenticación y estado del plan.',
					'Datos de uso: suscripciones añadidas (nombres, precios, monedas, fechas de pago, categorías, notas), ajustes de la app e idioma.',
					'Datos técnicos: identificadores de sesión/dispositivo, dirección IP, tipo de navegador y registros para seguridad y diagnóstico.',
					'Los datos de pago los procesan proveedores externos; almacenamos referencias de facturación, no números completos de tarjeta.',
				],
			},
			{
				heading: '3. Uso de los datos',
				paragraphs: ['Usamos los datos para prestar y mejorar el Servicio, calcular analíticas, enviar recordatorios opcionales, procesar pagos premium y brindar soporte.', 'No vendemos sus datos personales.'],
			},
			{
				heading: '4. Notificaciones',
				paragraphs: [
					'Si activa recordatorios, podemos enviar notificaciones por Telegram o correo según la configuración de la cuenta y de cada suscripción.',
					'Puede desactivar las notificaciones en el perfil o por suscripción.',
				],
			},
			{
				heading: '5. Almacenamiento y seguridad',
				paragraphs: [
					'Los datos se almacenan en servidores seguros con controles de acceso. Aplicamos medidas técnicas y organizativas razonables.',
					'Ningún método de transmisión o almacenamiento es totalmente seguro; use una contraseña robusta.',
				],
			},
			{
				heading: '6. Terceros',
				paragraphs: [
					'Podemos compartir datos limitados con proveedores de infraestructura, pagos, correo y mensajería únicamente para operar el Servicio.',
					'Los terceros deben proteger los datos según obligaciones contractuales y legales.',
				],
			},
			{
				heading: '7. Sus derechos',
				paragraphs: [
					'Según su ubicación, puede solicitar acceso, corrección, eliminación o exportación de datos personales y oponerse a ciertos tratamientos.',
					'Contacto: support@paylist.app. Puede eliminar su cuenta en ajustes del perfil.',
				],
			},
			{
				heading: '8. Cambios y contacto',
				paragraphs: ['Podemos actualizar esta Política. Los cambios importantes se comunicarán a través del Servicio cuando corresponda.', 'Consultas de privacidad: support@paylist.app'],
			},
		],
	},
};

const writePdf = (filePath, document) =>
	new Promise((resolve, reject) => {
		fs.mkdirSync(path.dirname(filePath), { recursive: true });

		const pdf = new PDFDocument({ margin: 56, size: 'A4' });
		const stream = fs.createWriteStream(filePath);

		stream.on('finish', () => resolve());
		stream.on('error', reject);
		pdf.on('error', reject);

		pdf.pipe(stream);
		pdf.font(fontPath);

		pdf.fontSize(20).text(document.title, { align: 'left' });
		pdf.moveDown(0.5);
		pdf.fontSize(10).fillColor('#666666').text(document.updated);
		pdf.moveDown(1.2);
		pdf.fillColor('#000000');

		for (const section of document.sections) {
			pdf.fontSize(13).text(section.heading, { continued: false });
			pdf.moveDown(0.4);

			for (const paragraph of section.paragraphs) {
				pdf.fontSize(11).text(paragraph, {
					align: 'left',
					lineGap: 4,
				});
				pdf.moveDown(0.6);
			}

			pdf.moveDown(0.4);
		}

		pdf.end();
	});

if (!fs.existsSync(fontPath)) {
	console.error(`Font not found: ${fontPath}. Run npm install first.`);
	process.exit(1);
}

for (const lang of languages) {
	await writePdf(path.join(outputDir, `terms-${lang}.pdf`), termsContent[lang]);
	await writePdf(path.join(outputDir, `privacy-${lang}.pdf`), privacyContent[lang]);
	console.log(`Generated legal PDFs for ${lang}`);
}

console.log(`Saved PDFs to ${outputDir}`);
