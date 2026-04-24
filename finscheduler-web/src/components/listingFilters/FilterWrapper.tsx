import {Button, Flex} from "@chakra-ui/react";
import type {ReactNode} from "react";
import {filterWidthProps} from "./shared.ts";

type FilterWrapperProps = {
    children: ReactNode;
    onReset: () => void;
};

export default function FilterWrapper(props: FilterWrapperProps) {
    return (
        <Flex
            mb={4}
            p={4}
            bg="bg.layer2"
            borderRadius="lg"
            border="1px solid"
            borderColor="glass.border"
            width="100%"
            align="center"
            gap={4}
            flexWrap="wrap"
            justifyContent="flex-start"
        >
            {props.children}
            <Button
                {...filterWidthProps}
                onClick={props.onReset}
                bg="bg.layer3"
                color="textMuted"
                borderColor="glass.border"
                border="1px solid"
                _hover={{bg: 'neon.pink', color: 'bg.base'}}
            >
                Сброс
            </Button>
        </Flex>
    );
}
