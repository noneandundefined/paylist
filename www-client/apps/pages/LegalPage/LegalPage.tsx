import { useEffect, useState } from 'react';
import { Helmet } from 'react-helmet-async';
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getLegalDocument, isLegalDocumentType } from '@/content/legal';
import { APP_NAME, getAppLanguage } from '@/constants/Language.constant';
import { ROUTES } from '@/constants/constants';
import PageHeader from '@/components/common/PageHeader/PageHeader';
import CookiePreferencesLink from '@/components/common/CookieConsent/CookiePreferencesLink';

const TOKEN_PATTERN = /(https?:\/\/[^\s]+)|([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})/g;

const LinkedText = ({ text }: { text: string }) => {
	const parts = text.split(TOKEN_PATTERN).filter((part): part is string => Boolean(part));

	return (
		<>
			{parts.map((part, index) => {
				if (part.startsWith('http://') || part.startsWith('https://')) {
					return (
						<a key={`${part}-${index}`} href={part} target="_blank" rel="noopener noreferrer" className="underline">
							{part}
						</a>
					);
				}

				if (part.includes('@') && !part.includes(' ')) {
					return (
						<a key={`${part}-${index}`} href={`mailto:${part}`} className="underline">
							{part}
						</a>
					);
				}

				return <span key={`${part}-${index}`}>{part}</span>;
			})}
		</>
	);
};

const LegalPage = () => {
	const { type } = useParams<{ type: string }>();
	const { t, i18n } = useTranslation();
	const navigate = useNavigate();
	const language = getAppLanguage(i18n.language);
	const [activeId, setActiveId] = useState('section-1');

	const legalDocument = isLegalDocumentType(type) ? getLegalDocument(type, language) : null;

	useEffect(() => {
		if (!legalDocument) {
			return;
		}

		const headings = legalDocument.sections.map((_, index) => window.document.getElementById(`section-${index + 1}`)).filter((element): element is HTMLElement => element !== null);

		if (headings.length === 0) {
			return;
		}

		const observer = new IntersectionObserver(
			(entries) => {
				const visible = entries.filter((entry) => entry.isIntersecting).sort((left, right) => right.intersectionRatio - left.intersectionRatio)[0];

				if (visible?.target.id) {
					setActiveId(visible.target.id);
				}
			},
			{ rootMargin: '-15% 0px -65% 0px', threshold: [0, 0.25, 1] }
		);

		headings.forEach((heading) => observer.observe(heading));

		return () => observer.disconnect();
	}, [legalDocument]);

	if (!isLegalDocumentType(type) || !legalDocument) {
		return <Navigate to={ROUTES.NOT_FOUND} replace />;
	}

	const onClose = () => {
		if (window.history.length > 1) {
			navigate(-1);
			return;
		}

		navigate(ROUTES.SIGNIN);
	};

	const onSectionClick = (event: React.MouseEvent<HTMLAnchorElement>, sectionId: string) => {
		event.preventDefault();
		window.document.getElementById(sectionId)?.scrollIntoView({ behavior: 'smooth', block: 'start' });
		setActiveId(sectionId);
	};

	const toc = (
		<nav aria-label={t('legal.on-this-page')}>
			<p className="mb-3 text-[13px] font-semibold tracking-wide gu-text-primary">{t('legal.on-this-page')}</p>
			<ul className="space-y-2.5">
				{legalDocument.sections.map((section, index) => {
					const sectionId = `section-${index + 1}`;
					const isActive = activeId === sectionId;

					return (
						<li key={sectionId}>
							<a
								href={`#${sectionId}`}
								onClick={(event) => onSectionClick(event, sectionId)}
								className={`block text-[13px] leading-snug no-underline hover:no-underline ${isActive ? 'font-medium gu-text-primary' : 'gu-text-muted hover:text-[var(--text-primary)]'}`}
							>
								{section.heading}
							</a>
						</li>
					);
				})}
			</ul>
		</nav>
	);

	return (
		<div className="min-h-screen bg-[var(--surface)] gu-text-primary">
			<Helmet>
				<title>{legalDocument.title}</title>
			</Helmet>

			<div className="mx-auto w-full max-w-6xl px-4 pb-16 pt-5 lg:px-8">
				<div className="mb-6 flex items-center justify-between">
					<Link to={ROUTES.HOME} className="text-[15px] font-semibold no-underline hover:no-underline gu-text-primary">
						{APP_NAME}
					</Link>
					<PageHeader variant="close" onClose={onClose} backLabel={t('action.close')} />
				</div>

				<div className="lg:grid lg:grid-cols-[minmax(0,1fr)_220px] lg:items-start lg:gap-16">
					<article className="min-w-0 max-w-3xl">
						<h1 className="text-[34px] font-bold leading-tight tracking-tight gu-text-primary">{legalDocument.title}</h1>
						<p className="mt-3 text-[13px] gu-text-muted">{legalDocument.updated}</p>
						<p className="mt-6 text-[16px] leading-7 gu-text-secondary">
							<LinkedText text={legalDocument.intro} />
						</p>

						<div className="mt-6 rounded-xl border border-[#e8d48b] bg-[#fff8e7] px-5 py-4 text-[15px] leading-6 text-[#5c4d12] dark:border-amber-500/25 dark:bg-amber-500/10 dark:text-amber-100">
							<LinkedText text={legalDocument.callout} />
						</div>

						<div className="mt-8 lg:hidden">{toc}</div>

						{legalDocument.sections.map((section, index) => {
							const sectionId = `section-${index + 1}`;

							return (
								<section key={sectionId} className="mt-10">
									<h2 id={sectionId} className="scroll-mt-8 text-[22px] font-semibold leading-snug gu-text-primary">
										{section.heading}
									</h2>
									<div className="mt-3 space-y-4 text-[15px] leading-7 gu-text-secondary">
										{section.paragraphs.map((paragraph) => (
											<p key={paragraph}>
												<LinkedText text={paragraph} />
											</p>
										))}
									</div>
								</section>
							);
						})}
					</article>

					<aside className="sticky top-8 hidden lg:block">{toc}</aside>
				</div>

				<CookiePreferencesLink className="mt-10 block w-full text-center" />
			</div>
		</div>
	);
};

export default LegalPage;
