import {Box, Flex, Text} from '@chakra-ui/react';

interface PriceHistorySeriesLegendProps {
    label: string;
    color: string;
}

export default function PriceHistorySeriesLegend({label, color}: PriceHistorySeriesLegendProps) {
    return (
        <Flex align="center" gap={2}>
            <Box
                width="2.5"
                height="2.5"
                borderRadius="full"
                bg={color}
                boxShadow={`0 0 14px ${color}`}
                flexShrink={0}
            />
            <Text color="fg.subtle" textStyle="sm">
                {label}
            </Text>
        </Flex>
    );
}
