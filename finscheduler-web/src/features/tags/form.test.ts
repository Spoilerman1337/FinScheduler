import {describe, expect, it} from "vitest";
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    validateTagFormData,
} from "./form.ts";

describe("tags form", () => {
    it("creates the default tag form data", () => {
        // Arrange
        const expectedFormData = {
            name: "",
            isActive: true,
        };

        // Act
        const formData = createDefaultTagFormData();

        // Assert
        expect(formData).toEqual(expectedFormData);
    });

    it("maps a tag dto into edit form state", () => {
        // Arrange
        const item = {
            name: "Transport",
            isActive: false,
        };

        // Act
        const formData = mapTagToFormData(item);

        // Assert
        expect(formData).toEqual({
            name: "Transport",
            isActive: false,
        });
    });

    it("returns an error when the name is blank", () => {
        // Arrange
        const formData = {
            ...createDefaultTagFormData(),
            name: "   ",
        };

        // Act
        const error = validateTagFormData(formData);

        // Assert
        expect(error).toBe("Название обязательно для заполнения");
    });

    it("builds a normalized tag modification from valid form data", () => {
        // Arrange
        const formData = {
            name: " Transport ",
            isActive: false,
        };

        // Act
        const payload = buildTagModification(formData);

        // Assert
        expect(payload).toEqual({
            name: "Transport",
            isActive: false,
        });
    });
});
