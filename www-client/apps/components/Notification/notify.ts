type NotificationType = 'success' | 'error';

export type NotificationItem = {
	id: string;
	type: NotificationType;
	message: string;
};

type Listener = (items: NotificationItem[]) => void;

class NotificationStore {
	private items: NotificationItem[] = [];
	private listeners: Listener[] = [];

	subscribe(listener: Listener) {
		this.listeners.push(listener);
		this.emit();

		return () => {
			this.listeners = this.listeners.filter((l) => l !== listener);
		};
	}

	private emit() {
		this.listeners.forEach((l) => l(this.items));
	}

	add(type: NotificationType, message: string) {
		this.items = this.items.filter((item) => !(item.type === type && item.message === message));

		const item: NotificationItem = {
			id: crypto.randomUUID(),
			type,
			message,
		};

		this.items = [item, ...this.items];
		this.emit();

		setTimeout(() => this.remove(item.id), 4000);
	}

	remove(id: string) {
		this.items = this.items.filter((i) => i.id !== id);
		this.emit();
	}
}

export const notify = {
	success: (msg: string) => notificationStore.add('success', msg),
	error: (msg: string) => notificationStore.add('error', msg),
};

export const notificationStore = new NotificationStore();
