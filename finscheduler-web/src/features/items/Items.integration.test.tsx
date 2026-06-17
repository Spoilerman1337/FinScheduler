import {fireEvent, screen, waitFor, within} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import type {RouteObject} from 'react-router-dom';
import {describe, expect, it, vi} from 'vitest';
import {API_BASE_URL} from '../../config/api.ts';
import type {ItemDetailedDto, ItemListingDto} from '../../api/items.types.ts';
import {renderWithProviders} from '../../test/render.tsx';
import {renderWithDataRouter} from '../../test/renderDataRouter.tsx';
import {server} from '../../test/msw/server.ts';
import {itemEditRoutePath, itemsListPath, newItemPath} from '../routes.ts';
import ItemDetailsPage from './ItemDetailsPage.tsx';
import Items from './Items.tsx';

type ItemTestData = ItemListingDto & ItemDetailedDto;

function buildItem(overrides: Partial<ItemTestData> = {}): ItemTestData {
    return {
        id: 'item-1',
        name: 'Coffee',
        price: 199.5,
        description: 'Morning drink',
        isActive: true,
        updatedAt: null,
        cashback: 5,
        category: 'FoodDrinks',
        tags: [],
        priceHistory: [],
        ...overrides,
    };
}

function renderItemsRoutes(initialEntries: string[] = [itemsListPath]) {
    const routes: RouteObject[] = [
        {path: itemsListPath, element: <Items />},
        {path: newItemPath, element: <ItemDetailsPage mode="create" />},
        {path: itemEditRoutePath, element: <ItemDetailsPage mode="edit" />},
    ];

    return renderWithDataRouter(routes, {initialEntries});
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
        const requestCountBeforeReset = requests.length;

        // Act
        await user.click(screen.getByRole('button', {name: 'Сброс'}));

        // Assert
        expect(await screen.findByText('Coffee')).toBeInTheDocument();
        expect(screen.queryByText('Tea')).not.toBeInTheDocument();
        await waitFor(() => expect(requests).toHaveLength(requestCountBeforeReset));
    });

    it('applies the price range from a single filter control and delays reload until apply', async () => {
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
        await user.click(screen.getByRole('button', {name: 'Цена'}));
        await user.type(await screen.findByLabelText('Цена от'), '300');
        await user.type(screen.getByLabelText('Цена до'), '100');

        // Assert
        expect(requests).toHaveLength(1);

        // Act
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(await screen.findByText('Filtered by Price')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('priceFrom')).toBe('100');
            expect(lastRequest?.searchParams.get('priceTo')).toBe('300');
        });
        expect(screen.getByRole('button', {name: 'Цена: 100 - 300 ₽'})).toBeInTheDocument();
    });

    it('applies the cashback range from a single filter control and delays reload until apply', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);

                requests.push(url);

                if (
                    url.searchParams.get('cashbackFrom') === '5' &&
                    url.searchParams.get('cashbackTo') === '15'
                ) {
                    return HttpResponse.json({
                        data: [
                            buildItem({id: 'item-5', name: 'Filtered by Cashback', cashback: 10}),
                        ],
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
        await user.click(screen.getByRole('button', {name: 'Кэшбэк'}));
        await user.type(await screen.findByLabelText('Кэшбэк от'), '15');
        await user.type(screen.getByLabelText('Кэшбэк до'), '5');

        // Assert
        expect(requests).toHaveLength(1);

        // Act
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(await screen.findByText('Filtered by Cashback')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('cashbackFrom')).toBe('5');
            expect(lastRequest?.searchParams.get('cashbackTo')).toBe('15');
        });
        expect(screen.getByRole('button', {name: 'Кэшбэк: 5 - 15 %'})).toBeInTheDocument();
    });

    it('applies the created date range from a single filter control and delays reload until apply', async () => {
        // Arrange
        const requests: URL[] = [];
        const expectedCreatedFrom = new Date('2025-02-10T00:00:00.000').toISOString();
        const expectedCreatedTo = new Date('2025-02-15T23:59:59.999').toISOString();

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);

                requests.push(url);

                if (
                    url.searchParams.get('createdFrom') === expectedCreatedFrom &&
                    url.searchParams.get('createdTo') === expectedCreatedTo
                ) {
                    return HttpResponse.json({
                        data: [
                            buildItem({
                                id: 'item-6',
                                name: 'Filtered by Date',
                            }),
                        ],
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
        await user.click(screen.getByRole('button', {name: 'Обновлён'}));
        await user.click(screen.getByRole('button', {name: 'Создан'}));
        fireEvent.change(await screen.findByLabelText('Дата создания от'), {
            target: {value: '15.02.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Дата создания от'));
        fireEvent.change(screen.getByLabelText('Дата создания до'), {
            target: {value: '10.02.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Дата создания до'));

        // Assert
        expect(requests).toHaveLength(1);

        // Act
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(await screen.findByText('Filtered by Date')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('createdFrom')).toBe(expectedCreatedFrom);
            expect(lastRequest?.searchParams.get('createdTo')).toBe(expectedCreatedTo);
        });
        expect(
            screen.getByRole('button', {name: 'Создан: 10.02.2025 - 15.02.2025'}),
        ).toBeInTheDocument();
    });

    it('applies the updated date range after switching the active date mode', async () => {
        // Arrange
        const requests: URL[] = [];
        const expectedUpdatedFrom = new Date('2025-03-10T00:00:00.000').toISOString();
        const expectedUpdatedTo = new Date('2025-03-12T23:59:59.999').toISOString();

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);

                requests.push(url);

                if (
                    url.searchParams.get('updatedFrom') === expectedUpdatedFrom &&
                    url.searchParams.get('updatedTo') === expectedUpdatedTo
                ) {
                    return HttpResponse.json({
                        data: [
                            buildItem({
                                id: 'item-7',
                                name: 'Filtered by Updated Date',
                                updatedAt: '2025-03-11T09:00:00.000Z',
                            }),
                        ],
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
        await user.click(screen.getByRole('button', {name: 'Обновлён'}));
        fireEvent.change(await screen.findByLabelText('Дата обновления от'), {
            target: {value: '10.03.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Дата обновления от'));
        fireEvent.change(screen.getByLabelText('Дата обновления до'), {
            target: {value: '12.03.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Дата обновления до'));

        // Assert
        expect(requests).toHaveLength(1);

        // Act
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(await screen.findByText('Filtered by Updated Date')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('updatedFrom')).toBe(expectedUpdatedFrom);
            expect(lastRequest?.searchParams.get('updatedTo')).toBe(expectedUpdatedTo);
            expect(lastRequest?.searchParams.get('createdFrom')).toBeNull();
            expect(lastRequest?.searchParams.get('createdTo')).toBeNull();
        });
        expect(
            screen.getByRole('button', {name: 'Обновлён: 10.03.2025 - 12.03.2025'}),
        ).toBeInTheDocument();
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

    it('returns to the list after going back from a clean edit page', async () => {
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
        await user.click(screen.getByRole('button', {name: 'Назад'}));

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

    it('updates cashback by tag when no items are selected', async () => {
        // Arrange
        const currentItems = [
            buildItem({id: 'item-1', name: 'Coffee', cashback: 5}),
            buildItem({id: 'item-2', name: 'Tea', cashback: 7}),
        ];
        const tagPayloads: {tagId: string; cashback: number}[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: currentItems,
                    count: currentItems.length,
                });
            }),
            http.get(`${API_BASE_URL}/tags/lookup`, () => {
                return HttpResponse.json({
                    data: [{label: 'Groceries', value: 'tag-1'}],
                    count: 1,
                });
            }),
            http.patch(`${API_BASE_URL}/items/cashback/tag`, async ({request}) => {
                const payload = (await request.json()) as {tagId: string; cashback: number};

                tagPayloads.push(payload);
                currentItems.forEach((item) => {
                    item.cashback = payload.cashback;
                });

                return new HttpResponse(null, {status: 204});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('Coffee');
        await user.click(screen.getByRole('button', {name: 'Массовое обновление кешбека'}));

        const dialog = await screen.findByRole('dialog');
        await user.click(within(dialog).getByPlaceholderText('Выберите тег'));
        await user.click(await screen.findByText('Groceries'));
        const cashbackInput = within(dialog).getByRole('spinbutton', {name: 'Кешбек (%)'});
        await user.clear(cashbackInput);
        await user.type(cashbackInput, '12');
        await user.click(within(dialog).getByRole('button', {name: 'Сохранить'}));

        // Assert
        await waitFor(() => {
            expect(tagPayloads).toEqual([{tagId: 'tag-1', cashback: 12}]);
        });
        await waitFor(() => {
            expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
        });
        expect(
            (await screen.findAllByText((_, element) => element?.textContent?.trim() === '12%'))
                .length,
        ).toBeGreaterThanOrEqual(2);
    });

    it('updates cashback for selected items when rows are selected', async () => {
        // Arrange
        const currentItems = [
            buildItem({id: 'item-1', name: 'Coffee', cashback: 5}),
            buildItem({id: 'item-2', name: 'Tea', cashback: 7}),
        ];
        const itemPayloads: {itemIds: string[]; cashback: number}[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: currentItems,
                    count: currentItems.length,
                });
            }),
            http.patch(`${API_BASE_URL}/items/cashback/items`, async ({request}) => {
                const payload = (await request.json()) as {
                    itemIds: string[];
                    cashback: number;
                };

                itemPayloads.push(payload);
                currentItems.forEach((item) => {
                    if (payload.itemIds.includes(item.id ?? '')) {
                        item.cashback = payload.cashback;
                    }
                });

                return new HttpResponse(null, {status: 204});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items />);
        await screen.findByText('Coffee');
        const checkboxes = await screen.findAllByRole('checkbox');
        await user.click(checkboxes[1]);
        await user.click(screen.getByRole('button', {name: 'Массовое обновление кешбека'}));

        const dialog = await screen.findByRole('dialog');
        expect(within(dialog).getByText('Coffee')).toBeInTheDocument();
        expect(within(dialog).queryByLabelText('Тег')).not.toBeInTheDocument();

        const cashbackInput = within(dialog).getByRole('spinbutton', {name: 'Кешбек (%)'});
        await user.clear(cashbackInput);
        await user.type(cashbackInput, '18');
        await user.click(within(dialog).getByRole('button', {name: 'Сохранить'}));

        // Assert
        await waitFor(() => {
            expect(itemPayloads).toEqual([{itemIds: ['item-1'], cashback: 18}]);
        });
        expect(await screen.findByRole('button', {name: 'Добавить'})).toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Удалить (1)'})).not.toBeInTheDocument();
        expect(await screen.findByText('18%')).toBeInTheDocument();
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
