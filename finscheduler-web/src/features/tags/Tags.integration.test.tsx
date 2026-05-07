import {screen, waitFor} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {http, HttpResponse} from "msw";
import {describe, expect, it, vi} from "vitest";
import type {TagDto, TagModification} from "../../api/types.ts";
import {renderWithProviders} from "../../test/render.tsx";
import {server} from "../../test/msw/server.ts";
import Tags from "./Tags.tsx";

const API_BASE_URL = "http://localhost:8081/api";

function buildTag(overrides: Partial<TagDto> = {}): TagDto {
    return {
        id: "tag-1",
        name: "Food",
        isActive: true,
        ...overrides,
    };
}

describe("Tags integration", () => {
    it("loads active tags on mount and renders the first page", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                requests.push(new URL(request.url));

                return HttpResponse.json({
                    data: [buildTag({name: "Groceries"})],
                    count: 1,
                });
            }),
        );

        // Act
        renderWithProviders(<Tags/>);

        // Assert
        expect(await screen.findByText("Groceries")).toBeInTheDocument();
        await waitFor(() => expect(requests).toHaveLength(1));
        expect(requests[0].searchParams.get("page")).toBe("0");
        expect(requests[0].searchParams.get("pageSize")).toBe("10");
        expect(requests[0].searchParams.get("isActive")).toBe("true");
    });

    it("applies search and status filters, then resets them", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const name = url.searchParams.get("name");
                const isActive = url.searchParams.get("isActive");

                requests.push(url);

                if (name === "Travel") {
                    return HttpResponse.json({
                        data: [buildTag({id: "tag-2", name: "Travel"})],
                        count: 1,
                    });
                }

                if (isActive === "false") {
                    return HttpResponse.json({
                        data: [buildTag({id: "tag-3", name: "Archive", isActive: false})],
                        count: 1,
                    });
                }

                if (isActive === "true") {
                    return HttpResponse.json({
                        data: [buildTag({name: "Food"})],
                        count: 1,
                    });
                }

                return HttpResponse.json({
                    data: [buildTag({id: "tag-4", name: "All Tags"})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags/>);
        await screen.findByText("Food");
        await user.click(screen.getByRole("button", {name: "Неактивные"}));

        // Assert
        expect(await screen.findByText("Archive")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("isActive")).toBe("false");
        });

        // Act
        await user.type(screen.getByPlaceholderText("Поиск по названию..."), "Travel");

        // Assert
        expect(await screen.findByText("Travel")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("name")).toBe("Travel");
            expect(lastRequest?.searchParams.get("isActive")).toBe("false");
        });

        // Act
        await user.click(screen.getByRole("button", {name: "Сброс"}));

        // Assert
        expect(await screen.findByText("All Tags")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("name")).toBeNull();
            expect(lastRequest?.searchParams.get("isActive")).toBeNull();
        });
    });

    it("creates a new tag and reloads the table", async () => {
        // Arrange
        const currentTags = [buildTag({name: "Existing Tag"})];
        let createdPayload: TagModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const isActive = new URL(request.url).searchParams.get("isActive");

                return HttpResponse.json({
                    data: isActive === "true"
                        ? currentTags.filter((tag) => tag.isActive)
                        : currentTags,
                    count: currentTags.length,
                });
            }),
            http.post(`${API_BASE_URL}/tags`, async ({request}) => {
                createdPayload = await request.json() as TagModification;
                currentTags.push(buildTag({
                    id: "tag-2",
                    name: createdPayload.name,
                    isActive: createdPayload.isActive,
                }));

                return HttpResponse.json("tag-2");
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags/>);
        await screen.findByText("Existing Tag");
        await user.click(screen.getByRole("button", {name: "Добавить"}));
        await user.type(screen.getByLabelText("Название"), "New Tag");
        await user.click(screen.getByRole("button", {name: "Сохранить"}));

        // Assert
        await waitFor(() => {
            expect(createdPayload).toEqual({
                name: "New Tag",
                isActive: true,
            });
        });
        expect(await screen.findByText("New Tag")).toBeInTheDocument();
    });

    it("edits an existing tag and reloads the table", async () => {
        // Arrange
        const currentTags = [buildTag({name: "Old Tag"})];
        let updatedPayload: TagModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: currentTags,
                    count: currentTags.length,
                });
            }),
            http.put(`${API_BASE_URL}/tags/:id`, async ({params, request}) => {
                updatedPayload = await request.json() as TagModification;
                currentTags[0] = {
                    ...currentTags[0],
                    id: String(params.id),
                    name: updatedPayload.name,
                    isActive: updatedPayload.isActive,
                };

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags/>);
        await user.click(await screen.findByText("Old Tag"));
        await user.clear(screen.getByLabelText("Название"));
        await user.type(screen.getByLabelText("Название"), "Updated Tag");
        await user.click(screen.getByRole("button", {name: "Сохранить"}));

        // Assert
        await waitFor(() => {
            expect(updatedPayload).toEqual({
                name: "Updated Tag",
                isActive: true,
            });
        });
        expect(await screen.findByText("Updated Tag")).toBeInTheDocument();
        expect(screen.queryByText("Old Tag")).not.toBeInTheDocument();
    });

    it("changes the page when the paginator is used", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get("page");

                requests.push(url);

                return HttpResponse.json({
                    data: page === "1"
                        ? [buildTag({id: "tag-2", name: "Second Page Tag"})]
                        : [buildTag({name: "First Page Tag"})],
                    count: 11,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags/>);
        await screen.findByText("First Page Tag");
        await user.click(screen.getByRole("button", {name: "last page, page 2"}));

        // Assert
        expect(await screen.findByText("Second Page Tag")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("page")).toBe("1");
        });
    });

    it("changes the page size and reloads the first page", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/tags`, ({request}) => {
                const url = new URL(request.url);
                const page = url.searchParams.get("page");
                const pageSize = url.searchParams.get("pageSize");

                requests.push(url);

                if (pageSize === "25") {
                    return HttpResponse.json({
                        data: [buildTag({id: "tag-3", name: "Twenty Five Per Page"})],
                        count: 30,
                    });
                }

                if (page === "1") {
                    return HttpResponse.json({
                        data: [buildTag({id: "tag-2", name: "Second Page Tag"})],
                        count: 30,
                    });
                }

                return HttpResponse.json({
                    data: [buildTag({name: "First Page Tag"})],
                    count: 30,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Tags/>);
        await screen.findByText("First Page Tag");
        await user.click(screen.getByRole("button", {name: /page 2/i}));
        await screen.findByText("Second Page Tag");
        await user.click(screen.getByRole("combobox"));
        await user.click(await screen.findByRole("option", {name: "25"}));

        // Assert
        expect(await screen.findByText("Twenty Five Per Page")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("page")).toBe("0");
            expect(lastRequest?.searchParams.get("pageSize")).toBe("25");
        });
    });

    it("shows an empty state when no tags are returned", async () => {
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
        renderWithProviders(<Tags/>);

        // Assert
        expect(await screen.findByText("Данные не найдены.")).toBeInTheDocument();
    });

    it("shows an error state when tags request fails", async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json(
                    {message: "boom"},
                    {status: 500, statusText: "Internal Server Error"},
                );
            }),
        );

        try {
            // Act
            renderWithProviders(<Tags/>);

            // Assert
            expect(await screen.findByText("Failed to fetch tags: Internal Server Error")).toBeInTheDocument();
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });

    it("shows a save error when creating a tag fails", async () => {
        // Arrange
        const consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => undefined);

        server.use(
            http.get(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json({
                    data: [buildTag({name: "Existing Tag"})],
                    count: 1,
                });
            }),
            http.post(`${API_BASE_URL}/tags`, () => {
                return HttpResponse.json(
                    {message: "boom"},
                    {status: 500, statusText: "Internal Server Error"},
                );
            }),
        );

        const user = userEvent.setup();

        try {
            // Act
            renderWithProviders(<Tags/>);
            await screen.findByText("Existing Tag");
            await user.click(screen.getByRole("button", {name: "Добавить"}));
            await user.type(screen.getByLabelText("Название"), "Broken Tag");
            await user.click(screen.getByRole("button", {name: "Сохранить"}));

            // Assert
            expect(await screen.findByText("Failed to create tag: Internal Server Error")).toBeInTheDocument();
            expect(screen.getByLabelText("Название")).toHaveValue("Broken Tag");
        } finally {
            consoleErrorSpy.mockRestore();
        }
    });
});
