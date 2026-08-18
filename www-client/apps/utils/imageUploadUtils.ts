const MAX_AVATAR_BYTES = 900 * 1024;
const MAX_AVATAR_EDGE = 1024;

const replaceExtension = (name: string, ext: string) => {
	const trimmed = name.trim();
	const dot = trimmed.lastIndexOf('.');
	const base = dot > 0 ? trimmed.slice(0, dot) : trimmed || 'avatar';
	return `${base}.${ext}`;
};

export const compressImageForUpload = async (file: File): Promise<File> => {
	if (!file.type.startsWith('image/') || file.type === 'image/gif') {
		return file;
	}

	let bitmap: ImageBitmap;
	try {
		bitmap = await createImageBitmap(file);
	} catch {
		return file;
	}

	try {
		const scale = Math.min(1, MAX_AVATAR_EDGE / Math.max(bitmap.width, bitmap.height));
		const width = Math.max(1, Math.round(bitmap.width * scale));
		const height = Math.max(1, Math.round(bitmap.height * scale));

		if (scale === 1 && file.size <= MAX_AVATAR_BYTES) {
			return file;
		}

		const canvas = document.createElement('canvas');
		canvas.width = width;
		canvas.height = height;
		const context = canvas.getContext('2d');
		if (!context) {
			return file;
		}

		context.drawImage(bitmap, 0, 0, width, height);

		let blob: Blob | null = null;
		for (let quality = 0.86; quality >= 0.5; quality -= 0.12) {
			blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', quality));
			if (blob && blob.size <= MAX_AVATAR_BYTES) {
				break;
			}
		}

		if (!blob) {
			return file;
		}

		return new File([blob], replaceExtension(file.name, 'jpg'), { type: 'image/jpeg', lastModified: Date.now() });
	} finally {
		bitmap.close();
	}
};
