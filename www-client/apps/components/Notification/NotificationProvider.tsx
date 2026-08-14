import { useEffect, useState } from 'react';
import { notificationStore, type NotificationItem } from './notify';

const NotificationProvider = () => {
	const [items, setItems] = useState<NotificationItem[]>([]);

	useEffect(() => {
		return notificationStore.subscribe(setItems);
	}, []);

	return (
		<div className="fixed bottom-4 left-1/2 -translate-x-1/2 w-[calc(100%-2rem)] max-w-xl sm:left-auto sm:translate-x-0 sm:right-4 sm:w-auto space-y-2 z-[10005]">
			{items.map((n) => (
				<div
					key={n.id}
					onClick={() => notificationStore.remove(n.id)}
					className={`
						px-5 py-3 rounded shadow text-white text-base sm:text-sm
						w-full cursor-pointer text-center sm:text-left
						${n.type === 'success' ? 'bg-green-500' : ''}
						${n.type === 'error' ? 'bg-red-500' : ''}
					`}
				>
					{n.message}
				</div>
			))}
		</div>
	);
};

export default NotificationProvider;
