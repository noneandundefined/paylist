import type { LegalDocumentsByLang } from './types';

export const offerDocuments: LegalDocumentsByLang = {
	ru: {
		title: 'Публичная оферта на автоплатежи',
		updated: 'Последнее обновление: 16 августа 2026 г.',
		intro: 'Настоящая публичная оферта (далее — Оферта) определяет условия автоматических периодических списаний («Автоплатежи») за тариф Paylist Premium, предоставляемый Власовым Артёмом Владимировичем (далее — Оператор) через веб-сайт https://paylist.site/',
		callout:
			'Отмечая согласие на странице тарифов и совершая первый платёж, Пользователь принимает Оферту. Автоплатежи можно отключить в разделе «Аккаунт» — кнопка «Отмена подписки». Актуальная версия: https://paylist.site/legal/offer. Контакт: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Общие положения',
				paragraphs: [
					'1.1. Оферта является публичной офертой в смысле статей 435 и 437 Гражданского кодекса Российской Федерации и адресована неопределённому кругу лиц.',
					'1.2. Оферта применяется совместно с Пользовательским соглашением https://paylist.site/legal/terms и Политикой в отношении обработки персональных данных https://paylist.site/legal/privacy.',
					'1.3. Оператор вправе изменять Оферту. Новая редакция вступает в силу с момента размещения по адресу https://paylist.site/legal/offer, если иной срок не указан в самой редакции.',
					'1.4. Платежи обрабатывает ЮKassa (ООО НКО «ЮМани»). Оператор не хранит полный номер банковской карты Пользователя.',
				],
			},
			{
				heading: '2. Основные понятия, используемые в Оферте',
				paragraphs: [
					'2.1. Оператор — Власов Артём Владимирович, лицо, предоставляющее сервис Paylist.',
					'2.2. Пользователь — физическое лицо, принявшее Оферту.',
					'2.3. Веб-сайт — сайт Оператора в сети интернет по адресу https://paylist.site/',
					'2.4. Сервис (Paylist) — веб-приложение Оператора для учёта подписок, аналитики расходов и напоминаний об оплате.',
					'2.5. Тариф Premium — возмездный тариф Сервиса, стоимость, период и состав функций которого отображаются на странице тарифов до оплаты.',
					'2.6. Автоплатежи — автоматические периодические списания денежных средств за тариф Premium с сохранённого способа оплаты Пользователя без дополнительного подтверждения каждого списания.',
					'2.7. Акцепт — полное и безоговорочное принятие Оферты путём проставления отметки о согласии на странице тарифов и совершения первого платежа за тариф Premium.',
					'2.8. Платёжный провайдер — ЮKassa (ООО НКО «ЮМани»), осуществляющий приём и обработку платежей.',
				],
			},
			{
				heading: '3. Предмет Оферты',
				paragraphs: [
					'3.1. Оператор обязуется предоставлять Пользователю доступ к тарифу Premium, а Пользователь обязуется оплачивать тариф Premium в порядке, установленном Офертой.',
					'3.2. Оферта регулирует исключительно Автоплатежи за тариф Premium. Использование Сервиса в остальной части регулируется Пользовательским соглашением.',
				],
			},
			{
				heading: '4. Акцепт Оферты',
				paragraphs: [
					'4.1. Акцептом Оферты признаются совокупность следующих действий Пользователя: ознакомление с Офертой; проставление отметки о согласии на странице тарифов; совершение первого платежа за тариф Premium.',
					'4.2. С момента Акцепта между Оператором и Пользователем считается заключённым договор на условиях Оферты.',
					'4.3. Совершая Акцепт, Пользователь подтверждает, что ему понятны сумма, периодичность и порядок отключения Автоплатежей.',
				],
			},
			{
				heading: '5. Подключение Автоплатежей',
				paragraphs: [
					'5.1. При оформлении тарифа Premium первый платёж проходит через Платёжного провайдера.',
					'5.2. Пользователь соглашается на сохранение способа оплаты у Платёжного провайдера для последующих списаний.',
					'5.3. После успешного первого платежа Автоплатежи включаются в Аккаунте Пользователя.',
				],
			},
			{
				heading: '6. Сумма и периодичность списаний',
				paragraphs: [
					'6.1. Сумма списания и период оплаты показываются на странице тарифов до оплаты и на карточке Premium в разделе «Аккаунт».',
					'6.2. Тариф Premium списывается один раз за период тарифа (как правило, каждые 30 дней) в размере, указанном при оформлении.',
					'6.3. Дата следующего списания — дата окончания оплаченного периода Premium, отображаемая в Аккаунте. Списание выполняется в 21:00 по московскому времени в день списания.',
					'6.4. Если Автоплатежи не отключены, в дату списания списывается та же сумма, и доступ к тарифу Premium продлевается на следующий период.',
				],
			},
			{
				heading: '7. Отключение Автоплатежей',
				paragraphs: [
					'7.1. Пользователь вправе отключить Автоплатежи в любой момент в разделе «Аккаунт»: на карточке Premium необходимо нажать «Отмена подписки».',
					'7.2. После отключения новые автоматические списания не выполняются.',
					'7.3. Тариф Premium остаётся активным до уже оплаченной даты окончания периода и далее не продлевается, если Пользователь не оформит его повторно.',
					'7.4. Повторное включение Автоплатежей возможно путём нового оформления тарифа Premium и Акцепта Оферты.',
				],
			},
			{
				heading: '8. Неуспешные списания',
				paragraphs: [
					'8.1. Если очередное списание не прошло, Оператор вправе повторить попытку списания либо приостановить Автоплатежи.',
					'8.2. Если оплаченный период закончился без успешного платежа, Аккаунт переводится на тариф Free.',
					'8.3. Оператор не несёт ответственности за отказ банка или Платёжного провайдера в проведении платежа по причинам, не зависящим от Оператора (недостаточно средств, ограничения банка, истечение срока карты и т.п.).',
				],
			},
			{
				heading: '9. Права и обязанности Сторон',
				paragraphs: [
					'9.1. Пользователь вправе получать доступ к тарифу Premium при успешной оплате, отключать Автоплатежи в порядке раздела 7 и обращаться к Оператору по вопросам оплаты.',
					'9.2. Пользователь обязан обеспечивать достаточный остаток средств и актуальность сохранённого способа оплаты, а также самостоятельно следить за суммой и датой следующего списания в Аккаунте.',
					'9.3. Оператор вправе привлекать Платёжного провайдера для приёма платежей и приостанавливать Автоплатежи при нарушении Оферты.',
					'9.4. Оператор обязан отображать сумму и период списания до оплаты и предоставлять возможность отключения Автоплатежей в Аккаунте.',
				],
			},
			{
				heading: '10. Ответственность',
				paragraphs: [
					'10.1. Оператор не хранит полный номер карты и не несёт ответственность за обработку платёжных данных Платёжным провайдером.',
					'10.2. Вся информация, собираемая Платёжным провайдером, обрабатывается им в соответствии с его правилами и политикой конфиденциальности. Пользователь обязан самостоятельно ознакомиться с указанными документами.',
					'10.3. В пределах, допустимых законодательством Российской Федерации, Оператор не несёт ответственности за косвенные убытки, связанные с Автоплатежами.',
				],
			},
			{
				heading: '11. Заключительные положения',
				paragraphs: [
					'11.1. К Оферте применяется законодательство Российской Федерации.',
					'11.2. Споры разрешаются в претензионном порядке. Претензия направляется на paylist.info@gmail.com. Срок рассмотрения — 30 календарных дней.',
					'11.3. Пользователь может получить разъяснения по Оферте, направив обращение на paylist.info@gmail.com.',
					'11.4. Оферта действует бессрочно до замены её новой версией либо до отключения Автоплатежей Пользователем.',
					'11.5. Актуальная версия Оферты в свободном доступе расположена в сети Интернет по адресу https://paylist.site/legal/offer.',
				],
			},
		],
	},
	en: {
		title: 'Public Autopayment Offer',
		updated: 'Last updated: August 16, 2026',
		intro: 'This public offer (the “Offer”) sets the terms of automatic recurring charges (“Autopayments”) for the Paylist Premium plan provided by Artem Vladimirovich Vlasov (the “Operator”) via the website https://paylist.site/',
		callout:
			'By checking consent on the plans page and making the first payment, the User accepts the Offer. Autopayments can be disabled in Account — the Cancel subscription button. Current version: https://paylist.site/legal/offer. Contact: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. General provisions',
				paragraphs: [
					'1.1. The Offer is a public offer within the meaning of Articles 435 and 437 of the Civil Code of the Russian Federation and is addressed to an indefinite number of persons.',
					'1.2. The Offer applies together with the User Agreement https://paylist.site/legal/terms and the Personal Data Processing Policy https://paylist.site/legal/privacy.',
					'1.3. The Operator may amend the Offer. A new version takes effect when posted at https://paylist.site/legal/offer, unless another effective date is stated in that version.',
					'1.4. Payments are processed by YooKassa (NCO YooMoney LLC). The Operator does not store the User’s full bank card number.',
				],
			},
			{
				heading: '2. Key terms used in the Offer',
				paragraphs: [
					'2.1. Operator — Artem Vladimirovich Vlasov, the person providing the Paylist service.',
					'2.2. User — a natural person who has accepted the Offer.',
					'2.3. Website — the Operator’s website on the internet at https://paylist.site/',
					'2.4. Service (Paylist) — the Operator’s web application for tracking subscriptions, spending analytics and payment reminders.',
					'2.5. Premium plan — the paid plan of the Service whose price, period and features are shown on the plans page before payment.',
					'2.6. Autopayments — automatic recurring charges for the Premium plan from the User’s saved payment method without additional confirmation of each charge.',
					'2.7. Acceptance — full and unconditional acceptance of the Offer by ticking the consent box on the plans page and making the first payment for the Premium plan.',
					'2.8. Payment provider — YooKassa (NCO YooMoney LLC), which accepts and processes payments.',
				],
			},
			{
				heading: '3. Subject of the Offer',
				paragraphs: [
					'3.1. The Operator undertakes to provide the User with access to the Premium plan, and the User undertakes to pay for the Premium plan in the manner set out in the Offer.',
					'3.2. The Offer governs only Autopayments for the Premium plan. Use of the Service otherwise is governed by the User Agreement.',
				],
			},
			{
				heading: '4. Acceptance of the Offer',
				paragraphs: [
					'4.1. Acceptance of the Offer is the combination of the following actions of the User: reviewing the Offer; ticking the consent box on the plans page; making the first payment for the Premium plan.',
					'4.2. From the moment of Acceptance, a contract on the terms of the Offer is deemed concluded between the Operator and the User.',
					'4.3. By making Acceptance, the User confirms that the amount, frequency and the procedure for disabling Autopayments are understood.',
				],
			},
			{
				heading: '5. Enabling Autopayments',
				paragraphs: [
					'5.1. When taking out the Premium plan, the first payment is made through the Payment provider.',
					'5.2. The User agrees to the payment method being saved with the Payment provider for subsequent charges.',
					'5.3. After a successful first payment, Autopayments are enabled in the User’s Account.',
				],
			},
			{
				heading: '6. Amount and frequency of charges',
				paragraphs: [
					'6.1. The charge amount and billing period are shown on the plans page before payment and on the Premium card in Account.',
					'6.2. The Premium plan is charged once per plan period (as a rule, every 30 days) in the amount stated at checkout.',
					'6.3. The next charge date is the end date of the paid Premium period shown in the Account. The charge is made at 21:00 Moscow time on the billing day.',
					'6.4. If Autopayments are not disabled, the same amount is charged on the charge date and access to the Premium plan is extended for the next period.',
				],
			},
			{
				heading: '7. Disabling Autopayments',
				paragraphs: [
					'7.1. The User may disable Autopayments at any time in Account: on the Premium card, tap Cancel subscription.',
					'7.2. After disabling, no further automatic charges are made.',
					'7.3. The Premium plan remains active until the already paid period end date and is not renewed thereafter unless the User takes it out again.',
					'7.4. Autopayments may be enabled again by taking out the Premium plan anew and accepting the Offer.',
				],
			},
			{
				heading: '8. Failed charges',
				paragraphs: [
					'8.1. If a scheduled charge fails, the Operator may retry the charge or suspend Autopayments.',
					'8.2. If the paid period ends without a successful payment, the Account is switched to the Free plan.',
					'8.3. The Operator is not liable for a refusal by a bank or the Payment provider to process a payment for reasons beyond the Operator’s control (insufficient funds, bank restrictions, expired card, etc.).',
				],
			},
			{
				heading: '9. Rights and duties of the Parties',
				paragraphs: [
					'9.1. The User may obtain access to the Premium plan upon successful payment, disable Autopayments in the manner of section 7 and contact the Operator about payment matters.',
					'9.2. The User must maintain a sufficient balance and an up-to-date saved payment method, and must independently monitor the amount and date of the next charge in the Account.',
					'9.3. The Operator may engage the Payment provider to accept payments and may suspend Autopayments if the Offer is breached.',
					'9.4. The Operator must display the charge amount and period before payment and provide a way to disable Autopayments in the Account.',
				],
			},
			{
				heading: '10. Liability',
				paragraphs: [
					'10.1. The Operator does not store the full card number and is not liable for the Payment provider’s processing of payment data.',
					'10.2. All information collected by the Payment provider is processed by it in accordance with its rules and privacy policy. The User must independently review those documents.',
					'10.3. To the extent permitted by the laws of the Russian Federation, the Operator is not liable for indirect losses related to Autopayments.',
				],
			},
			{
				heading: '11. Final provisions',
				paragraphs: [
					'11.1. The laws of the Russian Federation apply to the Offer.',
					'11.2. Disputes are resolved under a claims procedure. A claim is sent to paylist.info@gmail.com. The review period is 30 calendar days.',
					'11.3. The User may obtain explanations of the Offer by sending a request to paylist.info@gmail.com.',
					'11.4. The Offer remains in force indefinitely until replaced by a new version or until the User disables Autopayments.',
					'11.5. The current version of the Offer is freely available on the internet at https://paylist.site/legal/offer.',
				],
			},
		],
	},
	de: {
		title: 'Öffentliches Angebot für Autozahlungen',
		updated: 'Zuletzt aktualisiert: 16. August 2026',
		intro: 'Dieses öffentliche Angebot (im Folgenden das „Angebot“) bestimmt die Bedingungen automatischer wiederkehrender Abbuchungen („Autozahlungen“) für den Tarif Paylist Premium, bereitgestellt von Artem Wladimirowitsch Wlassow (im Folgenden der „Betreiber“) über die Website https://paylist.site/',
		callout:
			'Mit dem Setzen des Häkchens auf der Tarifseite und der ersten Zahlung nimmt der Nutzer das Angebot an. Autozahlungen können im Bereich „Konto“ deaktiviert werden — Schaltfläche „Abo kündigen“. Aktuelle Fassung: https://paylist.site/legal/offer. Kontakt: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Allgemeine Bestimmungen',
				paragraphs: [
					'1.1. Das Angebot ist ein öffentliches Angebot im Sinne der Artikel 435 und 437 des Zivilgesetzbuchs der Russischen Föderation und richtet sich an einen unbestimmten Personenkreis.',
					'1.2. Das Angebot gilt gemeinsam mit der Nutzervereinbarung https://paylist.site/legal/terms und der Richtlinie zur Verarbeitung personenbezogener Daten https://paylist.site/legal/privacy.',
					'1.3. Der Betreiber darf das Angebot ändern. Eine neue Fassung tritt mit der Veröffentlichung unter https://paylist.site/legal/offer in Kraft, sofern in der Fassung kein anderer Termin angegeben ist.',
					'1.4. Zahlungen verarbeitet YooKassa (ООО НКО «ЮМани»). Der Betreiber speichert nicht die vollständige Bankkartennummer des Nutzers.',
				],
			},
			{
				heading: '2. Wesentliche Begriffe, die im Angebot verwendet werden',
				paragraphs: [
					'2.1. Betreiber — Artem Wladimirowitsch Wlassow, die Person, die den Dienst Paylist bereitstellt.',
					'2.2. Nutzer — eine natürliche Person, die das Angebot angenommen hat.',
					'2.3. Website — die Website des Betreibers im Internet unter https://paylist.site/',
					'2.4. Dienst (Paylist) — die Webanwendung des Betreibers zur Erfassung von Abonnements, zur Ausgabenanalyse und für Zahlungserinnerungen.',
					'2.5. Tarif Premium — der entgeltliche Tarif des Dienstes, dessen Preis, Zeitraum und Funktionsumfang vor der Zahlung auf der Tarifseite angezeigt werden.',
					'2.6. Autozahlungen — automatische wiederkehrende Abbuchungen für den Tarif Premium von der gespeicherten Zahlungsmethode des Nutzers ohne zusätzliche Bestätigung jeder Abbuchung.',
					'2.7. Annahme — die vollständige und vorbehaltlose Annahme des Angebots durch Setzen des Häkchens auf der Tarifseite und Leistung der ersten Zahlung für den Tarif Premium.',
					'2.8. Zahlungsanbieter — YooKassa (ООО НКО «ЮМани»), der Zahlungen entgegennimmt und verarbeitet.',
				],
			},
			{
				heading: '3. Gegenstand des Angebots',
				paragraphs: [
					'3.1. Der Betreiber verpflichtet sich, dem Nutzer Zugang zum Tarif Premium zu gewähren, und der Nutzer verpflichtet sich, den Tarif Premium in der im Angebot festgelegten Weise zu bezahlen.',
					'3.2. Das Angebot regelt ausschließlich Autozahlungen für den Tarif Premium. Die übrige Nutzung des Dienstes wird durch die Nutzervereinbarung geregelt.',
				],
			},
			{
				heading: '4. Annahme des Angebots',
				paragraphs: [
					'4.1. Als Annahme des Angebots gilt die Gesamtheit der folgenden Handlungen des Nutzers: Kenntnisnahme des Angebots; Setzen des Häkchens auf der Tarifseite; Leistung der ersten Zahlung für den Tarif Premium.',
					'4.2. Ab dem Zeitpunkt der Annahme gilt zwischen dem Betreiber und dem Nutzer ein Vertrag zu den Bedingungen des Angebots als geschlossen.',
					'4.3. Mit der Annahme bestätigt der Nutzer, dass ihm Betrag, Rhythmus und das Verfahren zur Deaktivierung der Autozahlungen verständlich sind.',
				],
			},
			{
				heading: '5. Aktivierung der Autozahlungen',
				paragraphs: [
					'5.1. Beim Abschluss des Tarifs Premium erfolgt die erste Zahlung über den Zahlungsanbieter.',
					'5.2. Der Nutzer stimmt der Speicherung der Zahlungsmethode beim Zahlungsanbieter für Folgeabbuchungen zu.',
					'5.3. Nach erfolgreicher erster Zahlung werden Autozahlungen im Konto des Nutzers aktiviert.',
				],
			},
			{
				heading: '6. Betrag und Rhythmus der Abbuchungen',
				paragraphs: [
					'6.1. Abbuchungsbetrag und Zahlungszeitraum werden vor der Zahlung auf der Tarifseite und auf der Premium-Karte im Bereich „Konto“ angezeigt.',
					'6.2. Der Tarif Premium wird einmal je Tarifzeitraum (in der Regel alle 30 Tage) in der bei Abschluss angegebenen Höhe abgebucht.',
					'6.3. Das Datum der nächsten Abbuchung ist das im Konto angezeigte Endedatum des bezahlten Premium-Zeitraums. Die Abbuchung erfolgt um 21:00 Uhr Moskauer Zeit am Abbuchungstag.',
					'6.4. Sind Autozahlungen nicht deaktiviert, wird am Abbuchungstag derselbe Betrag abgebucht und der Zugang zum Tarif Premium um den nächsten Zeitraum verlängert.',
				],
			},
			{
				heading: '7. Deaktivierung der Autozahlungen',
				paragraphs: [
					'7.1. Der Nutzer darf Autozahlungen jederzeit im Bereich „Konto“ deaktivieren: auf der Premium-Karte ist „Abo kündigen“ zu tippen.',
					'7.2. Nach der Deaktivierung erfolgen keine weiteren automatischen Abbuchungen.',
					'7.3. Der Tarif Premium bleibt bis zum bereits bezahlten Endedatum des Zeitraums aktiv und wird danach nicht verlängert, sofern der Nutzer ihn nicht erneut abschließt.',
					'7.4. Eine erneute Aktivierung der Autozahlungen ist durch einen neuen Abschluss des Tarifs Premium und die Annahme des Angebots möglich.',
				],
			},
			{
				heading: '8. Fehlgeschlagene Abbuchungen',
				paragraphs: [
					'8.1. Ist eine fällige Abbuchung fehlgeschlagen, darf der Betreiber die Abbuchung wiederholen oder Autozahlungen aussetzen.',
					'8.2. Endet der bezahlte Zeitraum ohne erfolgreiche Zahlung, wird das Konto auf den Tarif Free umgestellt.',
					'8.3. Der Betreiber haftet nicht für die Ablehnung einer Zahlung durch die Bank oder den Zahlungsanbieter aus Gründen, die außerhalb des Einflusses des Betreibers liegen (unzureichende Deckung, Beschränkungen der Bank, Ablauf der Karte usw.).',
				],
			},
			{
				heading: '9. Rechte und Pflichten der Parteien',
				paragraphs: [
					'9.1. Der Nutzer darf bei erfolgreicher Zahlung Zugang zum Tarif Premium erhalten, Autozahlungen nach Abschnitt 7 deaktivieren und sich in Zahlungsfragen an den Betreiber wenden.',
					'9.2. Der Nutzer muss für ausreichende Deckung und eine aktuelle gespeicherte Zahlungsmethode sorgen und Betrag und Datum der nächsten Abbuchung selbst im Konto überwachen.',
					'9.3. Der Betreiber darf den Zahlungsanbieter zur Entgegennahme von Zahlungen einsetzen und Autozahlungen bei einem Verstoß gegen das Angebot aussetzen.',
					'9.4. Der Betreiber muss Betrag und Zeitraum der Abbuchung vor der Zahlung anzeigen und die Möglichkeit zur Deaktivierung der Autozahlungen im Konto bereitstellen.',
				],
			},
			{
				heading: '10. Haftung',
				paragraphs: [
					'10.1. Der Betreiber speichert nicht die vollständige Kartennummer und haftet nicht für die Verarbeitung von Zahlungsdaten durch den Zahlungsanbieter.',
					'10.2. Alle vom Zahlungsanbieter erhobenen Informationen werden von diesem gemäß seinen Regeln und seiner Datenschutzrichtlinie verarbeitet. Der Nutzer muss diese Dokumente selbst zur Kenntnis nehmen.',
					'10.3. Soweit nach dem Recht der Russischen Föderation zulässig, haftet der Betreiber nicht für indirekte Verluste im Zusammenhang mit Autozahlungen.',
				],
			},
			{
				heading: '11. Schlussbestimmungen',
				paragraphs: [
					'11.1. Auf das Angebot findet das Recht der Russischen Föderation Anwendung.',
					'11.2. Streitigkeiten werden im Reklamationsverfahren beigelegt. Die Reklamation ist an paylist.info@gmail.com zu richten. Die Prüfungsfrist beträgt 30 Kalendertage.',
					'11.3. Der Nutzer kann Erläuterungen zum Angebot erhalten, indem er eine Anfrage an paylist.info@gmail.com sendet.',
					'11.4. Das Angebot gilt unbefristet, bis es durch eine neue Fassung ersetzt oder bis Autozahlungen vom Nutzer deaktiviert werden.',
					'11.5. Die aktuelle Fassung des Angebots ist im Internet frei zugänglich unter https://paylist.site/legal/offer.',
				],
			},
		],
	},
	es: {
		title: 'Oferta pública de pagos automáticos',
		updated: 'Última actualización: 16 de agosto de 2026',
		intro: 'La presente oferta pública (en adelante, la «Oferta») determina las condiciones de los cobros periódicos automáticos («Pagos automáticos») del plan Paylist Premium, prestado por Artem Vladimirovich Vlasov (en adelante, el «Operador») a través del sitio web https://paylist.site/',
		callout:
			'Al marcar el consentimiento en la página de planes y realizar el primer pago, el Usuario acepta la Oferta. Los pagos automáticos pueden desactivarse en «Cuenta» — botón «Cancelar suscripción». Versión vigente: https://paylist.site/legal/offer. Contacto: paylist.info@gmail.com.',
		sections: [
			{
				heading: '1. Disposiciones generales',
				paragraphs: [
					'1.1. La Oferta es una oferta pública en el sentido de los artículos 435 y 437 del Código Civil de la Federación de Rusia y se dirige a un círculo indeterminado de personas.',
					'1.2. La Oferta se aplica juntamente con el Acuerdo de usuario https://paylist.site/legal/terms y la Política de tratamiento de datos personales https://paylist.site/legal/privacy.',
					'1.3. El Operador puede modificar la Oferta. La nueva redacción entra en vigor desde el momento de su publicación en https://paylist.site/legal/offer, salvo que en la propia redacción se indique otro plazo.',
					'1.4. Los pagos los procesa YooKassa (ООО НКО «ЮМани»). El Operador no almacena el número completo de la tarjeta bancaria del Usuario.',
				],
			},
			{
				heading: '2. Conceptos principales utilizados en la Oferta',
				paragraphs: [
					'2.1. Operador — Artem Vladimirovich Vlasov, la persona que presta el servicio Paylist.',
					'2.2. Usuario — la persona física que ha aceptado la Oferta.',
					'2.3. Sitio web — el sitio del Operador en internet en la dirección https://paylist.site/',
					'2.4. Servicio (Paylist) — la aplicación web del Operador para el registro de suscripciones, la analítica de gastos y los recordatorios de pago.',
					'2.5. Plan Premium — el plan oneroso del Servicio cuyo precio, periodo y conjunto de funciones se muestran en la página de planes antes del pago.',
					'2.6. Pagos automáticos — cobros periódicos automáticos del plan Premium desde el método de pago guardado del Usuario sin confirmación adicional de cada cobro.',
					'2.7. Aceptación — la aceptación plena e incondicional de la Oferta mediante la marca de consentimiento en la página de planes y el primer pago del plan Premium.',
					'2.8. Proveedor de pagos — YooKassa (ООО НКО «ЮМани»), que acepta y procesa los pagos.',
				],
			},
			{
				heading: '3. Objeto de la Oferta',
				paragraphs: [
					'3.1. El Operador se obliga a facilitar al Usuario el acceso al plan Premium, y el Usuario se obliga a pagar el plan Premium en la forma establecida en la Oferta.',
					'3.2. La Oferta regula exclusivamente los Pagos automáticos del plan Premium. El uso del Servicio en lo demás se rige por el Acuerdo de usuario.',
				],
			},
			{
				heading: '4. Aceptación de la Oferta',
				paragraphs: [
					'4.1. Se reconoce como aceptación de la Oferta el conjunto de las siguientes acciones del Usuario: toma de conocimiento de la Oferta; marca de consentimiento en la página de planes; realización del primer pago del plan Premium.',
					'4.2. Desde el momento de la Aceptación se considera celebrado entre el Operador y el Usuario un contrato en las condiciones de la Oferta.',
					'4.3. Al realizar la Aceptación, el Usuario confirma que comprende el importe, la periodicidad y el procedimiento de desactivación de los Pagos automáticos.',
				],
			},
			{
				heading: '5. Activación de los Pagos automáticos',
				paragraphs: [
					'5.1. Al contratar el plan Premium, el primer pago se realiza a través del Proveedor de pagos.',
					'5.2. El Usuario acepta que el método de pago se conserve en el Proveedor de pagos para cobros posteriores.',
					'5.3. Tras el primer pago correcto, los Pagos automáticos se activan en la Cuenta del Usuario.',
				],
			},
			{
				heading: '6. Importe y periodicidad de los cobros',
				paragraphs: [
					'6.1. El importe del cobro y el periodo de pago se muestran en la página de planes antes de pagar y en la tarjeta Premium de la sección «Cuenta».',
					'6.2. El plan Premium se cobra una vez por periodo del plan (por regla general, cada 30 días) por el importe indicado al contratar.',
					'6.3. La fecha del siguiente cobro es la fecha de fin del periodo Premium ya pagado que se muestra en la Cuenta. El cobro se realiza a las 21:00 hora de Moscú el día de cobro.',
					'6.4. Si los Pagos automáticos no están desactivados, en la fecha de cobro se cobra el mismo importe y el acceso al plan Premium se prorroga por el siguiente periodo.',
				],
			},
			{
				heading: '7. Desactivación de los Pagos automáticos',
				paragraphs: [
					'7.1. El Usuario puede desactivar los Pagos automáticos en cualquier momento en la sección «Cuenta»: en la tarjeta Premium debe pulsar «Cancelar suscripción».',
					'7.2. Tras la desactivación no se realizan nuevos cobros automáticos.',
					'7.3. El plan Premium permanece activo hasta la fecha de fin del periodo ya pagada y no se prorroga después, salvo que el Usuario lo contrate de nuevo.',
					'7.4. La reactivación de los Pagos automáticos es posible mediante una nueva contratación del plan Premium y la Aceptación de la Oferta.',
				],
			},
			{
				heading: '8. Cobros fallidos',
				paragraphs: [
					'8.1. Si el cobro correspondiente no se realiza, el Operador puede reiterar el intento de cobro o suspender los Pagos automáticos.',
					'8.2. Si el periodo pagado termina sin un pago correcto, la Cuenta pasa al plan Free.',
					'8.3. El Operador no responde de la denegación del banco o del Proveedor de pagos de tramitar el pago por causas ajenas al Operador (fondos insuficientes, restricciones del banco, caducidad de la tarjeta, etc.).',
				],
			},
			{
				heading: '9. Derechos y obligaciones de las Partes',
				paragraphs: [
					'9.1. El Usuario puede obtener acceso al plan Premium si el pago es correcto, desactivar los Pagos automáticos conforme a la sección 7 y dirigirse al Operador sobre cuestiones de pago.',
					'9.2. El Usuario debe garantizar saldo suficiente y la vigencia del método de pago guardado, y seguir por sí mismo el importe y la fecha del siguiente cobro en la Cuenta.',
					'9.3. El Operador puede recurrir al Proveedor de pagos para aceptar pagos y suspender los Pagos automáticos en caso de incumplimiento de la Oferta.',
					'9.4. El Operador debe mostrar el importe y el periodo de cobro antes del pago y facilitar la posibilidad de desactivar los Pagos automáticos en la Cuenta.',
				],
			},
			{
				heading: '10. Responsabilidad',
				paragraphs: [
					'10.1. El Operador no almacena el número completo de la tarjeta y no responde del tratamiento de los datos de pago por el Proveedor de pagos.',
					'10.2. Toda la información recogida por el Proveedor de pagos es tratada por este conforme a sus normas y política de privacidad. El Usuario debe tomar conocimiento de dichos documentos por sí mismo.',
					'10.3. En la medida permitida por la legislación de la Federación de Rusia, el Operador no responde de perjuicios indirectos relacionados con los Pagos automáticos.',
				],
			},
			{
				heading: '11. Disposiciones finales',
				paragraphs: [
					'11.1. A la Oferta se aplica la legislación de la Federación de Rusia.',
					'11.2. Las controversias se resuelven en vía de reclamación. La reclamación se envía a paylist.info@gmail.com. El plazo de examen es de 30 días naturales.',
					'11.3. El Usuario puede obtener aclaraciones sobre la Oferta enviando una consulta a paylist.info@gmail.com.',
					'11.4. La Oferta rige de forma indefinida hasta que se sustituya por una nueva versión o hasta que el Usuario desactive los Pagos automáticos.',
					'11.5. La versión vigente de la Oferta está disponible libremente en internet en https://paylist.site/legal/offer.',
				],
			},
		],
	},
};
