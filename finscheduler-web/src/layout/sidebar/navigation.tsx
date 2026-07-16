import type {LucideIcon} from 'lucide-react';
import {CalendarDays, LayoutDashboard, ShoppingBag, Tag} from 'lucide-react';
import {calendarPath, dashboardPath, itemsListPath, tagsListPath} from '../../features/routes.ts';
import {
    calendarNavigationLabel,
    catalogNavigationLabel,
    dashboardNavigationLabel,
    tagsNavigationLabel,
} from '../navigationLabels.ts';

export interface NavigationItem {
    id: string;
    label: string;
    path?: string;
    icon: LucideIcon;
    disabled?: boolean;
}

export const routedNavigationItems: NavigationItem[] = [
    {
        id: 'dashboard',
        label: dashboardNavigationLabel,
        path: dashboardPath,
        icon: LayoutDashboard,
    },
    {
        id: 'calendar',
        label: calendarNavigationLabel,
        path: calendarPath,
        icon: CalendarDays,
    },
    {
        id: 'catalog',
        label: catalogNavigationLabel,
        path: itemsListPath,
        icon: ShoppingBag,
    },
    {
        id: 'tags',
        label: tagsNavigationLabel,
        path: tagsListPath,
        icon: Tag,
    },
];
