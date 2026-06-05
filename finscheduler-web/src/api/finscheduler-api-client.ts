import {API_BASE_URL} from '../config/api.ts';

type QueryParamPrimitive = string | number | boolean;
type QueryParamValue = QueryParamPrimitive | null | undefined | readonly QueryParamPrimitive[];

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
}
