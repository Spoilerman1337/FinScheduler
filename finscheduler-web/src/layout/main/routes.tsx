import Dashboard from '../../features/dashboard/Dashboard.tsx';
import ItemDetailsPage from '../../features/items/ItemDetailsPage.tsx';
import Items from '../../features/items/Items.tsx';
import {newItemPath} from '../../features/items/routes.ts';
import Tags from '../../features/tags/Tags.tsx';

export const routes = [
    {
        path: '/',
        title: 'Дашборды',
        element: <Dashboard />,
    },
    {
        path: '/items',
        title: 'Каталог',
        element: <Items />,
    },
    {
        path: newItemPath,
        title: 'Каталог',
        element: <ItemDetailsPage mode="create" />,
    },
    {
        path: '/items/:itemId/edit',
        title: 'Каталог',
        element: <ItemDetailsPage mode="edit" />,
    },
    {
        path: '/tags',
        title: 'Теги',
        element: <Tags />,
    },
];
