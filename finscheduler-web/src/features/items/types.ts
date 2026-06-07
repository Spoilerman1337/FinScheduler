import type {
    DateRangeFilterMode,
    DateRangeFilterValue,
} from '../../components/listingFilters/DateRangeFilter.tsx';
import type {ItemDateFilterMode} from '../../api/items.types.ts';

export const itemDateModes: readonly [
    DateRangeFilterMode<ItemDateFilterMode>,
    ...DateRangeFilterMode<ItemDateFilterMode>[],
] = [
    {
        value: 'created',
        label: 'Создан',
        description: 'Показывать предметы, созданные в выбранный период.',
        fromInputLabel: 'Дата создания от',
        toInputLabel: 'Дата создания до',
    },
    {
        value: 'updated',
        label: 'Обновлён',
        description: 'Показывать предметы, обновлённые в выбранный период.',
        fromInputLabel: 'Дата обновления от',
        toInputLabel: 'Дата обновления до',
    },
];

export function createDefaultItemDateFilter(): DateRangeFilterValue<ItemDateFilterMode> {
    return {
        mode: 'created',
        from: '',
        to: '',
    };
}

export function getCashbackColor(cashback: number | undefined): string {
    let color = 'neon.pink';

    if (cashback) {
        color = cashback > 1 ? 'neon.green' : 'neon.yellow';
    }

    return color;
}
