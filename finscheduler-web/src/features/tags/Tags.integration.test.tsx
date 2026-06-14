import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {http, HttpResponse} from 'msw';
import type {RouteObject} from 'react-router-dom';
import {describe, expect, it, vi} from 'vitest';
import {API_BASE_URL} from '../../config/api.ts';
import type {TagDetailedDto, TagListingDto} from '../../api/tags.types.ts';
import {renderWithProviders} from '../../test/render.tsx';
import {renderWithDataRouter} from '../../test/renderDataRouter.tsx';
import {server} from '../../test/msw/server.ts';
import {buildEditTagPath, newTagPath, tagEditRoutePath, tagsListPath} from '../routes.ts';
import Tags from './Tags.tsx';
import TagDetailsPage from './TagDetailsPage.tsx';

type TagTestData = TagListingDto & TagDetailedDto;

function buildTag(overrides: Partial<TagTestData> = {}): TagTestData {
    return {
        id: 'tag-1',
        name: 'Food',
        isActive: true,
        ...overrides,
    };
}

function renderTagsRoutes(initialEntries: string[] = [tagsListPath]) {
    const routes: RouteObject[] = [
        {path: tagsListPath, element: <Tags />},
        {path: newTagPath, element: <TagDetailsPage mode="create" />},
        {path: tagEditRoutePath, element: <TagDetailsPage mode="edit" />},
    ];

    return renderWithDataRouter(routes, {initialEntries});
}

describe('Tags integration', () => {
    it('loads active tags on mount and renders the first page', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                requests.push(new URL(request.url));

                return HttpResponse.json({
                    data: [buildTag({name: 'Groceries'})],
                    count: 1,
                });
            }),
        );

        // Act
        renderWithProviders(<Tags />);

        // Assert
        expect(await screen.findByText('Groceries')).toBeInTheDocument();
        await waitFor(() => expect(requests).toHaveLength(1));
        expect(requests[0].searchParams.get('page')).toBe('0');
        expect(requests[0].searchParams.get('pageSize')).toBe('10');
        expect(requests[0].searchParams.get('isActive')).toBe('true');
    });

    it('applies search and status filters, then resets them', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const name = url.searchParams.get('name');
                const isActive = url.searchParams.get('isActive');

                requests.push(url);

                if (name === 'Travel') {
                    return HttpResponse.json({
                        data: [buildTag({id: 'tag-2', name: 'Travel'})],
                        count: 1,
                    });
                }

                if (isActive === 'false') {
                    return HttpResponse.json({
                        data: [buildTag({id: 'tag-3', name: 'Archive', isActive: false})],
                        count: 1,
                    });
                }

                if (isActive === 'true') {
                    return HttpResponse.json({
                        data: [buildTag({name: 'Food'})],
                        count: 1,
                    });
                }

                return HttpResponse.json({
                    data: [buildTag({id: 'tag-4', name: 'All Tags'})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags />);
        await screen.findByText('Food');
        await user.click(screen.getByRole('button', {name: 'Неактивные'}));

        // Assert
        expect(await screen.findByText('Archive')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('isActive')).toBe('false');
        });

        // Act
        await user.type(screen.getByPlaceholderText('Поиск по названию...'), 'Travel');

        // Assert
        expect(await screen.findByText('Travel')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('name')).toBe('Travel');
            expect(lastRequest?.searchParams.get('isActive')).toBe('false');
        });

        // Act
        await user.click(screen.getByRole('button', {name: 'Сброс'}));

        // Assert
        expect(await screen.findByText('All Tags')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('name')).toBeNull();
            expect(lastRequest?.searchParams.get('isActive')).toBeNull();
        });
    });

    it('navigates to the create page from the add button', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: 'Existing Tag'})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagsRoutes();
        await screen.findByText('Existing Tag');
        await user.click(screen.getByRole('button', {name: 'Добавить'}));

        // Assert
        expect(await screen.findByText('Новый тег')).toBeInTheDocument();
        expect(screen.getByText('Создание')).toBeInTheDocument();
    });

    it('navigates to the edit page after a row click', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: 'Old Tag'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json(buildTag({name: 'Old Tag'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagsRoutes();
        await user.click(await screen.findByText('Old Tag'));

        // Assert
        expect(await screen.findByText('Редактирование тега')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('Old Tag');
    });

    it('returns to the list after going back from a clean edit page', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: 'Food'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json(buildTag({name: 'Food'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagsRoutes();
        await user.click(await screen.findByText('Food'));
        expect(await screen.findByText('Редактирование тега')).toBeInTheDocument();
        await user.click(screen.getByRole('button', {name: 'Назад'}));

        // Assert
        expect(await screen.findByText('Food')).toBeInTheDocument();
        expect(screen.queryByText('Редактирование тега')).not.toBeInTheDocument();
    });

    it('hides the add button when a tag is selected', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [
                        buildTag({id: 'tag-1', name: 'Food'}),
                        buildTag({id: 'tag-2', name: 'Travel'}),
                    ],
                    count: 2,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags />);
        await screen.findByText('Food');
        const checkboxes = await screen.findAllByRole('checkbox');

        await user.click(checkboxes[1]);

        // Assert
        expect(screen.queryByRole('button', {name: 'Добавить'})).not.toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Редактировать'})).toBeInTheDocument();
    });

    it('changes the page when the paginator is used', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get('page');

                requests.push(url);

                return HttpResponse.json({
                    data:
                        page === '1'
                            ? [buildTag({id: 'tag-2', name: 'Second Page Tag'})]
                            : [buildTag({name: 'First Page Tag'})],
                    count: 11,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags />);
        await screen.findByText('First Page Tag');
        await user.click(screen.getByRole('button', {name: 'last page, page 2'}));

        // Assert
        expect(await screen.findByText('Second Page Tag')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('page')).toBe('1');
        });
    });

    it('changes the page size and reloads the first page', async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get('page');
                const pageSize = url.searchParams.get('pageSize');

                requests.push(url);

                if (pageSize === '25') {
                    return HttpResponse.json({
                        data: [buildTag({id: 'tag-3', name: 'Twenty Five Per Page'})],
                        count: 30,
                    });
                }

                if (page === '1') {
                    return HttpResponse.json({
                        data: [buildTag({id: 'tag-2', name: 'Second Page Tag'})],
                        count: 30,
                    });
                }

                return HttpResponse.json({
                    data: [buildTag({name: 'First Page Tag'})],
                    count: 30,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags />);
        await screen.findByText('First Page Tag');
        await user.click(screen.getByRole('button', {name: /page 2/i}));
        await screen.findByText('Second Page Tag');
        await user.click(screen.getByRole('combobox'));
        await user.click(await screen.findByRole('option', {name: '25'}));

        // Assert
        expect(await screen.findByText('Twenty Five Per Page')).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get('page')).toBe('0');
            expect(lastRequest?.searchParams.get('pageSize')).toBe('25');
        });
    });

    it('shows an empty state when no tags are returned', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [],
                    count: 0,
                });
            }),
        );

        // Act
        renderWithProviders(<Tags />);

        // Assert
        expect(await screen.findByText('Данные не найдены.')).toBeInTheDocument();
    });

    it('shows an error state when tags request fails', async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json(
                    {message: 'boom'},
                    {status: 500, statusText: 'Internal Server Error'},
                );
            }),
        );

        try {
            // Act
            renderWithProviders(<Tags />);

            // Assert
            expect(
                await screen.findByText('Failed to fetch tags: Internal Server Error'),
            ).toBeInTheDocument();
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });

    it('shows a save error when creating a tag fails', async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: 'Existing Tag'})],
                    count: 1,
                });
            }),
            http.post(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json(
                    {message: 'boom'},
                    {status: 500, statusText: 'Internal Server Error'},
                );
            }),
        );

        const user = userEvent.setup();

        try {
            // Act
            renderTagsRoutes();
            await screen.findByText('Existing Tag');
            await user.click(screen.getByRole('button', {name: 'Добавить'}));
            await user.type(screen.getByLabelText('Название'), 'Broken Tag');
            await user.click(screen.getByRole('button', {name: 'Сохранить'}));

            // Assert
            expect(
                await screen.findByText('Failed to create tag: Internal Server Error'),
            ).toBeInTheDocument();
            expect(screen.getByLabelText('Название')).toHaveValue('Broken Tag');
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });

    it('opens the edit page from the selected-row action button', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: 'Coffee'})],
                    count: 1,
                });
            }),
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json(buildTag({name: 'Coffee'}));
            }),
        );

        const user = userEvent.setup();

        // Act
        renderTagsRoutes();
        await screen.findByText('Coffee');
        const checkboxes = await screen.findAllByRole('checkbox');

        await user.click(checkboxes[1]);
        await user.click(screen.getByRole('button', {name: 'Редактировать'}));

        // Assert
        expect(await screen.findByText('Редактирование тега')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('Coffee');
    });

    it('keeps the list route working when opening the edit page directly', async () => {
        // Arrange
        server.use(
            http.get(`${API_BASE_URL}/tags/tag-1`, () => {
                return HttpResponse.json(buildTag({name: 'Travel'}));
            }),
        );

        // Act
        renderTagsRoutes([buildEditTagPath('tag-1')]);

        // Assert
        expect(await screen.findByText('Редактирование тега')).toBeInTheDocument();
        expect(screen.getByLabelText('Название')).toHaveValue('Travel');
    });
});
