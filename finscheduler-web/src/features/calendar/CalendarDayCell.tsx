import {Box, DatePicker, Flex, HStack, Text, VStack, useDatePickerContext} from '@chakra-ui/react';
import type {DateValue} from '@chakra-ui/react';
import type {CalendarDayDetails} from './shared.ts';
import {buildDayKey, formatDayButtonLabel, formatOverflowDayLabel} from './shared.ts';
import CalendarDayEventDots from './CalendarDayEventDots.tsx';
import CalendarDayTimeChip from './CalendarDayTimeChip.tsx';

interface CalendarDayCellProps {
    dayValue: DateValue;
    detailsByDay: Record<string, CalendarDayDetails>;
}

const fixedDayCellStyle = {
    height: '7rem',
    minHeight: '7rem',
    maxHeight: '7rem',
    boxSizing: 'border-box',
} as const;

function dateValueToDate(value: DateValue) {
    return value.toDate('UTC');
}

function buildDateValueKey(value: DateValue) {
    return buildDayKey(dateValueToDate(value));
}

function formatCountLabel(count: number, forms: [string, string, string]) {
    const mod10 = count % 10;
    const mod100 = count % 100;

    if (mod10 === 1 && mod100 !== 11) {
        return `${count} ${forms[0]}`;
    }

    if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) {
        return `${count} ${forms[1]}`;
    }

    return `${count} ${forms[2]}`;
}

function formatEventCountLabel(count: number) {
    return formatCountLabel(count, ['событие', 'события', 'событий']);
}

function toTimeHintSortValue(timeHint: string) {
    const [hours = '0', minutes = '0'] = timeHint.split(':');

    return Number(hours) * 60 + Number(minutes);
}

function getEarliestTimeHint(timeHints: string[]) {
    return [...timeHints].sort((left, right) => {
        return toTimeHintSortValue(left) - toTimeHintSortValue(right);
    })[0];
}

function getTimeChipAriaLabel(label: string, earliestTimeHint: string, hasMoreTimeHints: boolean) {
    if (hasMoreTimeHints) {
        return `У события ${label} самое раннее время ${earliestTimeHint}, есть еще слоты`;
    }

    return `У события ${label} время ${earliestTimeHint}`;
}

export default function CalendarDayCell(props: CalendarDayCellProps) {
    const {dayValue, detailsByDay} = props;
    const datePicker = useDatePickerContext();
    const date = dateValueToDate(dayValue);
    const dayKey = buildDateValueKey(dayValue);
    const details = detailsByDay[dayKey] ?? null;
    const eventCount = details?.markers.length ?? 0;
    const firstMarker = details?.markers[0] ?? null;
    const timeChipValue = firstMarker?.timeHints.length
        ? getEarliestTimeHint(firstMarker.timeHints)
        : null;
    const cellState = datePicker.getDayTableCellState({
        value: dayValue,
        visibleRange: datePicker.visibleRange,
    });
    const buttonLabel =
        eventCount > 0
            ? `Открыть ${formatDayButtonLabel(date)}, ${formatEventCountLabel(eventCount)}`
            : `Открыть ${formatDayButtonLabel(date)}, без событий`;

    return (
        <DatePicker.TableCell
            unstyled
            value={dayValue}
            visibleRange={datePicker.visibleRange}
            verticalAlign="top"
            px="0"
            py="0"
            style={fixedDayCellStyle}
        >
            <DatePicker.TableCellTrigger
                unstyled
                aria-label={buttonLabel}
                aria-pressed={cellState.selected}
                display="flex"
                flexDirection="column"
                justifyContent="space-between"
                alignItems="stretch"
                w="full"
                overflow="hidden"
                p={{base: '2.5', md: '3'}}
                borderWidth="1px"
                borderRadius="xl"
                borderColor={
                    cellState.selected
                        ? 'app.cardBorderActive'
                        : cellState.today
                          ? 'app.sidebarItemHoverBorder'
                          : 'app.cardBorder'
                }
                bg={
                    cellState.selected
                        ? 'rgba(32, 208, 255, 0.12)'
                        : cellState.outsideRange
                          ? 'rgba(6, 16, 34, 0.22)'
                          : 'rgba(6, 16, 34, 0.48)'
                }
                boxShadow={cellState.selected ? 'app.glowCyan' : 'none'}
                color={cellState.outsideRange ? 'fg.subtle' : 'fg'}
                opacity={cellState.outsideRange ? 0.82 : 1}
                textAlign="left"
                transition="transform 0.2s ease, border-color 0.2s ease, background 0.2s ease"
                style={fixedDayCellStyle}
                _hover={{
                    transform: cellState.disabled ? undefined : 'translateY(-1px)',
                    borderColor: cellState.disabled ? undefined : 'app.sidebarItemHoverBorder',
                    bg: cellState.selected ? 'rgba(32, 208, 255, 0.14)' : 'rgba(8, 20, 42, 0.64)',
                }}
            >
                <Flex justify="space-between" align="start" gap="2">
                    <VStack align="start" gap="0.5">
                        <Text
                            fontSize={{base: 'md', md: 'lg'}}
                            fontWeight={cellState.selected ? '800' : '700'}
                            color={
                                cellState.weekend && !cellState.outsideRange
                                    ? 'app.warning'
                                    : 'currentColor'
                            }
                            lineHeight="1"
                            fontVariantNumeric="tabular-nums"
                        >
                            {date.getDate()}
                        </Text>
                        {cellState.outsideRange ? (
                            <Text fontSize="2xs" color="fg.subtle">
                                {formatOverflowDayLabel(date)}
                            </Text>
                        ) : null}
                    </VStack>

                    <VStack align="end" gap="1">
                        {cellState.today ? (
                            <Box
                                px="2"
                                py="0.5"
                                borderRadius="full"
                                bg="app.accentSoft"
                                color="app.accent"
                            >
                                <Text fontSize="2xs" fontWeight="700">
                                    Сегодня
                                </Text>
                            </Box>
                        ) : null}
                    </VStack>
                </Flex>

                <VStack align="start" gap="1.5">
                    {eventCount > 0 ? (
                        <>
                            <CalendarDayEventDots markers={details?.markers.slice(0, 3) ?? []} />
                            <Text
                                fontSize="xs"
                                fontWeight="600"
                                color="fg"
                                w="full"
                                minW="0"
                                overflow="hidden"
                                whiteSpace="nowrap"
                                textOverflow="ellipsis"
                            >
                                {firstMarker?.label}
                            </Text>
                            <HStack gap="1.5" align="center">
                                {firstMarker && timeChipValue ? (
                                    <CalendarDayTimeChip
                                        ariaLabel={getTimeChipAriaLabel(
                                            firstMarker.label,
                                            timeChipValue,
                                            firstMarker.timeHints.length > 1,
                                        )}
                                        timeLabel={
                                            firstMarker.timeHints.length > 1
                                                ? `${timeChipValue}+`
                                                : timeChipValue
                                        }
                                        color={firstMarker.color}
                                    />
                                ) : null}
                            </HStack>
                        </>
                    ) : (
                        <Text fontSize="xs" color="fg.subtle">
                            Свободно
                        </Text>
                    )}
                </VStack>
            </DatePicker.TableCellTrigger>
        </DatePicker.TableCell>
    );
}
