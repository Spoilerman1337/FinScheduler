import Dashboard from '../../features/dashboard/Dashboard.tsx';
import {
    dashboardPath,
    itemEditRoutePath,
    itemsListPath,
    newItemPath,
    newTagPath,
    tagEditRoutePath,
    tagsListPath,
} from '../../features/routes.ts';
import ItemDetailsPage from '../../features/items/ItemDetailsPage.tsx';
import Items from '../../features/items/Items.tsx';
import TagDetailsPage from '../../features/tags/TagDetailsPage.tsx';
import Tags from '../../features/tags/Tags.tsx';

export const routes = [
    {
        path: dashboardPath,
        title: 'Дашборды',
        element: <Dashboard />,
    },
    {
        path: itemsListPath,
        title: 'Каталог',
        element: <Items />,
    },
    {
        path: newItemPath,
        title: 'Каталог',
        element: <ItemDetailsPage mode="create" />,
    },
    {
        path: itemEditRoutePath,
        title: 'Каталог',
        element: <ItemDetailsPage mode="edit" />,
    },
    {
        path: tagsListPath,
        title: 'Теги',
        element: <Tags />,
    },
    {
        path: newTagPath,
        title: 'Теги',
        element: <TagDetailsPage mode="create" />,
    },
    {
        path: tagEditRoutePath,
        title: 'Теги',
        element: <TagDetailsPage mode="edit" />,
    },
];
