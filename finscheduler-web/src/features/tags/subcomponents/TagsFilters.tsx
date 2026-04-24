import {Button, Flex} from "@chakra-ui/react";
import ActivityFilter from "../../../components/listingFilters/ActivityFilter.tsx";
import TextInputFilter from "../../../components/listingFilters/TextInputFilter.tsx";
import {filterWidthProps, type ActivityFilterValue} from "../../../components/listingFilters/shared.ts";

type TagsFiltersProps = {
    searchTerm: string;
    statusFilter: ActivityFilterValue;
    onSearchTermChange: (value: string) => void;
    onStatusFilterChange: (value: ActivityFilterValue) => void;
    onApply: () => void;
    onReset: () => void;
};

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
            <TextInputFilter
                value={props.searchTerm}
                placeholder="Поиск по названию..."
                onChange={props.onSearchTermChange}
                onApply={props.onApply}
            />
            <ActivityFilter value={props.statusFilter} onChange={props.onStatusFilterChange}/>
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
