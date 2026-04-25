import ActivityFilter from "../../../components/listingFilters/ActivityFilter.tsx";
import FilterWrapper from "../../../components/listingFilters/FilterWrapper.tsx";
import NumberInputFilter from "../../../components/listingFilters/NumberInputFilter.tsx";
import TextInputFilter from "../../../components/listingFilters/TextInputFilter.tsx";
import type {ActivityFilterValue} from "../../../components/listingFilters/shared.ts";

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
        <FilterWrapper onReset={props.onReset}>
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
        </FilterWrapper>
    );
}
