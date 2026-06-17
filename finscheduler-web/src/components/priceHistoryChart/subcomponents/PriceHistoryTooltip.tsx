import {Box, Flex, Stack, Text} from '@chakra-ui/react';
import {PRICE_HISTORY_TOOLTIP_WIDTH_PX} from '../shared.ts';
import type {PriceHistoryTooltipEntry} from '../shared.ts';

interface PriceHistoryTooltipProps {
    label: string;
    entries: PriceHistoryTooltipEntry[];
}

export default function PriceHistoryTooltip({label, entries}: PriceHistoryTooltipProps) {
    return (
        <Box
            width={`${PRICE_HISTORY_TOOLTIP_WIDTH_PX}px`}
            maxW="calc(100% - 16px)"
            borderRadius="xl"
            border="1px solid"
            borderColor="rgba(143, 120, 255, 0.26)"
            bg="rgba(8, 16, 32, 0.94)"
            boxShadow="0 20px 40px rgba(3, 8, 20, 0.42)"
            px={4}
            py={3}
        >
            <Text color="fg" fontWeight="600" mb={2}>
                {label}
            </Text>
            <Stack gap={1}>
                {entries.map((entry) => {
                    const ChangeIcon = entry.changeIcon;

                    return (
                        <Stack key={entry.key} gap={0.5}>
                            <Text color={entry.valueColor ?? 'fg'} fontWeight="700" textStyle="sm">
                                {entry.value}
                            </Text>
                            {entry.changeSummary && ChangeIcon ? (
                                <Flex align="center" gap={1.5}>
                                    <Box color={entry.changeColor} flexShrink={0}>
                                        <ChangeIcon size={12} strokeWidth={2.5} />
                                    </Box>
                                    <Text color={entry.changeColor} fontWeight="600" textStyle="xs">
                                        {entry.changeSummary}
                                    </Text>
                                </Flex>
                            ) : null}
                        </Stack>
                    );
                })}
            </Stack>
        </Box>
    );
}
