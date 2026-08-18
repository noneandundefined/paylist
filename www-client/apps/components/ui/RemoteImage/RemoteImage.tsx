import { useEffect, useState } from 'react';
import ImageSpinner from '@/components/ui/ImageSpinner/ImageSpinner';
import { getCachedObjectUrl, isImageFailed, isImageReady, markImageFailed, resolveCachedImage } from '@/utils/imageCacheUtils';

interface RemoteImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
	spinnerSize?: number;
	spinnerLight?: boolean;
}

const RemoteImage: React.FC<RemoteImageProps> = ({ className = '', src, spinnerSize = 16, spinnerLight = false, onLoad, onError, alt = '', ...props }) => {
	const [displaySrc, setDisplaySrc] = useState(() => getCachedObjectUrl(src) ?? src);
	const [loaded, setLoaded] = useState(() => isImageReady(src));
	const [failed, setFailed] = useState(() => isImageFailed(src));

	useEffect(() => {
		if (!src) {
			setDisplaySrc(undefined);
			setLoaded(false);
			setFailed(false);
			return;
		}

		if (isImageFailed(src)) {
			setFailed(true);
			setLoaded(true);
			return;
		}

		const cached = getCachedObjectUrl(src);
		if (cached) {
			setDisplaySrc(cached);
			setLoaded(true);
			setFailed(false);
			return;
		}

		let cancelled = false;
		setFailed(false);
		setLoaded(isImageReady(src));
		setDisplaySrc(src);

		void resolveCachedImage(src)
			.then((url) => {
				if (cancelled) {
					return;
				}

				setDisplaySrc(url);
				setLoaded(true);
			})
			.catch(() => {
				if (cancelled) {
					return;
				}

				setFailed(true);
				setLoaded(true);
				onError?.({} as React.SyntheticEvent<HTMLImageElement, Event>);
			});

		return () => {
			cancelled = true;
		};
	}, [src]);

	if (!src || failed) {
		return null;
	}

	return (
		<>
			{!loaded ? (
				<span className="absolute inset-0 z-[1] flex items-center justify-center">
					<ImageSpinner size={spinnerSize} light={spinnerLight} />
				</span>
			) : null}
			<img
				{...props}
				alt={alt}
				src={displaySrc}
				draggable={false}
				onDragStart={(event) => event.preventDefault()}
				className={className}
				onLoad={(event) => {
					setLoaded(true);
					onLoad?.(event);
				}}
				onError={(event) => {
					markImageFailed(src);
					setFailed(true);
					setLoaded(true);
					onError?.(event);
				}}
			/>
		</>
	);
};

export default RemoteImage;
