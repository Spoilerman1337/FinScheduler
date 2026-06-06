import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import {describe, expect, it, vi} from 'vitest';
import {Route, Routes} from 'react-router-dom';
import {API_BASE_URL} from '../../config/api.ts';
import type {ItemDto} from '../../api/types.ts';
import {renderWithProviders} from '../../test/render.tsx';
import {server} from '../../test/msw/server.ts';
import ItemDetailsPage from './ItemDetailsPage.tsx';
import Items from './Items.tsx';
import {itemsListPath, newItemPath} from './routes.ts';

function buildItem(overrides: Partial<ItemDto> = {}): ItemDto {
    return {
        id: 'item-1',
        name: 'Coffee',
        price: 199.5,
        description: 'Morning drink',
        isActive: true,
        cashback: 5,
        category: 'FoodDrinks',
        createdAt: '2025-01-01T00:00:00.000Z',
        tags: [],
        ...overrides,
    };
}

function renderItemsRoutes(initialEntries: string[] = [itemsListPath]) {
    return renderWithProviders(
        <Routes>
            <Route path={itemsListPath} element={<Items />} />
            <Route path={newItemPath} element={<ItemDetailsPage mode="create" />} />
            <Route path="/items/:itemId/edit" element={<ItemDetailsPage mode="edit" />} />
        </Routes>,
        {initialEntries},
    );
}

describe('Items integration', () => {
    it('loads active items on mount and renders the first page', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                requests.push(new URL(request.url));

                return HttpResponse.json({
                    data: [buildItem({name: 'Utility Bill'})],
                    count: 1,
                });
            }),
        );

        // Act
        renderWithProviders(<Items />);

        // Assert
        expect(await screen.findByText('Utility Bill')).toBeInTheDocument();
        await waitFor(() => expect(requests).toHaveLength(1));
        expect(requests[0].searchParams.get('page')).toBe('0');
        expect(requests[0].searchParams.get('pageSize')).toBe('12');
        expect(requests[0].searchParams.get('isActive')).toBe('true');
    });

    it('applies search and status filters, then resets them', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);
                const name = url.searchParams.get('name');
                const isActive = url.searchParams.get('isActive');

                requests.push(url);

                if (name === 'Tea') {
                    return HttpResponse.json({
                        data: [buildItem({id: 'item-2', name: 'Tea'})],
                        count: 1,
                    });
                }

                if (isActive === 'false') {
                    return HttpResponse.json({
                        data: [buildItem({id: 'item-3', name: 'Rent', isActive: false})],
                        count: 1,
                    });
                }

                if (isActive === 'true') {
                    return HttpResponse.json({
                        data: [buildItem({name: 'Coffee'})],
                        count: 1,
                    });
                }

                return HttpResponse.json({
                    data: [buildItem({id: 'item-4', name: 'All Results'})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('Coffee');
        await user.click(screen.getByRole('button', {name: 'Неактивные'}));

        // Assert
        expect(await screen.findByText('Rent')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('isActive')).toBe('false');
        });

        // Act
        await user.type(screen.getByPlaceholderText('Поиск по названию...'), 'Tea');

        // Assert
        expect(await screen.findByText('Tea')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('name')).toBe('Tea');
            expect(lastRequest?.searchParams.get('isActive')).toBe('false');
        });

        // Act
        await user.click(screen.getByRole('button', {name: 'Сброс'}));

        // Assert
        expect(await screen.findByText('All Results')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('name')).toBeNull();
            expect(lastRequest?.searchParams.get('isActive')).toBeNull();
        });
    });

    it('applies price filters and reloads the table', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);

                requests.push(url);

                if (
                    url.searchParams.get('priceFrom') === '100' &&
                    url.searchParams.get('priceTo') === '300'
                ) {
                    return HttpResponse.json({
                        data: [buildItem({id: 'item-2', name: 'Filtered by Price', price: 250})],
                        count: 1,
                    });
                }

                return HttpResponse.json({
                    data: [buildItem({name: 'Coffee'})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('Coffee');
        await user.type(screen.getByPlaceholderText('Цена от'), '100');
        await user.type(screen.getByPlaceholderText('Цена до'), '300{Enter}');

        // Assert
        expect(await screen.findByText('Filtered by Price')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('priceFrom')).toBe('100');
            expect(lastRequest?.searchParams.get('priceTo')).toBe('300');
        });
    });

    it('navigates to the create page from the add button', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Existing Item'})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemsRoutes();
        await screen.findByText('Existing Item');
        await user.click(screen.getByRole('button', {name: 'Добавить'}));

        // Assert
        expect(await screen.findByText('Новый предмет')).toBeInTheDocument();
        expect(screen.getByText('Создание')).toBeInTheDocument();
    });

    it('navigates to the edit page after a double click', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Old Item'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/items/item-1`, () => {
                return HttpResponse.json(buildItem({name: 'Old Item'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemsRoutes();
        await user.dblClick(await screen.findByText('Old Item'));

        // Assert
        expect(await screen.findByText('Редактирование предмета')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('Old Item');
    });

    it('opens the edit page from the showcase action button', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Coffee'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/items/item-1`, () => {
                return HttpResponse.json(buildItem({name: 'Coffee'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemsRoutes();
        await user.click((await screen.findAllByRole('button', {name: 'Открыть карточку'}))[0]);

        // Assert
        expect(await screen.findByText('Редактирование предмета')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('Coffee');
    });

    it('returns to the list after cancelling from the edit page', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Coffee'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/items/item-1`, () => {
                return HttpResponse.json(buildItem({name: 'Coffee'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemsRoutes();
        await user.dblClick(await screen.findByText('Coffee'));
        expect(await screen.findByText('Редактирование предмета')).toBeInTheDocument();
        await user.click(screen.getByRole('button', {name: 'Отмена'}));

        // Assert
        expect(await screen.findByText('Coffee')).toBeInTheDocument();
        expect(screen.queryByText('Редактирование предмета')).not.toBeInTheDocument();
    });

    it('deletes selected items and reloads the table', async () => {
        // Arrange
        const currentItems = [
            buildItem({id: 'item-1', name: 'Coffee'}),
            buildItem({id: 'item-2', name: 'Tea'}),
        ];
        const deletedIds: string[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: currentItems,
                    count: currentItems.length,
                });
            }),
            http.delete(`${API_BASE_URL}/items/:id`, ({params}) => {
                const id = String(params.id);
                const itemIndex = currentItems.findIndex((item) => item.id === id);

                deletedIds.push(id);

                if (itemIndex >= 0) {
                    currentItems.splice(itemIndex, 1);
                }

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('Coffee');
        const checkboxes = await screen.findAllByRole('checkbox');

        await user.click(checkboxes[1]);

        // Assert
        expect(screen.queryByRole('button', {name: 'Добавить'})).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Редактировать'})).not.toBeInTheDocument();

        // Act
        await user.click(screen.getByRole('button', {name: 'Удалить (1)'}));

        // Assert
        await waitFor(() => {
            expect(deletedIds).toEqual(['item-1']);
        });
        expect(screen.queryByText('Coffee')).not.toBeInTheDocument();
        expect(await screen.findByText('Tea')).toBeInTheDocument();
    });

    it('changes the page when the paginator is used', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get('page');

                requests.push(url);

                return HttpResponse.json({
                    data:
                        page === '1'
                            ? [buildItem({id: 'item-2', name: 'Second Page Item'})]
                            : [buildItem({name: 'First Page Item'})],
                    count: 13,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('First Page Item');
        await user.click(screen.getByRole('button', {name: 'last page, page 2'}));

        // Assert
        expect(await screen.findByText('Second Page Item')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('page')).toBe('1');
        });
    });

    it('changes the page size and reloads the first page', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get('page');
                const pageSize = url.searchParams.get('pageSize');

                requests.push(url);

                if (pageSize === '24') {
                    return HttpResponse.json({
                        data: [buildItem({id: 'item-3', name: 'Twenty Four Per Page'})],
                        count: 30,
                    });
                }

                if (page === '1') {
                    return HttpResponse.json({
                        data: [buildItem({id: 'item-2', name: 'Second Page Item'})],
                        count: 30,
                    });
                }

                return HttpResponse.json({
                    data: [buildItem({name: 'First Page Item'})],
                    count: 30,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('First Page Item');
        await user.click(screen.getByRole('button', {name: /page 2/i}));
        await screen.findByText('Second Page Item');
        await user.click(screen.getByRole('combobox'));
        await user.click(await screen.findByRole('option', {name: '24'}));

        // Assert
        expect(await screen.findByText('Twenty Four Per Page')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('page')).toBe('0');
            expect(lastRequest?.searchParams.get('pageSize')).toBe('24');
        });
    });

    it('loads tag lookup options when the tags field is opened', async () => {
        // Arrange
        const lookupRequests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Existing Item'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/tags/lookup`, ({request}) => {
                lookupRequests.push(new URL(request.url));

                return HttpResponse.json({
                    data: [
                        {label: 'Travel', value: 'tag-1'},
                        {label: 'Home', value: 'tag-2'},
                    ],
                    count: 2,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemsRoutes();
        await screen.findByText('Existing Item');
        await user.click(screen.getByRole('button', {name: 'Добавить'}));
        await user.click(screen.getByPlaceholderText('Выберите теги'));

        // Assert
        expect(await screen.findByText('Travel')).toBeInTheDocument();
        await waitFor(() => {
            expect(lookupRequests).toHaveLength(1);
        });
        expect(lookupRequests[0].searchParams.get('page')).toBe('0');
        expect(lookupRequests[0].searchParams.get('pageSize')).toBe('20');
        expect(lookupRequests[0].searchParams.get('name')).toBeNull();
    });

    it('shows an empty state when no items are returned', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [],
                    count: 0,
                });
            }),
        );

        // Act
        renderWithProviders(<Items />);

        // Assert
        expect(await screen.findByText('Данные не найдены.')).toBeInTheDocument();
    });

    it('renders ruble prices with the currency sign as a suffix', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({price: 1234.5})],
                    count: 1,
                });
            }),
        );

        // Act
        renderWithProviders(<Items />);

        // Assert
        expect(
            (
                await screen.findAllByText((_, element) => {
                    const normalizedText = element?.textContent?.replace(/\s+/g, ' ').trim();
                    return normalizedText === '1 234,50 ₽';
                })
            ).length,
        ).toBeGreaterThan(0);
    });

    it('shows an error state when items request fails', async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json(
                    {message: 'boom'},
                    {status: 500, statusText: 'Internal Server Error'},
                );
            }),
        );

        try {
            // Act
            renderWithProviders(<Items />);

            // Assert
            expect(
                await screen.findByText('Failed to fetch items: Internal Server Error'),
            ).toBeInTheDocument();
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });

    it('shows a save error when creating an item fails', async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: [buildItem({name: 'Existing Item'})],
                    count: 1,
                });
            }),
            http.post(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json(
                    {message: 'boom'},
                    {status: 500, statusText: 'Internal Server Error'},
                );
            }),
        );

        const user = userEvent.setup();

        try {
            // Act
            renderItemsRoutes();
            await screen.findByText('Existing Item');
            await user.click(screen.getByRole('button', {name: 'Добавить'}));
            await user.type(screen.getByLabelText('Название'), 'Broken Item');
            await user.click(screen.getByText('Выберите категорию'));
            await user.click(await screen.findByRole('option', {name: 'Еда и напитки'}));
            await user.click(screen.getByRole('button', {name: 'Сохранить'}));

            // Assert
            expect(
                await screen.findByText('Failed to create item: Internal Server Error'),
            ).toBeInTheDocument();
            expect(screen.getByLabelText('Название')).toHaveValue('Broken Item');
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });
});
