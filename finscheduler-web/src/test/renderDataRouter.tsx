import {Theme} from '@chakra-ui/react';
import {render} from '@testing-library/react';
import type {RouteObject} from 'react-router-dom';
import {createMemoryRouter, RouterProvider} from 'react-router-dom';
import {Provider} from '../components/ui/provider.tsx';

interface RenderWithDataRouterOptions {
    initialEntries?: string[];
}

let isRequestCompatibilityInstalled = false;

function installRequestCompatibility() {
    if (isRequestCompatibilityInstalled) {
        return;
    }

    const NativeRequest = globalThis.Request;

    class RequestCompatibility extends NativeRequest {
        constructor(...args: ConstructorParameters<typeof Request>) {
            const [input, init] = args;

            try {
                super(input, init);
            } catch (error) {
                if (
                    init?.signal != null &&
                    error instanceof TypeError &&
                    error.message.includes('Expected signal')
                ) {
                    super(input, {...init, signal: undefined});
                    return;
                }

                throw error;
            }
        }
    }

    globalThis.Request = RequestCompatibility as typeof Request;
    isRequestCompatibilityInstalled = true;
}

export function renderWithDataRouter(
    routes: RouteObject[],
    options: RenderWithDataRouterOptions = {},
) {
    installRequestCompatibility();

    const router = createMemoryRouter(routes, {
        initialEntries: options.initialEntries ?? ['/'],
    });

    return {
        router,
        ...render(
            <Provider>
                <Theme appearance="dark">
                    <RouterProvider router={router} />
                </Theme>
            </Provider>,
        ),
    };
}
