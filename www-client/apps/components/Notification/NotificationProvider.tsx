import { useEffect, useState } from 'react';
import { notificationStore, type NotificationItem } from './notify';

const NotificationProvider = () => {
	const [items, setItems] = useState<NotificationItem[]>([]);

	useEffect(() => {
		return notificationStore.subscribe(setItems);
	}, []);

	return (
		<div className="fixed bottom-4 right-4 space-y-2 z-[10005]">
			{items.map((n) => (
				<div
					key={n.id}
					onClick={() => notificationStore.remove(n.id)}
					className={`
						px-5 py-3 rounded shadow text-white text-sm
                        max-w-xl cursor-pointer
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
