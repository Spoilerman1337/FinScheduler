import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {describe, expect, it, vi} from 'vitest';
import {
    buildEditItemPath,
    buildEditTagPath,
    dashboardPath,
    itemsListPath,
    tagsListPath,
} from '../../features/routes.ts';
import {
    catalogNavigationLabel,
    dashboardNavigationLabel,
    dashboardRouteTitle,
    tagsNavigationLabel,
} from '../navigationLabels.ts';
import {appRoutes} from '../../router.tsx';
import {renderWithDataRouter} from '../../test/renderDataRouter.tsx';

vi.mock('framer-motion', async () => {
    const React = await import('react');

    const MotionDiv = React.forwardRef<
        HTMLDivElement,
        React.HTMLAttributes<HTMLDivElement> & {
            animate?: unknown;
            initial?: unknown;
            exit?: unknown;
            layout?: unknown;
            transition?: unknown;
        }
    >(function MotionDiv({animate, initial, exit, layout, transition, ...props}, ref) {
        void animate;
        void initial;
        void exit;
        void layout;
        void transition;

        return <div ref={ref} {...props} />;
    });

    return {
        AnimatePresence: ({children}: {children: React.ReactNode}) => <>{children}</>,
        motion: {
            div: MotionDiv,
            create: <TProps extends object>(Component: React.ComponentType<TProps>) =>
                React.forwardRef<unknown, TProps & {layout?: unknown}>(function MotionComponent(
                    {layout, ...props},
                    ref,
                ) {
                    void layout;

                    return <Component ref={ref as never} {...(props as TProps)} />;
                }),
        },
    };
});

vi.mock('../../features/dashboard/Dashboard.tsx', () => ({
    default: function DashboardPageMock() {
        return <div>Dashboard Page</div>;
    },
}));

vi.mock('../../features/items/Items.tsx', () => ({
    default: function ItemsPageMock() {
        return <div>Items Page</div>;
    },
}));

vi.mock('../../features/items/ItemDetailsPage.tsx', () => ({
    default: function ItemDetailsPageMock() {
        return <div>Item Details Page</div>;
    },
}));

vi.mock('../../features/tags/Tags.tsx', () => ({
    default: function TagsPageMock() {
        return <div>Tags Page</div>;
    },
}));

vi.mock('../../features/tags/TagDetailsPage.tsx', () => ({
    default: function TagDetailsPageMock() {
        return <div>Tag Details Page</div>;
    },
}));

vi.mock('../../components/ui/toaster.tsx', () => {
    const ToasterMock: typeof import('../../components/ui/toaster.tsx').Toaster = () => <></>;

    return {
        Toaster: ToasterMock,
    };
});

describe('Main and Sidebar integration', () => {
    it('renders the current route title, page content, and active sidebar link', () => {
        // Arrange
        renderWithDataRouter(appRoutes, {initialEntries: [itemsListPath]});

        // Assert
        expect(screen.getByText('Items Page')).toBeInTheDocument();
        expect(screen.getAllByText(catalogNavigationLabel)).toHaveLength(2);
        expect(screen.getByRole('link', {name: catalogNavigationLabel, hidden: true})).toHaveAttribute(
            'href',
            itemsListPath,
        );
        expect(screen.getByRole('link', {name: catalogNavigationLabel, hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('navigates between routes and updates the header and active link', async () => {
        // Arrange
        const user = userEvent.setup();
        renderWithDataRouter(appRoutes, {initialEntries: [itemsListPath]});

        // Act
        await user.click(screen.getByRole('link', {name: tagsNavigationLabel, hidden: true}));

        // Assert
        expect(await screen.findByText('Tags Page')).toBeInTheDocument();
        expect(await screen.findAllByText(tagsNavigationLabel)).toHaveLength(2);
        expect(screen.queryByText('Items Page')).not.toBeInTheDocument();
        expect(screen.getByRole('link', {name: tagsNavigationLabel, hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('keeps catalog route active for nested item detail pages', () => {
        // Arrange
        renderWithDataRouter(appRoutes, {initialEntries: [buildEditItemPath('item-1')]});

        // Assert
        expect(screen.getByText('Item Details Page')).toBeInTheDocument();
        expect(screen.getAllByText(catalogNavigationLabel)).toHaveLength(2);
        expect(screen.getByRole('link', {name: catalogNavigationLabel, hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('keeps tags route active for nested tag detail pages', () => {
        // Arrange
        renderWithDataRouter(appRoutes, {initialEntries: [buildEditTagPath('tag-1')]});

        // Assert
        expect(screen.getByText('Tag Details Page')).toBeInTheDocument();
        expect(screen.getAllByText(tagsNavigationLabel)).toHaveLength(2);
        expect(screen.getByRole('link', {name: tagsNavigationLabel, hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('keeps navigation usable after collapsing the sidebar', async () => {
        // Arrange
        const user = userEvent.setup();
        renderWithDataRouter(appRoutes, {initialEntries: [tagsListPath]});

        // Act
        await user.click(screen.getByRole('button', {hidden: true}));
        await user.click(screen.getByRole('link', {name: dashboardNavigationLabel, hidden: true}));

        // Assert
        expect(await screen.findByText('Dashboard Page')).toBeInTheDocument();
        expect(screen.getByText(dashboardRouteTitle)).toBeInTheDocument();
        expect(screen.getByRole('link', {name: dashboardNavigationLabel, hidden: true})).toHaveAttribute(
            'href',
            dashboardPath,
        );
    });
});
