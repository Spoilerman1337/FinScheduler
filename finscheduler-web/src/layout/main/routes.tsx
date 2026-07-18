import type {RouteObject} from 'react-router-dom';
import Calendar from '../../features/calendar/Calendar.tsx';
import Dashboard from '../../features/dashboard/Dashboard.tsx';
import Events from '../../features/events/Events.tsx';
import {
    calendarPath,
    dashboardPath,
    eventsPath,
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
import {
    calendarRouteTitle,
    catalogNavigationLabel,
    dashboardRouteTitle,
    eventsRouteTitle,
    tagsNavigationLabel,
} from '../navigationLabels.ts';

export const mainRouteChildren: RouteObject[] = [
    {
        path: dashboardPath,
        handle: {
            title: dashboardRouteTitle,
        },
        element: <Dashboard />,
    },
    {
        path: calendarPath,
        handle: {
            title: calendarRouteTitle,
        },
        element: <Calendar />,
    },
    {
        path: eventsPath,
        handle: {
            title: eventsRouteTitle,
        },
        element: <Events />,
    },
    {
        path: itemsListPath,
        handle: {
            title: catalogNavigationLabel,
        },
        element: <Items />,
    },
    {
        path: newItemPath,
        handle: {
            title: catalogNavigationLabel,
        },
        element: <ItemDetailsPage mode="create" />,
    },
    {
        path: itemEditRoutePath,
        handle: {
            title: catalogNavigationLabel,
        },
        element: <ItemDetailsPage mode="edit" />,
    },
    {
        path: tagsListPath,
        handle: {
            title: tagsNavigationLabel,
        },
        element: <Tags />,
    },
    {
        path: newTagPath,
        handle: {
            title: tagsNavigationLabel,
        },
        element: <TagDetailsPage mode="create" />,
    },
    {
        path: tagEditRoutePath,
        handle: {
            title: tagsNavigationLabel,
        },
        element: <TagDetailsPage mode="edit" />,
    },
];
