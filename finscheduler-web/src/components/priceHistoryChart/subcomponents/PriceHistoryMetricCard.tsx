import {Flex, Text} from '@chakra-ui/react';

interface PriceHistoryMetricCardProps {
    label: string;
    value: string;
    accentColor: string;
}

export default function PriceHistoryMetricCard({
    label,
    value,
    accentColor,
}: PriceHistoryMetricCardProps) {
    return (
        <Flex
            direction="column"
            gap={0.5}
            px={3}
            py={2.5}
            aria-label={`${label}: ${value}`}
            borderRadius="lg"
            border="1px solid"
            borderColor="app.cardBorder"
            bg="rgba(8, 16, 32, 0.46)"
            boxShadow="inset 0 1px 0 rgba(255, 255, 255, 0.03)"
        >
            <Text
                textStyle="2xs"
                color="fg.subtle"
                textTransform="uppercase"
                letterSpacing="0.08em"
            >
                {label}
            </Text>
            <Text color={accentColor} fontWeight="700" textStyle="md">
                {value}
            </Text>
        </Flex>
    );
}
