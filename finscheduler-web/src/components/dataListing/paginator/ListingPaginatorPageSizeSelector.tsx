import {Box, Flex, Select, Text} from "@chakra-ui/react";
import {CheckIcon, ChevronDownIcon} from "lucide-react";
import {createPageSizeCollection, defaultPageSizeOptions} from "../types.ts";

interface ListingPaginatorPageSizeSelectorProps {
    pageSize: number;
    pageSizeOptions?: readonly number[];
    onPageSizeChange: (pageSize: number) => void;
}

export default function ListingPaginatorPageSizeSelector(props: ListingPaginatorPageSizeSelectorProps) {
    const pageSizeOptions = props.pageSizeOptions ?? defaultPageSizeOptions;
    const pageSizeCollection = createPageSizeCollection(pageSizeOptions);
    const fallbackPageSize = pageSizeOptions[0] ?? defaultPageSizeOptions[0];

    return (
        <Select.Root
            collection={pageSizeCollection}
            width="60px"
            size="lg"
            value={[props.pageSize.toString()]}
            onValueChange={(event) => {
                const newPageSize = parseInt(event.value[0] || fallbackPageSize.toString(), 10);
                props.onPageSizeChange(newPageSize);
            }}
        >
            <Select.HiddenSelect/>
            <Select.Label
                color="neon.blue"
                mb="2"
                fontWeight="medium"
                _empty={{display: "none"}}
            />

            <Select.Control asChild>
                <Flex align="center" justify="space-between">
                    <Select.Trigger asChild>
                        <Box
                            flex="1"
                            borderColor="neon.blue"
                            boxShadow="0 0 10px rgba(0, 212, 255, 0.6)"
                            transition="all 0.3s ease-in-out"
                            backdropFilter="blur(12px)"
                            bg="glass.bg"
                            color="neon.blue"
                            p="1.5"
                            maxW="100px"
                            h="36px"
                            focusRing="none"
                            _hover={{
                                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                color: "neon.purple",
                                bg: "glass.bgHover",
                                borderColor: "neon.purple",
                                cursor: "pointer",
                            }}
                        >
                            <Select.ValueText
                                placeholder={props.pageSize.toString()}
                                color="currentColor"
                                _placeholder={{color: "text.placeholder"}}
                                fontSize="sm"
                            />
                        </Box>
                    </Select.Trigger>

                    <Select.IndicatorGroup asChild>
                        <Flex align="center">
                            <Select.Indicator asChild>
                                <ChevronDownIcon/>
                            </Select.Indicator>
                        </Flex>
                    </Select.IndicatorGroup>
                </Flex>
            </Select.Control>

            <Select.Positioner>
                <Select.Content
                    p="1"
                    borderRadius="sm"
                    zIndex="dropdown"
                    mt="1"
                    border="1px solid"
                    borderColor="glass.borderStrong"
                    backdropFilter="blur(16px)"
                    bg="glass.bg"
                    boxShadow="lg"
                    width="--trigger-width"
                    maxH="200px"
                    overflowY="auto"
                    className="custom-scrollbar"
                >
                    {pageSizeCollection.items.map((item) => {
                        const isSelected = item.value === props.pageSize;

                        return (
                            <Select.Item
                                item={item}
                                key={item.value}
                                display="flex"
                                justifyContent="space-between"
                                alignItems="center"
                                py="1.5"
                                px="2"
                                borderRadius="sm"
                                color="neon.blue"
                                bg={isSelected ? "glass.bgHover" : "transparent"}
                                filter={isSelected ? "drop-shadow(0 0 8px rgba(0,212,255,0.9))" : "none"}
                                fontSize="sm"
                                fontWeight={isSelected ? "semibold" : "normal"}
                                transition="all 0.2s"
                                cursor="pointer"
                                _hover={{
                                    filter: "drop-shadow(0 0 8px rgba(212, 0,255,0.9))",
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                }}
                                _highlighted={{
                                    filter: "drop-shadow(0 0 8px rgba(212, 0,255,0.9))",
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                }}
                                _selected={{
                                    filter: "drop-shadow(0 0 8px rgba(0,212,255,0.9))",
                                    color: "neon.blue",
                                    bg: "glass.bgHover",
                                    fontWeight: "semibold",
                                }}
                                _focus={{
                                    outline: "none",
                                    bg: "glass.bgHover",
                                    boxShadow: "0 0 0 2px rgba(0, 212, 255, 0.5)",
                                }}
                            >
                                <Text color="currentColor">{item.label}</Text>
                                <Select.ItemIndicator asChild>
                                    <CheckIcon/>
                                </Select.ItemIndicator>
                            </Select.Item>
                        );
                    })}
                </Select.Content>
            </Select.Positioner>
        </Select.Root>
    );
}
