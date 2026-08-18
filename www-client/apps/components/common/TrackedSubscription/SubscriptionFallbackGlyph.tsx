import Car from '@/components/@icons/car';
import Home from '@/components/@icons/home';
import TrainCar from '@/components/@icons/train-car';
import WalletOutline from '@/components/@icons/wallet';
import CloudOutline from '@/components/@icons/cloud-outline';
import Memory from '@/components/@icons/memory';
import BookOpenBlankVariantOutline from '@/components/@icons/book-open-blank-variant-outline';
import { getSubscriptionFallbackKind, type SubscriptionFallbackKind } from '@/utils/subscriptionFallbackIconUtils';

interface IconProps {
	fill: string;
	size: number;
}

type GlyphIcon = (props: IconProps) => React.ReactNode;

const PathIcon = ({ fill, size, d }: IconProps & { d: string }) => (
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width={size} height={size} fill={fill} aria-hidden>
		<path d={d} />
	</svg>
);

const Cellphone = ({ fill, size }: IconProps) => <PathIcon fill={fill} size={size} d="M17,19H7V5H17M17,1H7C5.89,1 5,1.89 5,3V21A2,2 0 0,0 7,23H17A2,2 0 0,0 19,21V3C19,1.89 18.1,1 17,1Z" />;

const Wifi = ({ fill, size }: IconProps) => (
	<PathIcon
		fill={fill}
		size={size}
		d="M12,21L15.6,16.2C14.6,15.45 13.35,15 12,15C10.65,15 9.4,15.45 8.4,16.2L12,21M12,3C7.95,3 4.21,4.34 1.2,6.6L3,9C5.5,7.12 8.62,6 12,6C15.38,6 18.5,7.12 21,9L22.8,6.6C19.79,4.34 16.05,3 12,3M12,9C9.3,9 6.81,9.89 4.8,11.4L6.6,13.8C8.1,12.67 9.97,12 12,12C14.03,12 15.9,12.67 17.4,13.8L19.2,11.4C17.19,9.89 14.7,9 12,9Z"
	/>
);

const Flash = ({ fill, size }: IconProps) => <PathIcon fill={fill} size={size} d="M7,2V13H10V22L17,10H13L17,2H7Z" />;

const MusicNote = ({ fill, size }: IconProps) => <PathIcon fill={fill} size={size} d="M12,3V12.26C11.5,12.09 11,12 10.5,12C8.57,12 7,13.57 7,15.5C7,17.43 8.57,19 10.5,19C12.43,19 14,17.43 14,15.5V6H18V3H12Z" />;

const TelevisionPlay = ({ fill, size }: IconProps) => <PathIcon fill={fill} size={size} d="M21,17H3V5H21M21,3H3A2,2 0 0,0 1,5V17A2,2 0 0,0 3,19H8V21H16V19H21A2,2 0 0,0 23,17V5A2,2 0 0,0 21,3M10,8V14L15,11L10,8Z" />;

const Controller = ({ fill, size }: IconProps) => (
	<PathIcon
		fill={fill}
		size={size}
		d="M21,6H3A2,2 0 0,0 1,8V15A2,2 0 0,0 3,17H6L8,20H10L8,17H16L14,20H16L18,17H21A2,2 0 0,0 23,15V8A2,2 0 0,0 21,6M11,13H9V15H7V13H5V11H7V9H9V11H11M18,13A1,1 0 0,1 17,12A1,1 0 0,1 18,11A1,1 0 0,1 19,12A1,1 0 0,1 18,13M20,10A1,1 0 0,1 19,9A1,1 0 0,1 20,8A1,1 0 0,1 21,9A1,1 0 0,1 20,10Z"
	/>
);

const Dumbbell = ({ fill, size }: IconProps) => (
	<PathIcon
		fill={fill}
		size={size}
		d="M20.57,14.86L22,13.43L20.57,12L17,15.57L8.43,7L12,3.43L10.57,2L9.14,3.43L7.71,2L5.57,4.14L4.14,2.71L2.71,4.14L4.14,5.57L2,7.71L3.43,9.14L2,10.57L3.43,12L7,8.43L15.57,17L12,20.57L13.43,22L14.86,20.57L16.29,22L18.43,19.86L19.86,21.29L21.29,19.86L19.86,18.43L22,16.29L20.57,14.86Z"
	/>
);

const Newspaper = ({ fill, size }: IconProps) => (
	<PathIcon
		fill={fill}
		size={size}
		d="M20,11H4V8H20M20,15H13V13H20M20,19H13V17H20M11,19H4V13H11M20.33,4.67L18.67,3L17,4.67L15.33,3L13.67,4.67L12,3L10.33,4.67L8.67,3L7,4.67L5.33,3L3.67,4.67L2,3V19A2,2 0 0,0 4,21H20A2,2 0 0,0 22,19V3L20.33,4.67Z"
	/>
);

const ICONS: Record<SubscriptionFallbackKind, GlyphIcon> = {
	phone: Cellphone,
	internet: Wifi,
	'transport-car': Car,
	'transport-train': TrainCar,
	home: Home,
	utilities: Flash,
	music: MusicNote,
	streaming: TelevisionPlay,
	gaming: Controller,
	cloud: CloudOutline,
	productivity: Memory,
	fitness: Dumbbell,
	news: Newspaper,
	education: BookOpenBlankVariantOutline,
	wallet: WalletOutline,
};

interface SubscriptionFallbackGlyphProps {
	name: string;
	categories?: string[];
	size: number;
	fill: string;
}

const SubscriptionFallbackGlyph = ({ name, categories, size, fill }: SubscriptionFallbackGlyphProps) => {
	const kind = getSubscriptionFallbackKind(name, categories);
	const Icon = ICONS[kind];

	return <Icon fill={fill} size={size} />;
};

export default SubscriptionFallbackGlyph;
