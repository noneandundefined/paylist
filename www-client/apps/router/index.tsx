import { Route, Routes } from 'react-router-dom';
import config from './config';
import RoutePage from './RoutePage';
import { ProtectedRoute, PublicRoute } from './routeGuards';

const protectedRoutes = config.filter((route) => route.loginRequired);
const publicRoutesWithRedirect = config.filter((route) => !route.loginRequired && route.redirectIfLogged !== false);
const publicRoutesWithoutRedirect = config.filter((route) => !route.loginRequired && route.redirectIfLogged === false);

const AppRouter: React.FC = () => {
	return (
		<Routes>
			<Route element={<ProtectedRoute />}>
				{protectedRoutes.map((route) => (
					<Route key={route.path} path={route.path} element={<RoutePage title={route.title} component={route.component} />} />
				))}
			</Route>

			<Route element={<PublicRoute redirectIfLogged />}>
				{publicRoutesWithRedirect.map((route) => (
					<Route key={route.path} path={route.path} element={<RoutePage title={route.title} component={route.component} />} />
				))}
			</Route>

			<Route element={<PublicRoute redirectIfLogged={false} />}>
				{publicRoutesWithoutRedirect.map((route) => (
					<Route key={route.path} path={route.path} element={<RoutePage title={route.title} component={route.component} />} />
				))}
			</Route>
		</Routes>
	);
};

export default AppRouter;
