import {Flex, Text} from "@chakra-ui/react";
import React from "react";

export interface LayoutProps {
    children?: React.ReactNode;
    headerName?: string;
}

export default function Layout(props: LayoutProps) {
    return (
        <>
            <Flex justify="space-between" align="center" mb={7}>
                <Text fontSize="3xl" color="white" fontFamily="body">{props.headerName}</Text>
            </Flex>
            <Flex minHeight={0}>
                {props.children}
            </Flex>
        </>
    )
}