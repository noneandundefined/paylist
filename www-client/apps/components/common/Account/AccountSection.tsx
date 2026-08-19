import { useState } from 'react';
import ChevronDown from '@/components/@icons/chevron-down';

interface AccountSectionProps {
	title: string;
	children: React.ReactNode;
	collapsible?: boolean;
	defaultOpen?: boolean;
}

const AccountSection: React.FC<AccountSectionProps> = ({ title, children, collapsible = false, defaultOpen = true }) => {
	const [open, setOpen] = useState(defaultOpen);

	return (
		<section className="space-y-2">
			{collapsible ? (
				<button type="button" className="flex w-full items-center justify-between gap-2 px-1" onClick={() => setOpen((prev) => !prev)} aria-expanded={open}>
					<h2 className="text-[13px] font-medium uppercase tracking-wide gu-text-muted">{title}</h2>
					<ChevronDown fill="currentColor" size={25} className={`gu-text-muted shrink-0 transition-transform ${open ? 'rotate-180' : ''}`} />
				</button>
			) : (
				<h2 className="px-1 text-[13px] font-medium uppercase tracking-wide gu-text-muted">{title}</h2>
			)}
			{(!collapsible || open) && children}
		</section>
	);
};

export default AccountSection;
