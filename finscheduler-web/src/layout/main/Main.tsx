import {Flex} from "@chakra-ui/react";
import {motion} from "framer-motion";
import {Route, Routes} from "react-router-dom";
import Items from "../../features/items/Items.tsx";
import Dashboard from "../../features/dashboard/Dashboard.tsx";
import Tags from "../../features/tags/Tags.tsx";

const MotionFlex = motion.create(Flex)

export default function Main() {
    return (
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
            <Routes>
                <Route path="/" element={<Dashboard/>}/>
                <Route path="/items" element={<Items/>}/>
                <Route path="/tags" element={<Tags/>}/>
            </Routes>
        </MotionFlex>
    )
}
