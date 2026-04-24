import {Button, Flex} from "@chakra-ui/react";
import ActivityFilter from "../../../components/listingFilters/ActivityFilter.tsx";
import NumberInputFilter from "../../../components/listingFilters/NumberInputFilter.tsx";
import TextInputFilter from "../../../components/listingFilters/TextInputFilter.tsx";
import {filterWidthProps, type ActivityFilterValue} from "../../../components/listingFilters/shared.ts";

type ItemsFiltersProps = {
    searchTerm: string;
    statusFilter: ActivityFilterValue;
    priceFrom: string;
    priceTo: string;
    onSearchTermChange: (value: string) => void;
    onStatusFilterChange: (value: ActivityFilterValue) => void;
    onPriceFromChange: (value: string) => void;
    onPriceToChange: (value: string) => void;
    onApply: () => void;
    onReset: () => void;
};

export default function ItemsFilters(props: ItemsFiltersProps) {
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
            <NumberInputFilter
                value={props.priceFrom}
                placeholder="Цена от"
                onChange={props.onPriceFromChange}
                onApply={props.onApply}
            />
            <NumberInputFilter
                value={props.priceTo}
                placeholder="Цена до"
                onChange={props.onPriceToChange}
                onApply={props.onApply}
            />
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
