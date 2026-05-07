import "@testing-library/jest-dom/vitest";
import {afterAll, afterEach, beforeAll} from "vitest";
import {server} from "./msw/server.ts";

Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: (query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: () => {
        },
        removeListener: () => {
        },
        addEventListener: () => {
        },
        removeEventListener: () => {
        },
        dispatchEvent: () => false,
    }),
});

class ResizeObserverMock {
    observe() {
    }

    unobserve() {
    }

    disconnect() {
    }
}

Object.defineProperty(window, "ResizeObserver", {
    writable: true,
    value: ResizeObserverMock,
});

Object.defineProperty(HTMLElement.prototype, "scrollTo", {
    writable: true,
    value: () => {
    },
});

beforeAll(() => {
    server.listen({onUnhandledRequest: "error"});
});

afterEach(() => {
    server.resetHandlers();
});

afterAll(() => {
    server.close();
});
