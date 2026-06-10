import {StrictMode} from 'react';
import {createRoot} from 'react-dom/client';
import {Theme} from '@chakra-ui/react';
import './index.css';
import {Provider} from './components/ui/provider.tsx';
import {initializeFaroSdk} from './observability/faro.ts';
import {RouterProvider} from 'react-router-dom';
import {appRouter} from './router.tsx';

initializeFaroSdk();

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Provider>
            <Theme appearance="dark">
                <RouterProvider router={appRouter} />
            </Theme>
        </Provider>
    </StrictMode>,
);
