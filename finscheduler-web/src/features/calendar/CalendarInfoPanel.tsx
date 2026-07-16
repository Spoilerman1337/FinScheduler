import {Box, Flex, HStack, Text, VStack} from '@chakra-ui/react';
import {Clock3} from 'lucide-react';
import type {CalendarDayDetails, CalendarMarkerTone} from './shared.ts';
import {formatSelectedDayLabel} from './shared.ts';
import {markerToneStyles} from './calendarUi.ts';

interface CalendarInfoPanelProps {
    selectedDate: Date;
    details: CalendarDayDetails | null;
}

interface CalendarEventTimeListProps {
    timeHints: string[];
    tone: CalendarMarkerTone;
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

function CalendarDetailsRail(props: CalendarInfoPanelProps) {
    const {selectedDate, details} = props;

    return (
        <Box
            as="aside"
            aria-label={'Сведения по выбранному дню'}
            borderWidth="1px"
            borderColor="app.cardBorder"
            bg="app.cardBg"
            borderRadius="2xl"
            boxShadow="card"
            p={{base: '5', md: '6'}}
            display="flex"
            flexDirection="column"
            gap="5"
            alignSelf="start"
            position={{base: 'static', xl: 'sticky'}}
            top="4"
        >
            <Flex justify="space-between" align="start" gap="4">
                <VStack align="start" gap="1.5">
                    <Text textStyle="eyebrow" color="app.accent">
                        {'Контекст дня'}
                    </Text>
                    <Text fontSize="2xl" fontWeight="700" lineHeight="1.2">
                        {formatSelectedDayLabel(selectedDate)}
                    </Text>
                </VStack>
            </Flex>

            {details ? (
                <VStack align="stretch" gap="3">
                    {details.markers.map((marker) => (
                        <Box
                            key={marker.id}
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
                                        <CalendarEventTimeList
                                            timeHints={marker.timeHints}
                                            tone={marker.tone}
                                        />
                                    </Box>
                                ) : null}
                            </VStack>
                        </Box>
                    ))}
                </VStack>
            ) : (
                <Box
                    borderWidth="1px"
                    borderStyle="dashed"
                    borderColor="app.cardBorder"
                    borderRadius="xl"
                    px="4"
                    py="5"
                    bg="rgba(6, 16, 34, 0.2)"
                >
                    <Text fontWeight="700" mb="1.5">
                        {'Событий на эту дату нет'}
                    </Text>
                    <Text color="fg.muted">
                        {'Эта дата пока не привязана ни к одному событию.'}
                    </Text>
                </Box>
            )}
        </Box>
    );
}

export default function CalendarInfoPanel(props: CalendarInfoPanelProps) {
    const {selectedDate, details} = props;

    return <CalendarDetailsRail selectedDate={selectedDate} details={details} />;
}
