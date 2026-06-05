import type {LucideIcon} from 'lucide-react';
import {LayoutDashboard, ShoppingBag, Tag} from 'lucide-react';

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
        path: '/',
        icon: LayoutDashboard,
    },
    {
        id: 'catalog',
        label: 'Каталог',
        path: '/items',
        icon: ShoppingBag,
    },
    {
        id: 'tags',
        label: 'Теги',
        path: '/tags',
        icon: Tag,
    },
];
