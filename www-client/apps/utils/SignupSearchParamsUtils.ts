export type SignupFormQuery = {
	email?: string;
	first_name?: string | null;
	last_name?: string | null;
};

export const readSignupFormFromSearchParams = (params: URLSearchParams): SignupFormQuery => ({
	email: params.get('email') ?? undefined,
	first_name: params.get('first_name') ?? undefined,
	last_name: params.get('last_name') ?? undefined,
});

export const writeSignupFormToSearchParams = (params: URLSearchParams, values: SignupFormQuery) => {
	const setOrDelete = (key: keyof SignupFormQuery, value?: string) => {
		const trimmed = value?.trim();
		if (trimmed) {
			params.set(key, trimmed);
		} else {
			params.delete(key);
		}
	};

	setOrDelete('email', values.email);
	setOrDelete('first_name', values.first_name ?? undefined);
	setOrDelete('last_name', values.last_name ?? undefined);
};

export const buildSentConfirmEmailPath = (email: string) => {
	const params = new URLSearchParams();
	params.set('email', email.trim());

	return params.toString();
};
