import { writable } from 'svelte/store';

interface Notification {
	id: number;
	type: 'success' | 'error' | 'info';
	message: string;
}

function createNotificationStore() {
	const { subscribe, update } = writable<Notification[]>([]);
	let id = 0;

	return {
		subscribe,
		add: (type: 'success' | 'error' | 'info', message: string) => {
			const notification = { id: id++, type, message };
			update((n) => [...n, notification]);
		},
		remove: (id: number) => {
			update((n) => n.filter((x) => x.id !== id));
		},
		success: (message: string) => {
			notificationStore.add('success', message);
		},
		error: (message: string) => {
			notificationStore.add('error', message);
		},
		info: (message: string) => {
			notificationStore.add('info', message);
		}
	};
}

export const notificationStore = createNotificationStore();
