import {Box, Flex, Text} from '@chakra-ui/react';

interface PriceHistorySeriesLegendProps {
    label: string;
    color: string;
    lineStyle: 'solid' | 'dashed';
}

export default function PriceHistorySeriesLegend({
    label,
    color,
    lineStyle,
}: PriceHistorySeriesLegendProps) {
    return (
        <Flex align="center" gap={2}>
            <Box
                width="6"
                height="0"
                borderTopWidth="2px"
                borderTopColor={color}
                borderTopStyle={lineStyle === 'dashed' ? 'dashed' : 'solid'}
                boxShadow={`0 0 14px ${color}`}
                flexShrink={0}
            />
            <Text color="fg.subtle" textStyle="sm">
                {label}
            </Text>
        </Flex>
    );
}
