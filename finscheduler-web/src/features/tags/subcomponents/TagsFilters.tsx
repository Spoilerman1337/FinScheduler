import {Box, Button, Flex, Input} from "@chakra-ui/react";
import {SearchIcon} from "lucide-react";

type TagsFiltersProps = {
    searchTerm: string;
    statusFilter: 'All' | 'Active' | 'Inactive';
    onSearchTermChange: (value: string) => void;
    onStatusFilterChange: (value: 'All' | 'Active' | 'Inactive') => void;
    onApply: () => void;
    onReset: () => void;
};

const filterWidthProps = {
    w: {base: '100%', md: 'calc(50% - var(--chakra-space-2))', xl: 'calc(20% - var(--chakra-space-3))'},
};

const statusOptions: Array<'All' | 'Active' | 'Inactive'> = ['All', 'Active', 'Inactive'];

export default function TagsFilters(props: TagsFiltersProps) {
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
            <Box {...filterWidthProps} position="relative">
                <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1"
                     pointerEvents="none">
                    <SearchIcon size={18} color="rgba(255,255,255,0.6)"/>
                </Box>
                <Input
                    placeholder="Поиск по названию..."
                    value={props.searchTerm}
                    onChange={(e) => props.onSearchTermChange(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            props.onApply();
                        }
                    }}
                    pl="10"
                    bg="bg.layer1"
                    borderColor="glass.border"
                    color="neon.blue"
                    _placeholder={{color: 'textMuted'}}
                />
            </Box>

            <Box {...filterWidthProps}>
                <Flex gap={1} borderRadius="md" p={1} bg="bg.layer1" border="1px solid" borderColor="glass.border">
                    {statusOptions.map((status) => (
                        <Button
                            key={status}
                            size="sm"
                            flex={1}
                            onClick={() => props.onStatusFilterChange(status)}
                            color={props.statusFilter === status ? "neon.blue" : "neon.blue"}
                            borderColor={props.statusFilter === status ? "neon.blue" : "glass.border"}
                            backdropFilter="blur(12px)"
                            bg={props.statusFilter === status ? "glass.bgHover" : "transparent"}
                            border="1px solid"
                            transition="all 0.3s ease-in-out"
                            filter={props.statusFilter === status ? "drop-shadow(0 0 16px rgba(0,212,255,0.9))" : "none"}
                            boxShadow={props.statusFilter === status ? "0 0 20px rgba(0,212,255,1)" : "none"}
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
