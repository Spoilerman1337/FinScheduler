import {Box, Flex, HStack, Text, VStack} from '@chakra-ui/react';
import {Clock3} from 'lucide-react';
import type {CalendarMarkerPreview} from './shared.ts';
import {markerToneStyles} from './calendarUi.ts';

interface CalendarInfoPanelEventCardProps {
    marker: CalendarMarkerPreview;
}

interface CalendarEventTimeListProps {
    timeHints: string[];
    tone: CalendarMarkerPreview['tone'];
}

function CalendarEventTimeList(props: CalendarEventTimeListProps) {
    const {timeHints, tone} = props;

    if (timeHints.length === 0) {
        return null;
    }

    return (
        <Flex wrap="wrap" gap="2">
            {timeHints.map((timeHint, index) => (
                <Box
                    key={`${timeHint}-${index}`}
                    px="2.5"
                    py="1"
                    borderWidth="1px"
                    borderColor="app.cardBorder"
                    borderRadius="full"
                    bg="rgba(6, 16, 34, 0.2)"
                >
                    <HStack gap="1.5" color={markerToneStyles[tone].text}>
                        <Clock3 size={12} />
                        <Text fontSize="xs" fontWeight="700">
                            {timeHint}
                        </Text>
                    </HStack>
                </Box>
            ))}
        </Flex>
    );
}

export default function CalendarInfoPanelEventCard(props: CalendarInfoPanelEventCardProps) {
    const {marker} = props;

    return (
        <Box
            borderWidth="1px"
            borderColor="app.cardBorder"
            borderRadius="xl"
            bg="rgba(6, 16, 34, 0.34)"
            px="4"
            py="3.5"
        >
            <VStack align="stretch" gap="3">
                <HStack align="start" gap="3">
                    <Box
                        mt="1"
                        boxSize="2.5"
                        borderRadius="full"
                        bg={markerToneStyles[marker.tone].bg}
                        flexShrink={0}
                    />
                    <VStack align="start" gap="1" flex="1">
                        <Text fontWeight="700">{marker.label}</Text>
                        <Text color="fg.muted" fontSize="sm">
                            {marker.description}
                        </Text>
                    </VStack>
                </HStack>

                {marker.timeHints.length > 0 ? (
                    <Box pl="5.5">
                        <CalendarEventTimeList timeHints={marker.timeHints} tone={marker.tone} />
                    </Box>
                ) : null}
            </VStack>
        </Box>
    );
}