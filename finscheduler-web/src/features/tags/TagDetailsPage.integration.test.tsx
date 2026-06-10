import {screen, waitFor, within} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import {describe, expect, it} from 'vitest';
import type {RouteObject} from 'react-router-dom';
import type {TagModification} from '../../api/tags.types.ts';
import {API_BASE_URL} from '../../config/api.ts';
import {renderWithDataRouter} from '../../test/renderDataRouter.tsx';
import {server} from '../../test/msw/server.ts';
import {buildEditTagPath, newTagPath, tagEditRoutePath, tagsListPath} from '../routes.ts';
import TagDetailsPage from './TagDetailsPage.tsx';

function renderTagDetailsRoutes(initialEntries: string[]) {
    const routes: RouteObject[] = [
        {path: tagsListPath, element: <div>Tags Listing Page</div>},
        {path: newTagPath, element: <TagDetailsPage mode="create" />},
        {path: tagEditRoutePath, element: <TagDetailsPage mode="edit" />},
    ];

    return renderWithDataRouter(routes, {
        initialEntries,
    });
}

describe('TagDetailsPage integration', () => {
    it('creates a new tag and stays on the detail page after saving', async () => {
        // Arrange
        let createdPayload: TagModification | null = null;

        server.use(
            http.post(`${API_BASE_URL}/tags`, async ({request}) => {
                createdPayload = (await request.json()) as TagModification;

                return HttpResponse.json('tag-2');
            }),
            http.get(`${API_BASE_URL}/tags/tag-2`, () => {
                return HttpResponse.json({
                    id: 'tag-2',
                    name: 'New Tag',
                    isActive: true,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagDetailsRoutes([newTagPath]);
        await user.type(screen.getByLabelText('Название'), 'New Tag');
        await user.click(screen.getByRole('button', {name: 'Сохранить'}));

        // Assert
        await waitFor(() => {
            expect(createdPayload).toEqual({
                name: 'New Tag',
                isActive: true,
            });
        });
        expect(await screen.findByText('Редактирование тега')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('New Tag');
    });

    it('updates an existing tag and returns to the list after save and close', async () => {
        // Arrange
        let updatedPayload: TagModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json({
                    id: 'tag-1',
                    name: 'Old Tag',
                    isActive: true,
                });
            }),
            http.put(`${API_BASE_URL}/tags/tag-1`, async ({request}) => {
                updatedPayload = (await request.json()) as TagModification;

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagDetailsRoutes([buildEditTagPath('tag-1')]);
        await screen.findByText('Редактирование тега');
        await user.clear(screen.getByLabelText('Название'));
        await user.type(screen.getByLabelText('Название'), 'Updated Tag');
        await user.click(screen.getByRole('button', {name: 'Сохранить и закрыть'}));

        // Assert
        await waitFor(() => {
            expect(updatedPayload).toEqual({
                name: 'Updated Tag',
                isActive: true,
            });
        });
        expect(await screen.findByText('Tags Listing Page')).toBeInTheDocument();
    });

    it('does not send the update until tag deactivation is confirmed', async () => {
        // Arrange
        let updateRequests = 0;

        server.use(
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json({
                    id: 'tag-1',
                    name: 'Old Tag',
                    isActive: true,
                });
            }),
            http.put(`${API_BASE_URL}/tags/tag-1`, async () => {
                updateRequests += 1;

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagDetailsRoutes([buildEditTagPath('tag-1')]);
        await screen.findByDisplayValue('Old Tag');
        const [saveButton] = screen.getAllByRole('button');
        await user.click(screen.getByRole('checkbox'));
        await user.click(saveButton);

        // Assert
        const dialog = await screen.findByRole('dialog');
        expect(within(dialog).getByText(/элементов каталога\./i)).toBeInTheDocument();
        expect(updateRequests).toBe(0);

        // Act
        const dialogButtons = within(dialog).getAllByRole('button');
        await user.click(dialogButtons[1]);

        // Assert
        await waitFor(() => {
            expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
        });
        expect(updateRequests).toBe(0);
    });

    it('sends the update after deactivation confirmation and keeps save-and-close flow', async () => {
        // Arrange
        let updatedPayload: TagModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json({
                    id: 'tag-1',
                    name: 'Old Tag',
                    isActive: true,
                });
            }),
            http.put(`${API_BASE_URL}/tags/tag-1`, async ({request}) => {
                updatedPayload = (await request.json()) as TagModification;

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagDetailsRoutes([buildEditTagPath('tag-1')]);
        await screen.findByDisplayValue('Old Tag');
        const actionButtons = screen.getAllByRole('button');
        await user.click(screen.getByRole('checkbox'));
        await user.click(actionButtons[1]);
        const dialog = await screen.findByRole('dialog');
        const dialogButtons = within(dialog).getAllByRole('button');
        await user.click(dialogButtons[2]);

        // Assert
        await waitFor(() => {
            expect(updatedPayload).toEqual({
                name: 'Old Tag',
                isActive: false,
            });
        });
        expect(await screen.findByText('Tags Listing Page')).toBeInTheDocument();
    });

    it('shows a warning for breadcrumb navigation when the form has unsaved changes', async () => {
        // Arrange
        const user = userEvent.setup();

        renderTagDetailsRoutes([newTagPath]);

        // Act
        await user.type(screen.getByLabelText('Название'), 'Draft Tag');
        await user.click(screen.getByRole('link', {name: 'Теги'}));

        // Assert
        expect(await screen.findByRole('dialog')).toBeInTheDocument();
        expect(screen.getByText('Есть несохранённые изменения')).toBeInTheDocument();
        expect(screen.queryByText('Tags Listing Page')).not.toBeInTheDocument();

        // Act
        await user.click(screen.getByRole('button', {name: 'Закрыть без сохранения'}));

        // Assert
        expect(await screen.findByText('Tags Listing Page')).toBeInTheDocument();
    });
});
