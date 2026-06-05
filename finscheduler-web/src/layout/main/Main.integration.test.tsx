import {Flex} from '@chakra-ui/react';
import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {describe, expect, it, vi} from 'vitest';
import Main from './Main.tsx';
import Sidebar from '../sidebar/Sidebar.tsx';
import {renderWithProviders} from '../../test/render.tsx';

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

vi.mock('../../features/tags/Tags.tsx', () => ({
    default: function TagsPageMock() {
        return <div>Tags Page</div>;
    },
}));

function ShellUnderTest() {
    return (
        <Flex w="100vw" h="100vh">
            <Sidebar />
            <Main />
        </Flex>
    );
}

describe('Main and Sidebar integration', () => {
    it('renders the current route title, page content, and active sidebar link', () => {
        // Arrange
        renderWithProviders(<ShellUnderTest />, {route: '/items'});

        // Assert
        expect(screen.getByText('Items Page')).toBeInTheDocument();
        expect(screen.getAllByText('Каталог')).toHaveLength(2);
        expect(screen.getByRole('link', {name: 'Каталог', hidden: true})).toHaveAttribute(
            'href',
            '/items',
        );
        expect(screen.getByRole('link', {name: 'Каталог', hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('navigates between routes and updates the header and active link', async () => {
        // Arrange
        const user = userEvent.setup();
        renderWithProviders(<ShellUnderTest />, {route: '/items'});

        // Act
        await user.click(screen.getByRole('link', {name: 'Теги', hidden: true}));

        // Assert
        expect(screen.getByText('Tags Page')).toBeInTheDocument();
        expect(screen.getAllByText('Теги')).toHaveLength(2);
        expect(screen.queryByText('Items Page')).not.toBeInTheDocument();
        expect(screen.getByRole('link', {name: 'Теги', hidden: true})).toHaveAttribute(
            'aria-current',
            'page',
        );
    });

    it('keeps navigation usable after collapsing the sidebar', async () => {
        // Arrange
        const user = userEvent.setup();
        renderWithProviders(<ShellUnderTest />, {route: '/tags'});
        const collapseButton = screen.getByRole('button', {hidden: true});
        const dashboardLink = screen.getByRole('link', {name: 'Dashboard', hidden: true});

        // Act
        await user.click(collapseButton);
        await user.click(dashboardLink);

        // Assert
        expect(screen.getByText('Dashboard Page')).toBeInTheDocument();
        expect(screen.getByText('Дашборды')).toBeInTheDocument();
    });
});
