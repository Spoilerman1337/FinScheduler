import ActivityFilter from '../../../components/listingFilters/ActivityFilter.tsx';
import DateRangeFilter from '../../../components/listingFilters/DateRangeFilter.tsx';
import FilterWrapper from '../../../components/listingFilters/FilterWrapper.tsx';
import NumberRangeFilter, {
    type NumberRangeValue,
} from '../../../components/listingFilters/NumberRangeFilter.tsx';
import TextInputFilter from '../../../components/listingFilters/TextInputFilter.tsx';
import type {DateRangeFilterValue} from '../../../components/listingFilters/DateRangeFilter.tsx';
import type {ItemDateFilterMode} from '../../../api/items.types.ts';
import type {ActivityFilterValue} from '../../../components/listingFilters/shared.ts';
import {itemDateModes} from '../types.ts';

type ItemsFiltersProps = {
    searchTerm: string;
    statusFilter: ActivityFilterValue;
    dateFilter: DateRangeFilterValue<ItemDateFilterMode>;
    priceRange: NumberRangeValue;
    cashbackRange: NumberRangeValue;
    onSearchTermChange: (value: string) => void;
    onStatusFilterChange: (value: ActivityFilterValue) => void;
    onDateFilterChange: (value: DateRangeFilterValue<ItemDateFilterMode>) => void;
    onPriceRangeChange: (value: NumberRangeValue) => void;
    onCashbackRangeChange: (value: NumberRangeValue) => void;
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
            <DateRangeFilter
                label="Дата"
                modes={itemDateModes}
                value={props.dateFilter}
                onChange={(value) => {
                    props.onDateFilterChange({
                        mode: value.mode as ItemDateFilterMode,
                        from: value.from,
                        to: value.to,
                    });
                }}
            />
            <NumberRangeFilter
                label="Цена"
                description="Укажите диапазон, который применится к списку."
                value={props.priceRange}
                onChange={props.onPriceRangeChange}
                min={0}
                suffix="₽"
            />
            <NumberRangeFilter
                label="Кэшбэк"
                description="Укажите диапазон кэшбэка, который применится к списку."
                value={props.cashbackRange}
                onChange={props.onCashbackRangeChange}
                min={0}
                suffix="%"
            />
        </FilterWrapper>
    );
}
