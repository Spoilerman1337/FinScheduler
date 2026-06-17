import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import type {RouteObject} from 'react-router-dom';
import {describe, expect, it} from 'vitest';
import type {ItemModification} from '../../api/items.types.ts';
import {API_BASE_URL} from '../../config/api.ts';
import {renderWithDataRouter} from '../../test/renderDataRouter.tsx';
import {server} from '../../test/msw/server.ts';
import {buildEditItemPath, itemEditRoutePath, itemsListPath, newItemPath} from '../routes.ts';
import ItemDetailsPage from './ItemDetailsPage.tsx';

function renderItemDetailsRoutes(initialEntries: string[]) {
    const routes: RouteObject[] = [
        {path: itemsListPath, element: <div>Items Listing Page</div>},
        {path: newItemPath, element: <ItemDetailsPage mode="create" />},
        {path: itemEditRoutePath, element: <ItemDetailsPage mode="edit" />},
    ];

    return renderWithDataRouter(routes, {
        initialEntries,
    });
}

describe('ItemDetailsPage integration', () => {
    it('shows only the back action while the form is clean', async () => {
        // Arrange
        renderItemDetailsRoutes([newItemPath]);

        // Act
        const backButton = screen.getByRole('button', {name: 'Назад'});

        // Assert
        expect(backButton).toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить'})).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить и закрыть'})).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Отмена'})).not.toBeInTheDocument();
    });

    it('shows a validation summary and inline field errors without sending a request', async () => {
        // Arrange
        let createRequests = 0;

        server.use(
            http.post(`${API_BASE_URL}/items`, async () => {
                createRequests += 1;

                return HttpResponse.json('item-2');
            }),
        );

        const user = userEvent.setup();

        renderItemDetailsRoutes([newItemPath]);

        // Act
        await user.type(screen.getByLabelText('Название'), '   ');
        await user.click(screen.getByRole('button', {name: 'Сохранить'}));

        // Assert
        expect(createRequests).toBe(0);
        expect(screen.getByText('Ошибка валидации')).toBeInTheDocument();
        expect(screen.getByText('Название обязательно для заполнения')).toBeInTheDocument();
        expect(
            screen.getByText('Выберите категорию', {selector: '[data-part="error-text"]'}),
        ).toBeInTheDocument();
    });

    it('creates a new item and stays on the detail page after saving', async () => {
        // Arrange
        let createdPayload: ItemModification | null = null;

        server.use(
            http.post(`${API_BASE_URL}/items`, async ({request}) => {
                createdPayload = (await request.json()) as ItemModification;

                return HttpResponse.json('item-2');
            }),
            http.get(`${API_BASE_URL}/items/item-2`, () => {
                return HttpResponse.json({
                    id: 'item-2',
                    name: 'New Item',
                    description: undefined,
                    price: 0,
                    cashback: 0,
                    isActive: true,
                    category: 'FoodDrinks',
                    tags: [],
                    priceHistory: [],
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemDetailsRoutes([newItemPath]);
        await user.type(screen.getByLabelText('Название'), 'New Item');
        await user.click(screen.getByText('Выберите категорию'));
        await user.click(await screen.findByRole('option', {name: 'Еда и напитки'}));
        await user.click(screen.getByRole('button', {name: 'Сохранить'}));

        // Assert
        await waitFor(() => {
            expect(createdPayload).toEqual({
                name: 'New Item',
                description: undefined,
                price: 0,
                cashback: 0,
                isActive: true,
                category: 'FoodDrinks',
                tagIds: [],
            });
        });
        expect(await screen.findByText('Редактирование предмета')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('New Item');
        expect(screen.getByRole('button', {name: 'Назад'})).toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить'})).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить и закрыть'})).not.toBeInTheDocument();
    });

    it('updates an existing item and returns to the list after save and close', async () => {
        // Arrange
        let updatedPayload: ItemModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/items/item-1`, () => {
                return HttpResponse.json({
                    id: 'item-1',
                    name: 'Old Item',
                    description: 'Morning drink',
                    price: 199.5,
                    cashback: 5,
                    isActive: true,
                    category: 'FoodDrinks',
                    tags: [],
                    priceHistory: [],
                });
            }),
            http.put(`${API_BASE_URL}/items/item-1`, async ({request}) => {
                updatedPayload = (await request.json()) as ItemModification;

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderItemDetailsRoutes([buildEditItemPath('item-1')]);
        await screen.findByText('Редактирование предмета');
        expect(screen.getByRole('button', {name: 'Назад'})).toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить'})).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сохранить и закрыть'})).not.toBeInTheDocument();
        await user.clear(screen.getByLabelText('Название'));
        await user.type(screen.getByLabelText('Название'), 'Updated Item');
        await user.click(screen.getByRole('button', {name: 'Сохранить и закрыть'}));

        // Assert
        await waitFor(() => {
            expect(updatedPayload).toEqual({
                name: 'Updated Item',
                description: 'Morning drink',
                price: 199.5,
                cashback: 5,
                isActive: true,
                category: 'FoodDrinks',
                tagIds: [],
            });
        });
        expect(await screen.findByText('Items Listing Page')).toBeInTheDocument();
    });

    it('renders the price history chart and summary for an existing item', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/items/item-1`, () => {
                return HttpResponse.json({
                    id: 'item-1',
                    name: 'History Item',
                    description: 'Tracked over time',
                    price: 160,
                    cashback: 3,
                    isActive: true,
                    category: 'FoodDrinks',
                    tags: [],
                    priceHistory: [
                        {point: '2026-01-12T00:00:00Z', value: 139.9},
                        {point: '2026-02-18T00:00:00Z', value: 149.5},
                        {point: '2026-03-20T00:00:00Z', value: 160},
                    ],
                });
            }),
        );

        // Act
        renderItemDetailsRoutes([buildEditItemPath('item-1')]);

        // Assert
        expect(await screen.findByText('История цены')).toBeInTheDocument();
        expect(screen.getByLabelText('Минимум: 139,90 ₽')).toBeInTheDocument();
        expect(screen.getByLabelText('Текущая: 160,00 ₽')).toBeInTheDocument();
        expect(screen.getByLabelText('Максимум: 160,00 ₽')).toBeInTheDocument();
        expect(screen.queryByText('Точек')).not.toBeInTheDocument();
        expect(screen.getByLabelText('График истории цены')).toBeInTheDocument();
    });

    it('shows a warning on cancel when the form has unsaved changes', async () => {
        // Arrange
        const user = userEvent.setup();

        renderItemDetailsRoutes([newItemPath]);

        // Act
        await user.type(screen.getByLabelText('Название'), 'Draft Item');
        await user.click(screen.getByRole('button', {name: 'Отмена'}));

        // Assert
        expect(await screen.findByRole('dialog')).toBeInTheDocument();
        expect(screen.getByText('Есть несохранённые изменения')).toBeInTheDocument();
        expect(screen.queryByText('Items Listing Page')).not.toBeInTheDocument();

        // Act
        await user.click(screen.getByRole('button', {name: 'Остаться'}));

        // Assert
        await waitFor(() => {
            expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
        });
        expect(screen.getByLabelText('Название')).toHaveValue('Draft Item');
    });
});
