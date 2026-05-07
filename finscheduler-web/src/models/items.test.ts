import {describe, expect, it} from "vitest";
import {categoryOptions, categoryTranslations} from "./items.ts";

describe("items model", () => {
    it("creates one value-label option for every category translation in the same order", () => {
        // Arrange
        const expectedOptions = Object.entries(categoryTranslations).map(([value, label]) => ({
            label,
            value,
        }));

        // Act
        const actualOptions = categoryOptions;

        // Assert
        expect(actualOptions).toEqual(expectedOptions);
    });

    it("contains only unique and non-empty option values and labels", () => {
        // Arrange
        const optionValues = categoryOptions.map((option) => option.value);
        const optionLabels = categoryOptions.map((option) => option.label);

        // Act
        const uniqueValuesCount = new Set(optionValues).size;
        const hasOnlyFilledValues = optionValues.every((value) => value.trim().length > 0);
        const hasOnlyFilledLabels = optionLabels.every((label) => label.trim().length > 0);

        // Assert
        expect(uniqueValuesCount).toBe(optionValues.length);
        expect(hasOnlyFilledValues).toBe(true);
        expect(hasOnlyFilledLabels).toBe(true);
    });
});
