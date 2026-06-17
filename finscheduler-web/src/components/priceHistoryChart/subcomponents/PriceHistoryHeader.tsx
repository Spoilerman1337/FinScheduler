import {Card, Flex, Stack} from '@chakra-ui/react';
import type {PriceHistoryLegendItem} from '../shared.ts';
import PriceHistorySeriesLegend from './PriceHistorySeriesLegend.tsx';

interface PriceHistoryHeaderProps {
    title: string;
    description: string;
    legends: PriceHistoryLegendItem[];
}

export default function PriceHistoryHeader({title, description, legends}: PriceHistoryHeaderProps) {
    return (
        <Card.Header position="relative">
            <Flex
                direction={{base: 'column', md: 'row'}}
                align={{base: 'flex-start', md: 'center'}}
                justify="space-between"
                gap={3}
            >
                <Stack gap={1}>
                    <Card.Title>{title}</Card.Title>
                    <Card.Description>{description}</Card.Description>
                </Stack>
                <Flex wrap="wrap" gap={3}>
                    {legends.map((legend) => (
                        <PriceHistorySeriesLegend
                            key={legend.label}
                            label={legend.label}
                            color={legend.color}
                            lineStyle={legend.lineStyle}
                        />
                    ))}
                </Flex>
            </Flex>
        </Card.Header>
    );
}
