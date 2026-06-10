import {createBrowserRouter, type RouteObject} from 'react-router-dom';
import App from './App.tsx';
import {mainRouteChildren} from './layout/main/routes.tsx';

export const appRoutes: RouteObject[] = [
    {
        element: <App />,
        children: mainRouteChildren,
    },
];

export const appRouter = createBrowserRouter(appRoutes);
