import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import {describe, expect, it} from 'vitest';
import {Route, Routes} from 'react-router-dom';
import type {ItemModification} from '../../api/types.ts';
import {API_BASE_URL} from '../../config/api.ts';
import {renderWithProviders} from '../../test/render.tsx';
import {server} from '../../test/msw/server.ts';
import {buildEditItemPath, itemEditRoutePath, itemsListPath, newItemPath} from '../routes.ts';
import ItemDetailsPage from './ItemDetailsPage.tsx';

function renderItemDetailsRoutes(initialEntries: string[]) {
    return renderWithProviders(
        <Routes>
            <Route path={itemsListPath} element={<div>Items Listing Page</div>} />
            <Route path={newItemPath} element={<ItemDetailsPage mode="create" />} />
            <Route path={itemEditRoutePath} element={<ItemDetailsPage mode="edit" />} />
        </Routes>,
        {initialEntries},
    );
}

describe('ItemDetailsPage integration', () => {
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
});
