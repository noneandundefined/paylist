import SubscriptionIcon from '@/components/common/TrackedSubscription/SubscriptionIcon';
import { PhoneStroke } from './StartLandingIcons';

const LIME = '#d4ef4f';

const StartPhoneShowcase = () => {
	return (
		<div className="relative mx-auto w-[280px] sm:w-[300px]">
			<div className="pointer-events-none absolute -inset-10 rounded-full bg-[#d4ef4f]/12 blur-3xl" />

			<div className="pointer-events-none relative overflow-hidden rounded-[42px] border-[8px] border-[#1a1a1a] bg-[#0a0a0a] shadow-[0_40px_80px_rgba(0,0,0,0.55)]">
				<div className="absolute left-1/2 top-2.5 z-20 h-[22px] w-[96px] -translate-x-1/2 rounded-full bg-black" />

				<div className="flex items-center justify-between px-5 pb-1 pt-3 text-[11px] font-semibold text-white/80">
					<span>9:41</span>
					<span className="flex items-center gap-1 text-white">
						<span className="inline-block h-2 w-3.5 rounded-[2px] border border-white/80">
							<span className="ml-[1px] mt-[1px] block h-1.5 w-2 rounded-[1px] bg-white/80" />
						</span>
					</span>
				</div>

				<div className="px-3.5 pb-5 pt-2">
					<div className="mb-4 flex items-center justify-between gap-2">
						<div className="flex min-w-0 items-center gap-2.5">
							<div className="flex h-9 w-9 items-center justify-center rounded-full bg-[#2a2a2a] text-[12px] font-bold text-white">A</div>
							<div className="min-w-0">
								<p className="text-[11px] text-white/45">Добрый вечер</p>
								<p className="truncate text-[15px] font-bold leading-tight text-white">artemiik</p>
							</div>
						</div>
						<span className="rounded-full bg-[#1d4ed8] px-2.5 py-1 text-[10px] font-semibold uppercase tracking-wide text-white">Premium</span>
					</div>

					<section className="mb-3 rounded-[20px] bg-[#141414] p-4">
						<p className="text-[11px] font-medium text-white/45">Всего расходов</p>
						<p className="mt-1 text-[40px] font-bold leading-none tracking-tight text-white">
							$85<span className="text-[22px] font-semibold text-white/45">.21</span>
						</p>
						<div className="mt-3 h-1.5 overflow-hidden rounded-full bg-white/10">
							<div className="h-full w-[62%] rounded-full bg-[#d4ef4f]" />
						</div>
					</section>

					<p className="mb-2 px-0.5 text-[13px] font-semibold text-white">Активные подписки</p>

					<div className="space-y-2">
						<div className="flex items-center gap-2.5 rounded-[16px] bg-[#141414] px-3 py-2.5">
							<div className="flex h-9 w-9 items-center justify-center rounded-xl bg-[#1c1c1c]">
								<PhoneStroke size={18} color={LIME} />
							</div>
							<div className="min-w-0 flex-1">
								<p className="truncate text-[13px] font-semibold text-white">Мобильная связь</p>
								<p className="text-[11px] text-white/40">завтра</p>
							</div>
							<p className="text-[13px] font-bold text-white">$12.00</p>
						</div>

						<div className="flex items-center gap-2.5 rounded-[16px] bg-[#141414] px-3 py-2.5">
							<SubscriptionIcon name="Tinkoff" size="xs" />
							<div className="min-w-0 flex-1">
								<p className="truncate text-[13px] font-semibold text-white">Tinkoff Bank Pro</p>
								<p className="text-[11px] text-white/40">через 5 дн.</p>
							</div>
							<p className="text-[13px] font-bold text-white">$9.99</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default StartPhoneShowcase;
