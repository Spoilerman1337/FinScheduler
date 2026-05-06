import {screen, waitFor} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {http, HttpResponse} from "msw";
import {describe, expect, it} from "vitest";
import type {ItemDto, ItemModification} from "../../api/types.ts";
import {renderWithProviders} from "../../test/render.tsx";
import {server} from "../../test/msw/server.ts";
import Items from "./Items.tsx";

const API_BASE_URL = "http://localhost:8081/api";

function buildItem(overrides: Partial<ItemDto> = {}): ItemDto {
    return {
        id: "item-1",
        name: "Coffee",
        price: 199.5,
        description: "Morning drink",
        isActive: true,
        cashback: 5,
        category: "FoodDrinks",
        createdAt: "2025-01-01T00:00:00.000Z",
        tags: [],
        ...overrides,
    };
}

describe("Items integration", () => {
    it("loads active items on mount and renders the first page", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                requests.push(new URL(request.url));

                return HttpResponse.json({
                    data: [buildItem({name: "Utility Bill"})],
                    count: 1,
                });
            }),
        );

        // Act
        renderWithProviders(<Items/>);

        // Assert
        expect(await screen.findByText("Utility Bill")).toBeInTheDocument();
        await waitFor(() => expect(requests).toHaveLength(1));
        expect(requests[0].searchParams.get("page")).toBe("0");
        expect(requests[0].searchParams.get("pageSize")).toBe("10");
        expect(requests[0].searchParams.get("isActive")).toBe("true");
    });

    it("applies search and status filters, then resets them", async () => {
        // Arrange
        const requests: URL[] = [];

        server.use(
            http.get(`${API_BASE_URL}/items`, ({request}) => {
                const url = new URL(request.url);
                const name = url.searchParams.get("name");
                const isActive = url.searchParams.get("isActive");

                requests.push(url);

                if (name === "Tea") {
                    return HttpResponse.json({
                        data: [buildItem({id: "item-2", name: "Tea"})],
                        count: 1,
                    });
                }

                if (isActive === "false") {
                    return HttpResponse.json({
                        data: [buildItem({id: "item-3", name: "Rent", isActive: false})],
                        count: 1,
                    });
                }

                if (isActive === "true") {
                    return HttpResponse.json({
                        data: [buildItem({name: "Coffee"})],
                        count: 1,
                    });
                }

                return HttpResponse.json({
                    data: [buildItem({id: "item-4", name: "All Results"})],
                    count: 1,
                });
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items/>);
        await screen.findByText("Coffee");
        await user.click(screen.getByRole("button", {name: "Неактивные"}));

        // Assert
        expect(await screen.findByText("Rent")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("isActive")).toBe("false");
        });

        // Act
        await user.type(screen.getByPlaceholderText("Поиск по названию..."), "Tea");

        // Assert
        expect(await screen.findByText("Tea")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("name")).toBe("Tea");
            expect(lastRequest?.searchParams.get("isActive")).toBe("false");
        });

        // Act
        await user.click(screen.getByRole("button", {name: "Сброс"}));

        // Assert
        expect(await screen.findByText("All Results")).toBeInTheDocument();
        await waitFor(() => {
            const lastRequest = requests.at(-1);

            expect(lastRequest?.searchParams.get("name")).toBeNull();
            expect(lastRequest?.searchParams.get("isActive")).toBeNull();
        });
    });

    it("edits an existing item and reloads the table", async () => {
        // Arrange
        const currentItems = [buildItem({name: "Old Item"})];
        let updatedPayload: ItemModification | null = null;

        server.use(
            http.get(`${API_BASE_URL}/items`, () => {
                return HttpResponse.json({
                    data: currentItems,
                    count: currentItems.length,
                });
            }),
            http.put(`${API_BASE_URL}/items/:id`, async ({params, request}) => {
                updatedPayload = await request.json() as ItemModification;
                currentItems[0] = {
                    ...currentItems[0],
                    id: String(params.id),
                    name: updatedPayload.name,
                    description: updatedPayload.description,
                    price: updatedPayload.price,
                    cashback: updatedPayload.cashback,
                    isActive: updatedPayload.isActive,
                    category: updatedPayload.category,
                };

                return new HttpResponse(null, {status: 200});
            }),
        );

        const user = userEvent.setup();

        // Act
        renderWithProviders(<Items/>);
        await user.click(await screen.findByText("Old Item"));
        await user.clear(screen.getByLabelText("Название"));
        await user.type(screen.getByLabelText("Название"), "Updated Item");
        await user.click(screen.getByRole("button", {name: "Сохранить"}));

        // Assert
        await waitFor(() => {
            expect(updatedPayload).toEqual({
                name: "Updated Item",
                description: "Morning drink",
                price: 199.5,
                cashback: 5,
                isActive: true,
                category: "FoodDrinks",
                tagIds: [],
            });
        });
        expect(await screen.findByText("Updated Item")).toBeInTheDocument();
        expect(screen.queryByText("Old Item")).not.toBeInTheDocument();
    });
});
