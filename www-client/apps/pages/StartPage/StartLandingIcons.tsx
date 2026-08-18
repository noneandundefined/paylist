interface StrokeIconProps {
	size?: number;
	color?: string;
	className?: string;
	strokeWidth?: number;
}

const Svg = ({ size = 24, className, children }: { size?: number; className?: string; children: React.ReactNode }) => (
	<svg xmlns="http://www.w3.org/2000/svg" width={size} height={size} viewBox="0 0 24 24" fill="none" className={className} aria-hidden>
		{children}
	</svg>
);

const stroke = (color: string, width: number) => ({
	stroke: color,
	strokeWidth: width,
	strokeLinecap: 'round' as const,
	strokeLinejoin: 'round' as const,
});

export const WalletStroke = ({ size = 28, color = '#d4ef4f', className, strokeWidth = 1.6 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<rect x="3" y="6.5" width="18" height="13" rx="2.4" {...stroke(color, strokeWidth)} />
		<path d="M3 10.2h18" {...stroke(color, strokeWidth)} />
		<circle cx="16.4" cy="14.6" r="1.15" fill={color} />
	</Svg>
);

export const BellStroke = ({ size = 28, color = '#d4ef4f', className, strokeWidth = 1.6 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<path d="M6.2 16.6V11a5.8 5.8 0 0 1 11.6 0v5.6" {...stroke(color, strokeWidth)} />
		<path d="M4.6 16.6h14.8" {...stroke(color, strokeWidth)} />
		<path d="M10 19.2a2 2 0 0 0 4 0" {...stroke(color, strokeWidth)} />
		<path d="M12 3.6V5.2" {...stroke(color, strokeWidth)} />
	</Svg>
);

export const PieStroke = ({ size = 28, color = '#d4ef4f', className, strokeWidth = 1.6 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<circle cx="12" cy="12" r="8.2" {...stroke(color, strokeWidth)} />
		<path d="M12 12V3.8" {...stroke(color, strokeWidth)} />
		<path d="M12 12l7.1 4.1" {...stroke(color, strokeWidth)} />
	</Svg>
);

export const ShieldStroke = ({ size = 28, color = '#d4ef4f', className, strokeWidth = 1.6 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<path d="M12 3.4 19.4 6v5.3c0 4.4-3 7.3-7.4 9.1C7.6 18.6 4.6 15.7 4.6 11.3V6L12 3.4Z" {...stroke(color, strokeWidth)} />
		<path d="m8.8 12.1 2.1 2.1 4.4-4.6" {...stroke(color, strokeWidth)} />
	</Svg>
);

export const SendStroke = ({ size = 20, color = '#111', className, strokeWidth = 1.8 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<path d="M4 11.2 20 4.6 13.4 20l-2.3-6.6L4 11.2Z" {...stroke(color, strokeWidth)} />
		<path d="m11.1 13.4 8.9-8.8" {...stroke(color, strokeWidth)} />
	</Svg>
);

export const SparkleStroke = ({ size = 18, color = '#d4ef4f', className }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<path d="M12 2.4 13.6 9 20.2 10.6 13.6 12.2 12 18.8 10.4 12.2 3.8 10.6 10.4 9 12 2.4Z" fill={color} />
	</Svg>
);

export const DotsStroke = ({ size = 28, color = '#d4ef4f', className }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<circle cx="6" cy="12" r="1.7" fill={color} />
		<circle cx="12" cy="12" r="1.7" fill={color} />
		<circle cx="18" cy="12" r="1.7" fill={color} />
	</Svg>
);

export const PhoneStroke = ({ size = 22, color = '#d4ef4f', className, strokeWidth = 1.6 }: StrokeIconProps) => (
	<Svg size={size} className={className}>
		<rect x="7.2" y="3.2" width="9.6" height="17.6" rx="2.2" {...stroke(color, strokeWidth)} />
		<path d="M10.4 17.8h3.2" {...stroke(color, strokeWidth)} />
	</Svg>
);
