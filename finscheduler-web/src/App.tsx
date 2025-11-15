import {Provider} from "./components/ui/provider.tsx";
import {Box, AbsoluteCenter, Text, Theme} from "@chakra-ui/react";

export default function App() {
    return (<Provider forcedTheme={"dark"}>
        <Theme>
            <Box h="100vh" w="100vw" bg="bg">
                <AbsoluteCenter>
                    <Text textStyle="7xl" fontFamily="body">
                        Hello, World! This is a skeleton of my frontend app - finscheduler
                    </Text>
                </AbsoluteCenter>
            </Box>
        </Theme>
    </Provider>);
}
