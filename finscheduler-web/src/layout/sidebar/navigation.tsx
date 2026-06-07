import type {LucideIcon} from 'lucide-react';
import {LayoutDashboard, ShoppingBag, Tag} from 'lucide-react';
import {dashboardPath, itemsListPath, tagsListPath} from '../../features/routes.ts';

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
        label: 'Dashboard',
        path: dashboardPath,
        icon: LayoutDashboard,
    },
    {
        id: 'catalog',
        label: 'Каталог',
        path: itemsListPath,
        icon: ShoppingBag,
    },
    {
        id: 'tags',
        label: 'Теги',
        path: tagsListPath,
        icon: Tag,
    },
];
