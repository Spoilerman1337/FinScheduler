import {HStack, Text} from '@chakra-ui/react';
import {Timer} from 'lucide-react';
import type {CalendarEventColorPreset} from './shared.ts';
import {calendarEventColorStyles} from './calendarUi.ts';

interface CalendarDayTimeChipProps {
    ariaLabel: string;
    timeLabel: string;
    color: CalendarEventColorPreset;
}

export default function CalendarDayTimeChip(props: CalendarDayTimeChipProps) {
    const {ariaLabel, timeLabel, color} = props;

    return (
        <HStack
            as="span"
            aria-label={ariaLabel}
            gap="1"
            h="1.25rem"
            px="1.5"
            borderWidth="1px"
            borderRadius="full"
            borderColor={calendarEventColorStyles[color].text}
            bg="rgba(6, 16, 34, 0.2)"
            color={calendarEventColorStyles[color].text}
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
