import {SimpleGrid} from '@chakra-ui/react';
import type {PriceHistoryMetricItem} from '../shared.ts';
import PriceHistoryMetricCard from './PriceHistoryMetricCard.tsx';

interface PriceHistoryMetricsProps {
    items: PriceHistoryMetricItem[];
}

export default function PriceHistoryMetrics({items}: PriceHistoryMetricsProps) {
    return (
        <SimpleGrid columns={3} gap={2.5}>
            {items.map((item) => (
                <PriceHistoryMetricCard
                    key={item.label}
                    label={item.label}
                    value={item.value}
                    accentColor={item.accentColor}
                />
            ))}
        </SimpleGrid>
    );
}
