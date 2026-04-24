import {Box, Button, Flex} from "@chakra-ui/react";
import {filterWidthProps, type ActivityFilterValue} from "./shared.ts";

type ActivityFilterProps = {
    value: ActivityFilterValue;
    onChange: (value: ActivityFilterValue) => void;
};

const statusOptions: ActivityFilterValue[] = ['All', 'Active', 'Inactive'];

export default function ActivityFilter(props: ActivityFilterProps) {
    return (
        <Box {...filterWidthProps}>
            <Flex gap={1} borderRadius="md" p={1} bg="bg.layer1" border="1px solid" borderColor="glass.border">
                {statusOptions.map((status) => (
                    <Button
                        key={status}
                        size="sm"
                        flex={1}
                        onClick={() => props.onChange(status)}
                        color="neon.blue"
                        borderColor={props.value === status ? "neon.blue" : "glass.border"}
                        backdropFilter="blur(12px)"
                        bg={props.value === status ? "glass.bgHover" : "transparent"}
                        border="1px solid"
                        transition="all 0.3s ease-in-out"
                        filter={props.value === status ? "drop-shadow(0 0 16px rgba(0,212,255,0.9))" : "none"}
                        boxShadow={props.value === status ? "0 0 20px rgba(0,212,255,1)" : "none"}
                        _hover={{
                            filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                            boxShadow: "0 0 20px rgba(212, 0,255,1)",
                            color: "neon.purple",
                            bg: "glass.bgHover",
                            backdropFilter: "blur(12px)",
                            borderColor: "neon.purple",
                        }}
                        _active={{
                            filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                            boxShadow: "0 0 20px rgba(0,212,255,1)",
                            color: "neon.blue",
                            bg: "glass.bgHover",
                            backdropFilter: "blur(12px)",
                            borderColor: "neon.blue",
                        }}
                        focusRing="none"
                    >
                        {status === 'All' ? 'Все' : status === 'Active' ? 'Активные' : 'Неактивные'}
                    </Button>
                ))}
            </Flex>
        </Box>
    );
}
