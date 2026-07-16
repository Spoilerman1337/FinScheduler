import {Box, Flex, Text, VStack} from '@chakra-ui/react';
import type {CalendarDayDetails} from './shared.ts';
import {formatSelectedDayLabel} from './shared.ts';
import CalendarInfoPanelEventCard from './CalendarInfoPanelEventCard.tsx';

interface CalendarInfoPanelProps {
    selectedDate: Date;
    details: CalendarDayDetails | null;
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
                        Контекст дня
                    </Text>
                    <Text fontSize="2xl" fontWeight="700" lineHeight="1.2">
                        {formatSelectedDayLabel(selectedDate)}
                    </Text>
                </VStack>
            </Flex>

            {details ? (
                <VStack align="stretch" gap="3">
                    {details.markers.map((marker) => (
                        <CalendarInfoPanelEventCard key={marker.id} marker={marker} />
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
                        Событий на эту дату нет
                    </Text>
                    <Text color="fg.muted">Эта дата пока не привязана ни к одному событию.</Text>
                </Box>
            )}
        </Box>
    );
}

export default function CalendarInfoPanel(props: CalendarInfoPanelProps) {
    const {selectedDate, details} = props;

    return <CalendarDetailsRail selectedDate={selectedDate} details={details} />;
}