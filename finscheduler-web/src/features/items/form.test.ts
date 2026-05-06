import {describe, expect, it} from "vitest";
import {
    buildItemModification,
    createDefaultItemFormData,
    mapItemToFormData,
    validateItemFormData,
} from "./form.ts";

describe("items form", () => {
    it("creates the default item form data", () => {
        // Arrange
        const expectedFormData = {
            name: "",
            description: "",
            price: "",
            cashback: "",
            isActive: true,
            category: "",
            tagIds: [],
        };

        // Act
        const formData = createDefaultItemFormData();

        // Assert
        expect(formData).toEqual(expectedFormData);
    });

    it("maps an item dto into edit form state", () => {
        // Arrange
        const item = {
            name: "Coffee",
            description: "Morning coffee",
            price: 199.5,
            cashback: 5,
            isActive: false,
            category: "FoodDrinks",
            tags: [
                {value: "tag-1", label: "Food"},
                {value: "tag-2"},
                {value: ""},
            ],
        };

        // Act
        const formData = mapItemToFormData(item);

        // Assert
        expect(formData).toEqual({
            name: "Coffee",
            description: "Morning coffee",
            price: "199.5",
            cashback: "5",
            isActive: false,
            category: "FoodDrinks",
            tagIds: ["tag-1", "tag-2"],
        });
    });

    it("returns an error when the name is blank", () => {
        // Arrange
        const formData = {
            ...createDefaultItemFormData(),
            category: "FoodDrinks",
            name: "   ",
        };

        // Act
        const error = validateItemFormData(formData);

        // Assert
        expect(error).toBe("Название обязательно для заполнения");
    });

    it("returns an error when the price is not numeric", () => {
        // Arrange
        const formData = {
            ...createDefaultItemFormData(),
            name: "Coffee",
            category: "FoodDrinks",
            price: "abc",
        };

        // Act
        const error = validateItemFormData(formData);

        // Assert
        expect(error).toBe("Цена должна быть числом");
    });

    it("returns an error when cashback is out of range", () => {
        // Arrange
        const formData = {
            ...createDefaultItemFormData(),
            name: "Coffee",
            category: "FoodDrinks",
            cashback: "101",
        };

        // Act
        const error = validateItemFormData(formData);

        // Assert
        expect(error).toBe("Кэшбэк должен быть числом от 0 до 100");
    });

    it("returns an error when the category is missing or unknown", () => {
        // Arrange
        const formData = {
            ...createDefaultItemFormData(),
            name: "Coffee",
            category: "UnknownCategory",
        };

        // Act
        const error = validateItemFormData(formData);

        // Assert
        expect(error).toBe("Выберите категорию");
    });

    it("builds a normalized item payload from valid form data", () => {
        // Arrange
        const formData = {
            name: " Coffee ",
            description: " Morning coffee ",
            price: "",
            cashback: "3",
            isActive: true,
            category: "FoodDrinks",
            tagIds: ["tag-1", "tag-2"],
        };

        // Act
        const payload = buildItemModification(formData);

        // Assert
        expect(payload).toEqual({
            name: "Coffee",
            description: "Morning coffee",
            price: 0,
            cashback: 3,
            isActive: true,
            category: "FoodDrinks",
            tagIds: ["tag-1", "tag-2"],
        });
    });
});
