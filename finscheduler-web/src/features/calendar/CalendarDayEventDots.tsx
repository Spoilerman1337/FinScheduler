import {Box, HStack} from '@chakra-ui/react';
import type {CalendarMarkerPreview} from './shared.ts';
import {markerToneStyles} from './calendarUi.ts';

interface CalendarDayEventDotsProps {
    markers: CalendarMarkerPreview[];
}

export default function CalendarDayEventDots(props: CalendarDayEventDotsProps) {
    const {markers} = props;

    if (markers.length === 0) {
        return null;
    }

    return (
        <HStack gap="1.5">
            {markers.map((marker) => (
                <Box
                    key={marker.id}
                    boxSize="2"
                    borderRadius="full"
                    bg={markerToneStyles[marker.tone].bg}
                    boxShadow={marker.tone === 'accent' ? 'app.glowCyan' : 'none'}
                />
            ))}
        </HStack>
    );
}
