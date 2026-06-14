import {Theme} from '@chakra-ui/react';
import {QueryClientProvider} from '@tanstack/react-query';
import {render} from '@testing-library/react';
import type {ReactElement} from 'react';
import {MemoryRouter} from 'react-router-dom';
import {Provider} from '../components/ui/provider.tsx';
import {createAppQueryClient} from '../query/query-client.ts';

interface RenderWithProvidersOptions {
    route?: string;
    initialEntries?: string[];
}

export function renderWithProviders(ui: ReactElement, options: RenderWithProvidersOptions = {}) {
    const initialEntries = options.initialEntries ?? [options.route ?? '/'];
    const queryClient = createAppQueryClient();

    return render(
        <QueryClientProvider client={queryClient}>
            <Provider>
                <Theme>
                    <MemoryRouter initialEntries={initialEntries}>{ui}</MemoryRouter>
                </Theme>
            </Provider>
        </QueryClientProvider>,
    );
}
