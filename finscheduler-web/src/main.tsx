import {StrictMode} from 'react';
import {createRoot} from 'react-dom/client';
import {Theme} from '@chakra-ui/react';
import './index.css';
import {Provider} from './components/ui/provider.tsx';
import {initializeFaroSdk} from './observability/faro.ts';
import {RouterProvider} from 'react-router-dom';
import {appRouter} from './router.tsx';
import {QueryClientProvider} from '@tanstack/react-query';
import {appQueryClient} from './query/query-client.ts';

initializeFaroSdk();

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={appQueryClient}>
            <Provider>
                <Theme appearance="dark">
                    <RouterProvider router={appRouter} />
                </Theme>
            </Provider>
        </QueryClientProvider>
    </StrictMode>,
);
