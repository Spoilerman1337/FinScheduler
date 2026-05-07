import {describe, expect, it} from "vitest";
import {getCashbackColor} from "./types.ts";

describe("items feature types", () => {
    it("getCashbackColor returns the expected color buckets", () => {
        // Arrange
        const cashbackValues = {
            missing: undefined,
            zero: 0,
            low: 1,
            high: 2,
        };

        // Act
        const colors = {
            missing: getCashbackColor(cashbackValues.missing),
            zero: getCashbackColor(cashbackValues.zero),
            low: getCashbackColor(cashbackValues.low),
            high: getCashbackColor(cashbackValues.high),
        };

        // Assert
        expect(colors).toEqual({
            missing: "neon.pink",
            zero: "neon.pink",
            low: "neon.yellow",
            high: "neon.green",
        });
    });
});
