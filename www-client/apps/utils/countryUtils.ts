export const getCountryFlagUrl = (code: string): string | null => {
	const normalized = code.trim().toLowerCase();

	if (!/^[a-z]{2}$/.test(normalized)) {
		return null;
	}

	return `https://flagcdn.com/w40/${normalized}.png`;
};
