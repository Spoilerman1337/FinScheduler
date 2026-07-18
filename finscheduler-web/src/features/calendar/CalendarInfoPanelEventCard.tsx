import {Badge, Box, Flex, HStack, Text, VStack} from '@chakra-ui/react';
import {Clock3} from 'lucide-react';
import type {CalendarEventColorPreset} from './shared.ts';
import {calendarEventColorStyles} from './calendarUi.ts';

interface CalendarInfoPanelEventCardProps {
    color: CalendarEventColorPreset;
    description: string;
    label: string;
    timeHints: string[];
    triggerTags: string[];
}

interface CalendarEventTimeListProps {
    color: CalendarEventColorPreset;
    timeHints: string[];
}

function CalendarEventTimeList(props: CalendarEventTimeListProps) {
    const {color, timeHints} = props;

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
                    <HStack gap="1.5" color={calendarEventColorStyles[color].text}>
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

function CalendarEventTriggerTags(props: {
    color: CalendarEventColorPreset;
    triggerTags: string[];
}) {
    const {color, triggerTags} = props;

    if (triggerTags.length === 0) {
        return null;
    }

    return (
        <Flex wrap="wrap" gap="2">
            {triggerTags.map((triggerTag) => (
                <Badge
                    key={triggerTag}
                    px="2.5"
                    py="1"
                    borderRadius="full"
                    borderWidth="1px"
                    borderColor="app.cardBorder"
                    bg="rgba(6, 16, 34, 0.2)"
                    color={calendarEventColorStyles[color].text}
                >
                    {triggerTag}
                </Badge>
            ))}
        </Flex>
    );
}

export default function CalendarInfoPanelEventCard(props: CalendarInfoPanelEventCardProps) {
    const {color, description, label, timeHints, triggerTags} = props;

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
                        bg={calendarEventColorStyles[color].bg}
                        flexShrink={0}
                    />
                    <VStack align="start" gap="1" flex="1">
                        <Text fontWeight="700">{label}</Text>
                        <Text color="fg.muted" fontSize="sm">
                            {description}
                        </Text>
                    </VStack>
                </HStack>

                <Box pl="5.5">
                    <VStack align="stretch" gap="3">
                        <CalendarEventTriggerTags color={color} triggerTags={triggerTags} />
                        <CalendarEventTimeList color={color} timeHints={timeHints} />
                    </VStack>
                </Box>
            </VStack>
        </Box>
    );
}
