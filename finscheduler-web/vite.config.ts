import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import * as path from 'node:path';

export default defineConfig({
    plugins: [react()],
    build: {
        rollupOptions: {
            output: {
                manualChunks(id) {
                    const normalizedId = id.replace(/\\/g, '/');

                    if (!normalizedId.includes('/node_modules/')) {
                        return undefined;
                    }

                    if (normalizedId.includes('/@grafana/')) {
                        return 'observability-vendor';
                    }

                    if (normalizedId.includes('/@chakra-ui/react/dist/esm/components/')) {
                        return 'chakra-components-vendor';
                    }

                    if (
                        normalizedId.includes('/@chakra-ui/react/dist/esm/hooks/') ||
                        normalizedId.includes('/@chakra-ui/react/dist/esm/styled-system/') ||
                        normalizedId.includes('/@chakra-ui/react/dist/esm/theme/') ||
                        normalizedId.includes('/@chakra-ui/react/dist/esm/utils/')
                    ) {
                        return 'chakra-core-vendor';
                    }

                    if (normalizedId.includes('/@chakra-ui/')) {
                        return 'chakra-vendor';
                    }

                    if (normalizedId.includes('/@emotion/')) {
                        return 'emotion-vendor';
                    }

                    if (normalizedId.includes('/@ark-ui/') || normalizedId.includes('/@zag-js/')) {
                        return 'ui-state-vendor';
                    }

                    if (normalizedId.includes('/next-themes/')) {
                        return 'theme-vendor';
                    }

                    if (normalizedId.includes('/framer-motion/')) {
                        return 'motion-vendor';
                    }

                    if (
                        normalizedId.includes('/react-icons/') ||
                        normalizedId.includes('/lucide-react/')
                    ) {
                        return 'icon-vendor';
                    }

                    if (normalizedId.includes('/react-router-dom/')) {
                        return 'router-vendor';
                    }

                    if (normalizedId.includes('/react-dom/') || normalizedId.includes('/react/')) {
                        return 'react-vendor';
                    }

                    return 'vendor';
                },
            },
        },
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src'),
        },
    },
});
