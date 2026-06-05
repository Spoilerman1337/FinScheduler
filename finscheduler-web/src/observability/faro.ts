import {
    createReactRouterV7Options,
    getWebInstrumentations,
    initializeFaro,
    ReactIntegration,
} from '@grafana/faro-react';
import {TracingInstrumentation} from '@grafana/faro-web-tracing';
import {
    createRoutesFromChildren,
    matchRoutes,
    Routes,
    useLocation,
    useNavigationType,
} from 'react-router-dom';

const faroCollectUrl = import.meta.env.VITE_FARO_COLLECT_URL?.trim() ?? '';
const appEnvironment = import.meta.env.VITE_APP_ENVIRONMENT?.trim() ?? import.meta.env.MODE;

export const isFaroEnabled = faroCollectUrl.length > 0;

let isFaroInitialized = false;

function escapeRegExp(value: string) {
    return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function createApiTraceTargets() {
    return [new RegExp(`^${escapeRegExp(window.location.origin)}/api(?:/|$)`)];
}

export function initializeFaroSdk() {
    if (!isFaroEnabled || isFaroInitialized) {
        return;
    }

    initializeFaro({
        url: faroCollectUrl,
        app: {
            name: 'finscheduler-web',
            environment: appEnvironment,
        },
        instrumentations: [
            ...getWebInstrumentations(),
            new ReactIntegration({
                router: createReactRouterV7Options({
                    createRoutesFromChildren,
                    matchRoutes,
                    Routes,
                    useLocation,
                    useNavigationType,
                }),
            }),
            new TracingInstrumentation({
                instrumentationOptions: {
                    propagateTraceHeaderCorsUrls: createApiTraceTargets(),
                },
            }),
        ],
    });

    isFaroInitialized = true;
}
