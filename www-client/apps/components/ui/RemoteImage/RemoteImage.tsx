import { useEffect, useRef, useState } from 'react';
import ImageSpinner from '@/components/ui/ImageSpinner/ImageSpinner';

interface RemoteImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
	spinnerSize?: number;
	spinnerLight?: boolean;
}

const RemoteImage: React.FC<RemoteImageProps> = ({ className = '', src, spinnerSize = 16, spinnerLight = false, onLoad, onError, alt = '', ...props }) => {
	const imageRef = useRef<HTMLImageElement>(null);
	const [loaded, setLoaded] = useState(false);

	useEffect(() => {
		setLoaded(false);

		const image = imageRef.current;
		if (image?.complete && image.naturalWidth > 0) {
			setLoaded(true);
		}
	}, [src]);

	return (
		<>
			{!loaded && src ? (
				<span className="absolute inset-0 z-[1] flex items-center justify-center">
					<ImageSpinner size={spinnerSize} light={spinnerLight} />
				</span>
			) : null}
			<img
				{...props}
				ref={imageRef}
				alt={alt}
				src={src}
				draggable={false}
				onDragStart={(event) => event.preventDefault()}
				className={className}
				onLoad={(event) => {
					setLoaded(true);
					onLoad?.(event);
				}}
				onError={(event) => {
					setLoaded(true);
					onError?.(event);
				}}
			/>
		</>
	);
};

export default RemoteImage;
