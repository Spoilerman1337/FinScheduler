import {DatePicker, Flex, Grid, parseDate} from '@chakra-ui/react';
import {useState} from 'react';
import type {DateValue} from '@chakra-ui/react';
import {
    buildDayKey,
    buildMonthPreviewDetails,
    getInitialSelectedDayKey,
    startOfDay,
    startOfMonth,
} from './shared.ts';
import CalendarGridPanel from './CalendarGridPanel.tsx';
import CalendarMonthHeader from './CalendarMonthHeader.tsx';
import CalendarInfoPanel from './CalendarInfoPanel.tsx';

function dateValueToDate(value: DateValue) {
    return value.toDate('UTC');
}

function buildDateValueKey(value: DateValue) {
    return buildDayKey(dateValueToDate(value));
}

export default function Calendar() {
    const [today] = useState(() => startOfDay(new Date()));
    const [visibleMonthStart, setVisibleMonthStart] = useState(() => startOfMonth(today));
    const [selectedDayKey, setSelectedDayKey] = useState(() => {
        const initialMonth = startOfMonth(today);
        const initialDetails = buildMonthPreviewDetails(initialMonth, today);

        return getInitialSelectedDayKey(initialMonth, initialDetails, today);
    });
    const [focusedValue, setFocusedValue] = useState<DateValue>(() => parseDate(selectedDayKey));

    const detailsByDay = buildMonthPreviewDetails(visibleMonthStart, today);
    const selectedDate = dateValueToDate(parseDate(selectedDayKey));
    const selectedDetails = detailsByDay[selectedDayKey] ?? null;

    function selectMonthDefaults(nextMonthStart: Date) {
        const nextDetails = buildMonthPreviewDetails(nextMonthStart, today);
        const nextSelectedDayKey = getInitialSelectedDayKey(nextMonthStart, nextDetails, today);

        setVisibleMonthStart(nextMonthStart);
        setSelectedDayKey(nextSelectedDayKey);
        setFocusedValue(parseDate(nextSelectedDayKey));
    }

    function handleTodayClick() {
        const todayMonthStart = startOfMonth(today);
        const todayKey = buildDayKey(today);

        setVisibleMonthStart(todayMonthStart);
        setSelectedDayKey(todayKey);
        setFocusedValue(parseDate(todayKey));
    }

    return (
        <DatePicker.Root
            inline
            open
            unstyled
            locale="ru-RU"
            closeOnSelect={false}
            startOfWeek={1}
            selectionMode="single"
            defaultView="day"
            minView="day"
            maxView="day"
            value={[parseDate(selectedDayKey)]}
            focusedValue={focusedValue}
            onValueChange={({value}) => {
                if (value[0]) {
                    const nextSelectedDayKey = buildDateValueKey(value[0]);

                    setSelectedDayKey(nextSelectedDayKey);
                    setFocusedValue(value[0]);
                }
            }}
            onFocusChange={({focusedValue: nextFocusedValue}) => {
                setFocusedValue(nextFocusedValue);
            }}
            onVisibleRangeChange={({visibleRange}) => {
                const nextMonthStart = startOfMonth(dateValueToDate(visibleRange.start));

                if (buildDayKey(nextMonthStart) !== buildDayKey(visibleMonthStart)) {
                    selectMonthDefaults(nextMonthStart);
                }
            }}
        >
            <Flex direction="column" gap="5" w="full" minW="0" pb="4">
                <CalendarMonthHeader
                    visibleMonthStart={visibleMonthStart}
                    onTodayClick={handleTodayClick}
                />

                <Grid
                    templateColumns={{base: '1fr', xl: 'minmax(0, 1fr) 24rem'}}
                    gap="5"
                    alignItems="start"
                >
                    <CalendarGridPanel detailsByDay={detailsByDay} />

                    <CalendarInfoPanel selectedDate={selectedDate} details={selectedDetails} />
                </Grid>
            </Flex>
        </DatePicker.Root>
    );
}
