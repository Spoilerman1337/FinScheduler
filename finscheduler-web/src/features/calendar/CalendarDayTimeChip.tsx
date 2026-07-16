import {HStack, Text} from '@chakra-ui/react';
import {Timer} from 'lucide-react';
import type {CalendarMarkerTone} from './shared.ts';
import {markerToneStyles} from './calendarUi.ts';

interface CalendarDayTimeChipProps {
    ariaLabel: string;
    timeLabel: string;
    tone: CalendarMarkerTone;
}

export default function CalendarDayTimeChip(props: CalendarDayTimeChipProps) {
    const {ariaLabel, timeLabel, tone} = props;

    return (
        <HStack
            as="span"
            aria-label={ariaLabel}
            gap="1"
            h="1.25rem"
            px="1.5"
            borderWidth="1px"
            borderRadius="full"
            borderColor={markerToneStyles[tone].text}
            bg="rgba(6, 16, 34, 0.2)"
            color={markerToneStyles[tone].text}
            flexShrink={0}
            whiteSpace="nowrap"
            style={{boxSizing: 'border-box'}}
        >
            <Timer size={11} />
            <Text fontSize="2xs" fontWeight="700" lineHeight="1">
                {timeLabel}
            </Text>
        </HStack>
    );
}
