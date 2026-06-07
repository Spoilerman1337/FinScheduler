import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import {describe, expect, it} from 'vitest';
import {Route, Routes} from 'react-router-dom';
import type {TagModification} from '../../api/tags.types.ts';
import {API_BASE_URL} from '../../config/api.ts';
import {renderWithProviders} from '../../test/render.tsx';
import {server} from '../../test/msw/server.ts';
import {buildEditTagPath, newTagPath, tagEditRoutePath, tagsListPath} from '../routes.ts';
import TagDetailsPage from './TagDetailsPage.tsx';

function renderTagDetailsRoutes(initialEntries: string[]) {
    return renderWithProviders(
        <Routes>
            <Route path={tagsListPath} element={<div>Tags Listing Page</div>} />
            <Route path={newTagPath} element={<TagDetailsPage mode="create" />} />
            <Route path={tagEditRoutePath} element={<TagDetailsPage mode="edit" />} />
        </Routes>,
        {initialEntries},
    );
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
});
