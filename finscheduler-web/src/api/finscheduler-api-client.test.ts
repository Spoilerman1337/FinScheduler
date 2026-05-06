import {describe, expect, it} from "vitest";
import {FinschedulerApiClient} from "./finscheduler-api-client.ts";

class TestApiClient extends FinschedulerApiClient {
    public build(params: Record<string, unknown>) {
        return this.buildQueryString(params);
    }
}

describe("FinschedulerApiClient.buildQueryString", () => {
    const client = new TestApiClient();

    it("returns an empty string when there are no serializable params", () => {
        expect(client.build({})).toBe("");
        expect(client.build({name: undefined, page: null})).toBe("");
    });

    it("serializes scalar values without dropping 0 or false", () => {
        const query = client.build({
            page: 0,
            pageSize: 25,
            isActive: false,
            name: "coffee",
        });
        const params = new URLSearchParams(query.slice(1));

        expect(query.startsWith("?")).toBe(true);
        expect(params.get("page")).toBe("0");
        expect(params.get("pageSize")).toBe("25");
        expect(params.get("isActive")).toBe("false");
        expect(params.get("name")).toBe("coffee");
    });

    it("skips null and undefined values", () => {
        const query = client.build({
            page: 1,
            pageSize: undefined,
            isActive: null,
            name: "tea",
        });
        const params = new URLSearchParams(query.slice(1));

        expect(params.get("page")).toBe("1");
        expect(params.has("pageSize")).toBe(false);
        expect(params.has("isActive")).toBe(false);
        expect(params.get("name")).toBe("tea");
    });

    it("serializes arrays as repeated query params", () => {
        const query = client.build({
            ids: ["1", "2", "3"],
            categories: ["FoodDrinks", "Transport"],
        });
        const params = new URLSearchParams(query.slice(1));

        expect(params.getAll("ids")).toEqual(["1", "2", "3"]);
        expect(params.getAll("categories")).toEqual(["FoodDrinks", "Transport"]);
    });
});
