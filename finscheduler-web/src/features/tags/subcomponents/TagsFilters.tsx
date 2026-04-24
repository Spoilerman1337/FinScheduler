import ActivityFilter from "../../../components/listingFilters/ActivityFilter.tsx";
import FilterWrapper from "../../../components/listingFilters/FilterWrapper.tsx";
import TextInputFilter from "../../../components/listingFilters/TextInputFilter.tsx";
import type {ActivityFilterValue} from "../../../components/listingFilters/shared.ts";

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
        <FilterWrapper onReset={props.onReset}>
            <TextInputFilter
                value={props.searchTerm}
                placeholder="Поиск по названию..."
                onChange={props.onSearchTermChange}
                onApply={props.onApply}
            />
            <ActivityFilter value={props.statusFilter} onChange={props.onStatusFilterChange}/>
        </FilterWrapper>
    );
}
