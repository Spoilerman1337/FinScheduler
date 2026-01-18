import {Provider} from "./components/ui/provider.tsx";
import {Theme, Flex} from "@chakra-ui/react";
import {motion} from "framer-motion";
import Sidebar from "./components/features/sideMenu/Sidebar.tsx";
import Main from "./components/features/main/Main.tsx";

const MotionFlex = motion.create(Flex)

export default function App() {
    return (<Provider>
        <Theme>
            <MotionFlex layout w="100vw" h="100vh" bg={"background"} minHeight={0}>
                <Sidebar/>
                <Main/>
            </MotionFlex>
        </Theme>
    </Provider>);
}
