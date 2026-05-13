import {Provider} from "./components/ui/provider.tsx";
import {Theme, Flex} from "@chakra-ui/react";
import {motion} from "framer-motion";
import Sidebar from "./layout/sidebar/Sidebar.tsx";
import Main from "./layout/main/Main.tsx";

const MotionFlex = motion.create(Flex)

export default function App() {
    return (<Provider>
        <Theme appearance="dark">
            <MotionFlex layout w="100vw" h="100vh" bg={"background"} minHeight={0}>
                <Sidebar/>
                <Main/>
            </MotionFlex>
        </Theme>
    </Provider>);
}
