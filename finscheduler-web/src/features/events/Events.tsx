import {Badge, Box, Flex, HStack, Text, VStack} from '@chakra-ui/react';
import {useMemo, useState} from 'react';
import DateRangeFilter, {
    type DateRangeFilterMode,
    type DateRangeFilterValue,
} from '../../components/listingFilters/DateRangeFilter.tsx';
import TextInputFilter from '../../components/listingFilters/TextInputFilter.tsx';
import FilterWrapper from '../../components/listingFilters/FilterWrapper.tsx';
import DataTable, {type DataListingColumn} from '../../components/dataTable/DataTable.tsx';
import {calendarEventColorStyles} from '../calendar/calendarUi.ts';
import {buildEventCatalogEntries, startOfDay, type EventCatalogEntry} from '../calendar/shared.ts';

type TimelineState = 'today' | 'upcoming' | 'past';
type EventDateFilterMode = 'eventDate';

const eventDateModes: readonly [
    DateRangeFilterMode<EventDateFilterMode>,
    ...DateRangeFilterMode<EventDateFilterMode>[],
] = [
    {
        value: 'eventDate',
        label: 'Дата',
        description: 'Показывать события, которые попадают в выбранный период.',
        fromInputLabel: 'Дата от',
        toInputLabel: 'Дата до',
    },
];

function normalizeSearchValue(value: string) {
    return value.trim().toLocaleLowerCase('ru-RU');
}

function startOfYear(date: Date) {
    return new Date(date.getFullYear(), 0, 1);
}

function getDayTimestamp(date: Date) {
    return startOfDay(date).getTime();
}

function getTimelineState(date: Date, today: Date): TimelineState {
    const eventTimestamp = getDayTimestamp(date);
    const todayTimestamp = getDayTimestamp(today);

    if (eventTimestamp === todayTimestamp) {
        return 'today';
    }

    return eventTimestamp > todayTimestamp ? 'upcoming' : 'past';
}

function toTimeHintSortValue(timeHint: string) {
    const [hours = '0', minutes = '0'] = timeHint.split(':');

    return Number(hours) * 60 + Number(minutes);
}

function getEarliestTimeSortValue(event: EventCatalogEntry) {
    if (event.timeHints.length === 0) {
        return Number.POSITIVE_INFINITY;
    }

    return event.timeHints.reduce((smallest, timeHint) => {
        return Math.min(smallest, toTimeHintSortValue(timeHint));
    }, Number.POSITIVE_INFINITY);
}

function compareEventEntries(left: EventCatalogEntry, right: EventCatalogEntry, today: Date) {
    const leftState = getTimelineState(left.date, today);
    const rightState = getTimelineState(right.date, today);
    const timelineRank: Record<TimelineState, number> = {
        today: 0,
        upcoming: 1,
        past: 2,
    };

    if (timelineRank[leftState] !== timelineRank[rightState]) {
        return timelineRank[leftState] - timelineRank[rightState];
    }

    const leftTimestamp = getDayTimestamp(left.date);
    const rightTimestamp = getDayTimestamp(right.date);

    if (leftState === 'past') {
        if (leftTimestamp !== rightTimestamp) {
            return rightTimestamp - leftTimestamp;
        }
    } else if (leftTimestamp !== rightTimestamp) {
        return leftTimestamp - rightTimestamp;
    }

    return getEarliestTimeSortValue(left) - getEarliestTimeSortValue(right);
}

function createDefaultEventDateFilter(): DateRangeFilterValue<EventDateFilterMode> {
    return {
        mode: 'eventDate',
        from: '',
        to: '',
    };
}

function buildYearEventCatalogEntries(referenceDate: Date) {
    const yearStart = startOfYear(referenceDate);

    return Array.from({length: 12}, (_, monthIndex) => {
        return new Date(yearStart.getFullYear(), monthIndex, 1);
    }).flatMap((monthStart) => buildEventCatalogEntries(monthStart, referenceDate));
}

function EventTimeBadges(props: {event: EventCatalogEntry}) {
    const {event} = props;
    const colorStyles = calendarEventColorStyles[event.color];

    if (event.timeHints.length === 0) {
        return <Text color="fg.subtle">Без времени</Text>;
    }

    return (
        <Flex wrap="wrap" gap="2">
            {event.timeHints.map((timeHint) => (
                <Badge
                    key={`${event.id}-${timeHint}`}
                    px="2.5"
                    py="1"
                    borderRadius="full"
                    borderWidth="1px"
                    borderColor="app.cardBorder"
                    bg="rgba(6, 16, 34, 0.24)"
                    color={colorStyles.text}
                >
                    {timeHint}
                </Badge>
            ))}
        </Flex>
    );
}

export default function Events() {
    const [today] = useState(() => startOfDay(new Date()));
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [dateFilter, setDateFilter] = useState<DateRangeFilterValue<EventDateFilterMode>>(() =>
        createDefaultEventDateFilter(),
    );

    const allEvents = useMemo(() => {
        return buildYearEventCatalogEntries(today).sort((left, right) => {
            return compareEventEntries(left, right, today);
        });
    }, [today]);

    const filteredEvents = useMemo(() => {
        const normalizedSearchTerm = normalizeSearchValue(searchTerm);

        return allEvents.filter((event) => {
            if (dateFilter.from && event.dayKey < dateFilter.from) {
                return false;
            }

            if (dateFilter.to && event.dayKey > dateFilter.to) {
                return false;
            }

            if (!normalizedSearchTerm) {
                return true;
            }

            return normalizeSearchValue(event.label).includes(normalizedSearchTerm);
        });
    }, [allEvents, dateFilter.from, dateFilter.to, searchTerm]);

    const pagedEvents = useMemo(() => {
        const startIndex = (page - 1) * pageSize;

        return filteredEvents.slice(startIndex, startIndex + pageSize);
    }, [filteredEvents, page, pageSize]);

    const eventColumns: DataListingColumn<EventCatalogEntry>[] = [
        {
            header: 'Событие',
            key: 'label',
            render: (row: EventCatalogEntry) => (
                <VStack align="start" gap="1" minW="0">
                    <HStack gap="2.5" align="start">
                        <Box
                            mt="1.5"
                            boxSize="2.5"
                            borderRadius="full"
                            bg={calendarEventColorStyles[row.color].bg}
                            boxShadow={calendarEventColorStyles[row.color].glow}
                            flexShrink={0}
                        />
                        <VStack align="start" gap="1" minW="0">
                            <Text fontWeight="semibold" color="neon.blue">
                                {row.label}
                            </Text>
                            <Text color="fg.muted" fontSize="sm">
                                {row.description}
                            </Text>
                        </VStack>
                    </HStack>
                </VStack>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Дата',
            key: 'date',
            render: (row: EventCatalogEntry) => (
                <Text color="neon.blue" fontWeight="medium">
                    {row.fullDateLabel}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Время',
            key: 'timeHints',
            render: (row: EventCatalogEntry) => <EventTimeBadges event={row} />,
            headerProps: {textAlign: 'left'},
        },
    ];

    const handleReset = () => {
        setSearchTerm('');
        setDateFilter(createDefaultEventDateFilter());
        setPage(1);
    };

    return (
        <Flex direction="column" width="100%">
            <FilterWrapper onReset={handleReset}>
                <TextInputFilter
                    value={searchTerm}
                    placeholder="Поиск по наименованию..."
                    onChange={(value) => {
                        setSearchTerm(value);
                        setPage(1);
                    }}
                    onApply={() => {
                        setPage(1);
                    }}
                />
                <DateRangeFilter
                    label="Дата"
                    modes={eventDateModes}
                    value={dateFilter}
                    onChange={(value) => {
                        setDateFilter({
                            mode: value.mode as EventDateFilterMode,
                            from: value.from,
                            to: value.to,
                        });
                        setPage(1);
                    }}
                />
            </FilterWrapper>

            <DataTable
                data={pagedEvents}
                columns={eventColumns}
                total={filteredEvents.length}
                page={page}
                pageSize={pageSize}
                onPageChange={setPage}
                onPageSizeChange={(newSize) => {
                    setPageSize(newSize);
                    setPage(1);
                }}
                selectable={false}
                getRowId={(row) => row.id}
            />
        </Flex>
    );
}