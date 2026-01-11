import {Flex} from "@chakra-ui/react";
import {motion} from "framer-motion";
import {Route, Routes} from "react-router-dom";
import Items from "../items/Items.tsx";
import Dashboard from "../dashboard/Dashboard.tsx";

const MotionFlex = motion.create(Flex)

export default function Main() {
    return (
        <MotionFlex layout
                    flex={1}
                    h="100%"
                    overflow="hidden"
                    pl="4"
                    pr="4"
                    pt="4"
                    pb="2"
                    direction="column"
                    backdropFilter="blur(8px)"
                    bg="rgba(0,0,0,0.25)"
                    minHeight={0}>
            <Routes>
                <Route path="/" element={<Dashboard/>}/>
                <Route path="/items" element={<Items/>}/>
            </Routes>
        </MotionFlex>
    )
}