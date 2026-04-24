import {Flex} from "@chakra-ui/react";
import {motion} from "framer-motion";
import {Route, Routes, matchPath, useLocation} from "react-router-dom";
import Layout from "./subcomponents/Layout.tsx";
import {routes} from "./routes.tsx";
import {Toaster} from "../../components/ui/toaster.tsx";

const MotionFlex = motion.create(Flex)

export default function Main() {
    const location = useLocation();
    const currentRoute = routes.find((route) => matchPath({path: route.path, end: true}, location.pathname));

    return (
        <>
            <MotionFlex layout
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
                        className="custom-scrollbar">
                <Layout headerName={currentRoute?.title}>
                    <Routes>
                        {routes.map((route) => (
                            <Route key={route.path} path={route.path} element={route.element}/>
                        ))}
                    </Routes>
                </Layout>
            </MotionFlex>
            <Toaster />
        </>
    )
}
