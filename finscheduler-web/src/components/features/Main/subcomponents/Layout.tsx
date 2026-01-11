import {Flex, IconButton, Text} from "@chakra-ui/react";
import {FiSearch} from "react-icons/fi";

export interface LayoutProps {
    children?: React.ReactNode;
    headerName?: string;
}

export default function Layout(props: LayoutProps) {
    return (
        <>
            <Flex justify="space-between" align="center" mb={7}>
                <Text fontSize="3xl" color="white" fontFamily="body">{props.headerName}</Text>
                <IconButton aria-label="search"
                            children={<FiSearch/>}
                            bg="whiteAlpha.300"
                            _hover={{bg: "whiteAlpha.400"}}/>
            </Flex>
            <Flex minHeight={0}>
                {props.children}
            </Flex>
        </>
    )
}