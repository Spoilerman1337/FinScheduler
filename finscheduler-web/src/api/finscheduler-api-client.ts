import {API_BASE_URL} from "../config/api.ts";

export class FinschedulerApiClient {
    protected baseUrl: string;

    constructor() {
        this.baseUrl = API_BASE_URL;
    }

    protected buildQueryString(params: Record<string, any>): string {
        const searchParams = new URLSearchParams();
        
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                if (Array.isArray(value)) {
                    value.forEach(v => searchParams.append(key, String(v)));
                } else {
                    searchParams.append(key, String(value));
                }
            }
        });

        const queryString = searchParams.toString();
        return queryString ? `?${queryString}` : '';
    }
}

