import {API_BASE_URL} from '../config/api.ts';

type QueryParamPrimitive = string | number | boolean;
type QueryParamValue = QueryParamPrimitive | null | undefined | readonly QueryParamPrimitive[];

type NumericRange = {
    from: number | null;
    to: number | null;
};

type DateRange = {
    from: string | null;
    to: string | null;
};

export class FinschedulerApiClient {
    protected baseUrl: string;

    constructor() {
        this.baseUrl = API_BASE_URL;
    }

    protected buildQueryString<T extends object>(params: T): string {
        const searchParams = new URLSearchParams();

        (Object.entries(params) as [string, QueryParamValue][]).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                if (Array.isArray(value)) {
                    value.forEach((v) => searchParams.append(key, String(v)));
                } else {
                    searchParams.append(key, String(value));
                }
            }
        });

        const queryString = searchParams.toString();
        return queryString ? `?${queryString}` : '';
    }

    static buildNonNegativeRange(fromValue: string, toValue: string): NumericRange {
        const from = this.parseNonNegativeNumberValue(fromValue);
        const to = this.parseNonNegativeNumberValue(toValue);

        if (from !== null && to !== null && from > to) {
            return {
                from: to,
                to: from,
            };
        }

        return {from, to};
    }

    static buildDateRange(fromValue?: string, toValue?: string): DateRange {
        const from = this.parseDateInputValue(fromValue ?? '');
        const to = this.parseDateInputValue(toValue ?? '');

        if (from !== null && to !== null && from > to) {
            return {
                from: to,
                to: from,
            };
        }

        return {from, to};
    }

    static toLocalDayBoundaryIso(value: string, endOfDay: boolean): string {
        const time = endOfDay ? '23:59:59.999' : '00:00:00.000';
        return new Date(`${value}T${time}`).toISOString();
    }

    private static parseNonNegativeNumberValue(value: string): number | null {
        if (!value) {
            return null;
        }

        const normalized = value.replace(',', '.').trim();

        if (!normalized) {
            return null;
        }

        const parsed = Number(normalized);

        if (!Number.isFinite(parsed) || parsed < 0) {
            return null;
        }

        return parsed;
    }

    private static parseDateInputValue(value: string): string | null {
        if (!value) {
            return null;
        }

        const normalized = value.trim();
        const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(normalized);

        if (!match) {
            return null;
        }

        const year = Number(match[1]);
        const month = Number(match[2]);
        const day = Number(match[3]);
        const parsed = new Date(Date.UTC(year, month - 1, day));

        const isValidDate =
            parsed.getUTCFullYear() === year &&
            parsed.getUTCMonth() === month - 1 &&
            parsed.getUTCDate() === day;

        return isValidDate ? normalized : null;
    }
}
