import { Navigate, Outlet } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { ROUTES } from '@/constants/constants';
import { getAuthState } from '@/private-route';
import { readAuthSession } from '@/utils/authSessionUtils';

export const ProtectedRoute = () => {
	const [isLoggedIn, setIsLoggedIn] = useState(readAuthSession);

	useEffect(() => {
		void getAuthState().then(setIsLoggedIn);
	}, []);

	if (!isLoggedIn) {
		return <Navigate to={ROUTES.START} replace />;
	}

	return <Outlet />;
};

export const PublicRoute = ({ redirectIfLogged = true }: { redirectIfLogged?: boolean }) => {
	const [isLoggedIn, setIsLoggedIn] = useState(readAuthSession);

	useEffect(() => {
		void getAuthState().then(setIsLoggedIn);
	}, []);

	if (isLoggedIn && redirectIfLogged) {
		return <Navigate to={ROUTES.HOME} replace />;
	}

	return <Outlet />;
};
