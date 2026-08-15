import { useCallback, useRef } from 'react';

interface ShareWedgeSliderProps {
	value: number;
	min: number;
	max: number;
	marker?: number;
	markerLabel?: string;
	minLabel?: string;
	maxLabel?: string;
	onChange: (value: number) => void;
	ariaLabel: string;
}

const clamp = (value: number, min: number, max: number) => Math.min(max, Math.max(min, value));

const toRatio = (value: number, min: number, max: number) => {
	if (max <= min) {
		return 0;
	}

	return clamp((value - min) / (max - min), 0, 1);
};

const ShareWedgeSlider: React.FC<ShareWedgeSliderProps> = ({ value, min, max, marker, markerLabel, minLabel, maxLabel, onChange, ariaLabel }) => {
	const trackRef = useRef<HTMLDivElement>(null);

	const ratio = toRatio(value, min, max);
	const markerRatio = marker == null ? null : toRatio(marker, min, max);

	const updateFromClientX = useCallback(
		(clientX: number) => {
			const track = trackRef.current;

			if (!track) {
				return;
			}

			const rect = track.getBoundingClientRect();
			const nextRatio = clamp((clientX - rect.left) / rect.width, 0, 1);
			onChange(Math.round((min + nextRatio * (max - min)) * 10) / 10);
		},
		[max, min, onChange]
	);

	const onPointerDown = (event: React.PointerEvent<HTMLDivElement>) => {
		event.currentTarget.setPointerCapture(event.pointerId);
		updateFromClientX(event.clientX);
	};

	const onPointerMove = (event: React.PointerEvent<HTMLDivElement>) => {
		if (!event.currentTarget.hasPointerCapture(event.pointerId)) {
			return;
		}

		updateFromClientX(event.clientX);
	};

	return (
		<div className="relative px-1 pt-2 pb-1">
			<div
				ref={trackRef}
				role="slider"
				tabIndex={0}
				aria-label={ariaLabel}
				aria-valuemin={min}
				aria-valuemax={max}
				aria-valuenow={value}
				onPointerDown={onPointerDown}
				onPointerMove={onPointerMove}
				onKeyDown={(event) => {
					const step = event.key === 'ArrowLeft' || event.key === 'ArrowDown' ? -0.5 : event.key === 'ArrowRight' || event.key === 'ArrowUp' ? 0.5 : 0;

					if (!step) {
						return;
					}

					event.preventDefault();
					onChange(clamp(Math.round((value + step) * 10) / 10, min, max));
				}}
				className="relative h-11 cursor-pointer touch-none select-none"
			>
				<div
					className="absolute inset-x-0 top-1/2 h-2 -translate-y-1/2 overflow-hidden rounded-full"
					style={{
						background: 'linear-gradient(90deg, #7dd3fc 0%, #67e8f9 16%, #fde047 48%, #fb923c 78%, #ef4444 100%)',
					}}
				>
					<div
						className="absolute inset-0"
						style={{
							backgroundImage: 'repeating-linear-gradient(-32deg, transparent, transparent 7px, rgba(255,255,255,0.22) 7px, rgba(255,255,255,0.22) 8px)',
						}}
					/>
				</div>

				{markerRatio != null ? <span className="absolute bottom-0 top-1 w-px bg-[var(--text-primary)]/35" style={{ left: `${markerRatio * 100}%` }} /> : null}

				<span
					className="pointer-events-none absolute top-1/2 h-7 w-7 -translate-x-1/2 -translate-y-1/2 rounded-full"
					style={{
						left: `${ratio * 100}%`,
						background: 'radial-gradient(circle at 32% 28%, #fff7ed 0%, #fdba74 42%, #f97316 100%)',
						boxShadow: '0 0 0 6px rgba(249, 115, 22, 0.22), 0 8px 18px rgba(249, 115, 22, 0.4)',
					}}
				/>
			</div>

			{(minLabel || maxLabel || markerLabel) && (
				<div className="relative mt-1 h-4 text-[11px] font-semibold tracking-[0.12em] gu-text-muted">
					{minLabel ? <span className="absolute left-0">{minLabel}</span> : null}
					{markerLabel && markerRatio != null ? (
						<span className="absolute -translate-x-1/2" style={{ left: `${markerRatio * 100}%` }}>
							{markerLabel}
						</span>
					) : null}
					{maxLabel ? <span className="absolute right-0">{maxLabel}</span> : null}
				</div>
			)}
		</div>
	);
};

export default ShareWedgeSlider;
