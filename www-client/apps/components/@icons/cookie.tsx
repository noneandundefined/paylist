interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: () => void;
}

export const Cookie: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox="0 0 24 24" fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M12,3A9,9 0 0,0 3,12A9,9 0 0,0 12,21A9,9 0 0,0 21,12A9,9 0 0,0 12,3M9.5,8A1.5,1.5 0 0,1 11,9.5A1.5,1.5 0 0,1 9.5,11A1.5,1.5 0 0,1 8,9.5A1.5,1.5 0 0,1 9.5,8M16.5,11A1.5,1.5 0 0,1 18,12.5A1.5,1.5 0 0,1 16.5,14A1.5,1.5 0 0,1 15,12.5A1.5,1.5 0 0,1 16.5,11M8,14A1,1 0 0,1 9,15A1,1 0 0,1 8,16A1,1 0 0,1 7,15A1,1 0 0,1 8,14M12.5,16A1.5,1.5 0 0,1 14,17.5A1.5,1.5 0 0,1 12.5,19A1.5,1.5 0 0,1 11,17.5A1.5,1.5 0 0,1 12.5,16Z" />
		</svg>
	);
};

export default Cookie;
