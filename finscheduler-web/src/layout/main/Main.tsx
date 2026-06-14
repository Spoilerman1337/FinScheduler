import {Flex} from '@chakra-ui/react';
import {motion} from 'framer-motion';
import {Outlet} from 'react-router-dom';
import {useCurrentRouteTitle} from '../../hooks/useCurrentRouteTitle.ts';
import Layout from './subcomponents/Layout.tsx';
import {Toaster} from '../../components/ui/toaster.tsx';

const MotionFlex = motion.create(Flex);

export default function Main() {
    const currentRouteTitle = useCurrentRouteTitle();

    return (
        <>
            <MotionFlex
                layout
                flex={1}
                h="100%"
                overflowY="auto"
                overflowX="hidden"
                pl="4"
                pr="4"
                pt="4"
                pb="0"
                direction="column"
                backdropFilter="blur(8px)"
                bg="rgba(0,0,0,0.25)"
                minHeight={0}
                className="custom-scrollbar"
            >
                <Layout headerName={currentRouteTitle}>
                    <Outlet />
                </Layout>
            </MotionFlex>
            <Toaster />
        </>
    );
}
