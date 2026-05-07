import {describe, expect, it} from "vitest";
import {buildTagFilter} from "./tags.ts";

describe("tags api", () => {
    it("buildTagFilter maps page, search, and active status to TagFilter", () => {
        // Arrange
        const params = {
            page: 4,
            pageSize: 25,
            searchTerm: "groceries",
            statusFilter: "Active" as const,
        };

        // Act
        const filter = buildTagFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 3,
            pageSize: 25,
            name: "groceries",
            isActive: true,
        });
    });

    it("buildTagFilter omits optional criteria when they are inactive", () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: "",
            statusFilter: "All" as const,
        };

        // Act
        const filter = buildTagFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
        });
    });
});
