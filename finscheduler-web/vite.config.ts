import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import * as path from 'node:path';

const nodeModulesSegment = String.raw`[\\/]node_modules[\\/]`;

export default defineConfig({
    plugins: [react()],
    build: {
        rolldownOptions: {
            output: {
                codeSplitting: {
                    groups: [
                        {
                            name: 'observability-vendor',
                            test: new RegExp(`${nodeModulesSegment}@grafana[\\\\/]`),
                            priority: 110,
                        },
                        {
                            name: 'chakra-components-vendor',
                            test: new RegExp(
                                `${nodeModulesSegment}@chakra-ui[\\\\/]react[\\\\/]dist[\\\\/]esm[\\\\/]components[\\\\/]`
                            ),
                            priority: 100,
                        },
                        {
                            name: 'chakra-core-vendor',
                            test: new RegExp(
                                `${nodeModulesSegment}@chakra-ui[\\\\/]react[\\\\/]dist[\\\\/]esm[\\\\/](hooks|styled-system|theme|utils)[\\\\/]`
                            ),
                            priority: 90,
                        },
                        {
                            name: 'chakra-vendor',
                            test: new RegExp(`${nodeModulesSegment}@chakra-ui[\\\\/]`),
                            priority: 80,
                        },
                        {
                            name: 'emotion-vendor',
                            test: new RegExp(`${nodeModulesSegment}@emotion[\\\\/]`),
                            priority: 70,
                        },
                        {
                            name: 'ui-state-vendor',
                            test: new RegExp(`${nodeModulesSegment}(@ark-ui|@zag-js)[\\\\/]`),
                            priority: 60,
                        },
                        {
                            name: 'theme-vendor',
                            test: new RegExp(`${nodeModulesSegment}next-themes[\\\\/]`),
                            priority: 50,
                        },
                        {
                            name: 'motion-vendor',
                            test: new RegExp(`${nodeModulesSegment}framer-motion[\\\\/]`),
                            priority: 40,
                        },
                        {
                            name: 'icon-vendor',
                            test: new RegExp(`${nodeModulesSegment}(react-icons|lucide-react)[\\\\/]`),
                            priority: 30,
                        },
                        {
                            name: 'router-vendor',
                            test: new RegExp(`${nodeModulesSegment}react-router-dom[\\\\/]`),
                            priority: 20,
                        },
                        {
                            name: 'react-vendor',
                            test: new RegExp(`${nodeModulesSegment}(react-dom|react)[\\\\/]`),
                            priority: 10,
                        },
                        {
                            name: 'vendor',
                            test: new RegExp(nodeModulesSegment),
                        },
                    ],
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
