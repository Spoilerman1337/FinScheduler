import ActivityFilter from '../../../components/listingFilters/ActivityFilter.tsx';
import FilterWrapper from '../../../components/listingFilters/FilterWrapper.tsx';
import NumberRangeFilter, {
    type NumberRangeValue,
} from '../../../components/listingFilters/NumberRangeFilter.tsx';
import TextInputFilter from '../../../components/listingFilters/TextInputFilter.tsx';
import type {ActivityFilterValue} from '../../../components/listingFilters/shared.ts';

type ItemsFiltersProps = {
    searchTerm: string;
    statusFilter: ActivityFilterValue;
    priceRange: NumberRangeValue;
    onSearchTermChange: (value: string) => void;
    onStatusFilterChange: (value: ActivityFilterValue) => void;
    onPriceRangeChange: (value: NumberRangeValue) => void;
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
            <ActivityFilter value={props.statusFilter} onChange={props.onStatusFilterChange} />
            <NumberRangeFilter
                label="Цена"
                description="Укажите диапазон, который применится к списку."
                value={props.priceRange}
                onChange={props.onPriceRangeChange}
                min={0}
                suffix="₽"
            />
        </FilterWrapper>
    );
}
