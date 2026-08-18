import { apiGet, apiPatch, apiPost } from '@/rest/apiClient';
import type { PaymentCheckoutRequest } from '@/interface/payment/paymentCheckoutRequest.interface';

const apiPath = '/payments';

export interface PaymentBillingResponse {
	auto_renew_enabled: boolean;
	has_payment_method: boolean;
	payment_method_type?: string | null;
	payment_method_title?: string | null;
	payment_method_saved_at?: string | null;
}

export interface PaymentHistoryResponse {
	id: number;
	created_at: string;
	plan_name: string;
	amount: number;
	currency: string;
	status: 'pending' | 'waiting_for_capture' | 'succeeded' | 'canceled' | 'failed';
	payment_kind: 'initial' | 'renewal' | 'manual';
	description?: string | null;
	paid_at?: string | null;
}

export interface AutoRenewUpdateRequest {
	enabled: boolean;
}

export interface PaymentCheckoutResponse {
	payment_id: string;
	confirmation_url: string;
}

export interface PaymentConfirmResponse {
	paid: boolean;
	status: string;
}

export const basicPaymentCheckout = async (payload: PaymentCheckoutRequest): Promise<PaymentCheckoutResponse> => apiPost(`${apiPath}/checkout`, payload);

export const basicPaymentConfirm = async (paymentId: string): Promise<PaymentConfirmResponse> => apiGet(`${apiPath}/confirm`, { params: { payment_id: paymentId }, skipErrorHandler: true });

export const basicPaymentBilling = async (): Promise<PaymentBillingResponse> => apiGet(`${apiPath}/billing`);

export const basicPaymentHistory = async (limit = 50): Promise<PaymentHistoryResponse[]> => apiGet(`${apiPath}/history`, { params: { limit } });

export const basicPaymentAutoRenewUpdate = async (payload: AutoRenewUpdateRequest): Promise<PaymentBillingResponse> => apiPatch(`${apiPath}/auto-renew`, payload);
