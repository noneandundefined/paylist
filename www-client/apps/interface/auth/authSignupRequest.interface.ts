export interface AuthSignupRequest {
	first_name: string | null;
	last_name: string | null;
	email: string;
	password: string;
	referral_code?: string;
}
