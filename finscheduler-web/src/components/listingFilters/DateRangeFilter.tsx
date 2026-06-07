import {Button, DatePicker, Flex, Popover, Portal, Text, parseDate} from '@chakra-ui/react';
import {ChevronDown} from 'lucide-react';
import {useEffect, useMemo, useState} from 'react';
import {filterWidthProps} from './shared.ts';

export type DateRangeFilterValue<TMode extends string = string> = {
    mode: TMode;
    from: string;
    to: string;
};

export type DateRangeFilterMode<TMode extends string = string> = {
    value: TMode;
    label: string;
    summaryLabel?: string;
    description?: string;
    fromInputLabel?: string;
    toInputLabel?: string;
};

type DateRangeFilterProps = {
    label: string;
    modes: readonly [DateRangeFilterMode, ...DateRangeFilterMode[]];
    value: DateRangeFilterValue;
    onChange: (value: DateRangeFilterValue) => void;
    description?: string;
    fromInputLabel?: string;
    toInputLabel?: string;
};

function normalizeIsoDateValue(value: string): string {
    const normalized = value.trim();

    if (!normalized) {
        return '';
    }

    const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(normalized);

    if (!match) {
        return '';
    }

    const year = Number(match[1]);
    const month = Number(match[2]);
    const day = Number(match[3]);
    const parsed = new Date(Date.UTC(year, month - 1, day));

    const isValidDate =
        parsed.getUTCFullYear() === year &&
        parsed.getUTCMonth() === month - 1 &&
        parsed.getUTCDate() === day;

    return isValidDate ? normalized : '';
}

function getActiveMode(
    modes: readonly [DateRangeFilterMode, ...DateRangeFilterMode[]],
    modeValue: string,
): DateRangeFilterMode {
    return modes.find((mode) => mode.value === modeValue) ?? modes[0];
}

function normalizeDateRangeValue(
    value: DateRangeFilterValue,
    modes: readonly [DateRangeFilterMode, ...DateRangeFilterMode[]],
): DateRangeFilterValue {
    const activeMode = getActiveMode(modes, value.mode);
    const from = normalizeIsoDateValue(value.from);
    const to = normalizeIsoDateValue(value.to);

    if (!from || !to) {
        return {
            mode: activeMode.value,
            from,
            to,
        };
    }

    return from <= to
        ? {
              mode: activeMode.value,
              from,
              to,
          }
        : {
              mode: activeMode.value,
              from: to,
              to: from,
          };
}

function formatDateValue(value: string): string | null {
    const normalized = normalizeIsoDateValue(value);

    if (!normalized) {
        return null;
    }

    const [year, month, day] = normalized.split('-');
    return `${day}.${month}.${year}`;
}

function getModeSummaryLabel(mode: DateRangeFilterMode): string {
    return mode.summaryLabel ?? mode.label;
}

function getDateRangeSummary(mode: DateRangeFilterMode, value: DateRangeFilterValue): string {
    const label = getModeSummaryLabel(mode);
    const from = formatDateValue(value.from);
    const to = formatDateValue(value.to);

    if (from && to) {
        return `${label}: ${from} - ${to}`;
    }

    if (from) {
        return `${label} от ${from}`;
    }

    if (to) {
        return `${label} до ${to}`;
    }

    return label;
}

export default function DateRangeFilter(props: DateRangeFilterProps) {
    const normalizedValue = useMemo(
        () => normalizeDateRangeValue(props.value, props.modes),
        [props.modes, props.value],
    );
    const [isOpen, setIsOpen] = useState(false);
    const [draftValue, setDraftValue] = useState<DateRangeFilterValue>(normalizedValue);

    useEffect(() => {
        setDraftValue(normalizedValue);
    }, [normalizedValue]);

    const activeMode = useMemo(
        () => getActiveMode(props.modes, draftValue.mode),
        [props.modes, draftValue.mode],
    );
    const appliedMode = useMemo(
        () => getActiveMode(props.modes, props.value.mode),
        [props.modes, props.value.mode],
    );

    const pickerValue = useMemo(() => {
        return [draftValue.from, draftValue.to]
            .map(normalizeIsoDateValue)
            .filter(Boolean)
            .map((value) => parseDate(value));
    }, [draftValue.from, draftValue.to]);

    const summary = useMemo(
        () => getDateRangeSummary(appliedMode, props.value),
        [appliedMode, props.value],
    );
    const isActive = Boolean(props.value.from || props.value.to);
    const showModeSwitch = props.modes.length > 1;

    const modeSummaryLabel = getModeSummaryLabel(activeMode);
    const fromInputLabel =
        activeMode.fromInputLabel ?? props.fromInputLabel ?? `${modeSummaryLabel} от`;
    const toInputLabel = activeMode.toInputLabel ?? props.toInputLabel ?? `${modeSummaryLabel} до`;
    const description = activeMode.description ?? props.description;

    const handleApply = () => {
        props.onChange(normalizeDateRangeValue(draftValue, props.modes));
        setIsOpen(false);
    };

    const handleResetDraft = () => {
        setDraftValue((current) => ({...current, from: '', to: ''}));
    };

    const handleOpenChange = (details: {open: boolean}) => {
        setIsOpen(details.open);

        if (!details.open) {
            setDraftValue(normalizedValue);
        }
    };

    return (
        <Popover.Root
            open={isOpen}
            onOpenChange={handleOpenChange}
            positioning={{placement: 'bottom-start'}}
        >
            <Popover.Trigger asChild>
                <Button
                    {...filterWidthProps}
                    aria-label={summary}
                    justifyContent="space-between"
                    bg="bg.layer1"
                    color="neon.blue"
                    border="1px solid"
                    borderColor={isOpen || isActive ? 'neon.blue' : 'glass.border'}
                    backdropFilter="blur(12px)"
                    transition="all 0.3s ease-in-out"
                    _hover={{
                        borderColor: 'neon.blue',
                        color: 'neon.blue',
                        filter: 'drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))',
                        boxShadow: '0 0 12px rgba(0, 212, 255, 0.45)',
                    }}
                    focusRing="none"
                >
                    <Text
                        flex="1"
                        minW={0}
                        textAlign="left"
                        overflow="hidden"
                        textOverflow="ellipsis"
                        whiteSpace="nowrap"
                    >
                        {summary}
                    </Text>
                    <ChevronDown size={16} />
                </Button>
            </Popover.Trigger>

            <Portal>
                <Popover.Positioner>
                    <Popover.Content
                        bg="bg.layer1"
                        border="1px solid"
                        borderColor="glass.border"
                        boxShadow="0 0 24px rgba(0, 0, 0, 0.35)"
                        width={{base: 'min(calc(100vw - 1rem), 420px)', sm: '420px'}}
                        maxW="calc(100vw - 1rem)"
                        zIndex="popover"
                    >
                        <Popover.Arrow>
                            <Popover.ArrowTip />
                        </Popover.Arrow>

                        <Popover.Body p={4}>
                            <DatePicker.Root
                                inline
                                open
                                selectionMode="range"
                                closeOnSelect={false}
                                locale="ru-RU"
                                timeZone="UTC"
                                value={pickerValue}
                                onValueChange={(details) => {
                                    setDraftValue((current) => ({
                                        ...current,
                                        from: details.value[0]?.toString() ?? '',
                                        to: details.value[1]?.toString() ?? '',
                                    }));
                                }}
                            >
                                <Flex direction="column" gap={4}>
                                    <Flex justify="space-between" align="center" gap={3}>
                                        <Flex direction="column" gap={1}>
                                            <Text color="neon.blue" fontWeight="semibold">
                                                {props.label}
                                            </Text>
                                            {description ? (
                                                <Text color="textMuted" fontSize="sm">
                                                    {description}
                                                </Text>
                                            ) : null}
                                        </Flex>

                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            color="textMuted"
                                            _hover={{bg: 'bg.layer2', color: 'neon.blue'}}
                                            onClick={handleResetDraft}
                                        >
                                            Очистить
                                        </Button>
                                    </Flex>

                                    {showModeSwitch ? (
                                        <Flex gap={2} wrap="wrap">
                                            {props.modes.map((mode) => {
                                                const isSelected = draftValue.mode === mode.value;

                                                return (
                                                    <Button
                                                        key={mode.value}
                                                        size="sm"
                                                        aria-pressed={isSelected}
                                                        bg={isSelected ? 'neon.blue' : 'bg.layer1'}
                                                        color={isSelected ? 'bg.base' : 'textMuted'}
                                                        border="1px solid"
                                                        borderColor={
                                                            isSelected ? 'neon.blue' : 'glass.border'
                                                        }
                                                        _hover={{
                                                            bg: isSelected ? 'neon.blue' : 'bg.layer2',
                                                            color: isSelected ? 'bg.base' : 'neon.blue',
                                                            borderColor: 'neon.blue',
                                                        }}
                                                        onClick={() => {
                                                            setDraftValue((current) => ({
                                                                ...current,
                                                                mode: mode.value,
                                                            }));
                                                        }}
                                                    >
                                                        {mode.label}
                                                    </Button>
                                                );
                                            })}
                                        </Flex>
                                    ) : null}

                                    <DatePicker.Control
                                        display="grid"
                                        gridTemplateColumns={{base: '1fr', sm: '1fr 1fr'}}
                                        gap={3}
                                    >
                                        <Flex direction="column" gap={2}>
                                            <Text color="textMuted" fontSize="sm">
                                                С
                                            </Text>
                                            <DatePicker.Input
                                                index={0}
                                                aria-label={fromInputLabel}
                                                placeholder="дд.мм.гггг"
                                                bg="bg.layer1"
                                                border="1px solid"
                                                borderColor="glass.border"
                                                borderRadius="md"
                                                color="neon.blue"
                                                px={3}
                                                _placeholder={{color: 'text.placeholder'}}
                                            />
                                        </Flex>

                                        <Flex direction="column" gap={2}>
                                            <Text color="textMuted" fontSize="sm">
                                                По
                                            </Text>
                                            <DatePicker.Input
                                                index={1}
                                                aria-label={toInputLabel}
                                                placeholder="дд.мм.гггг"
                                                bg="bg.layer1"
                                                border="1px solid"
                                                borderColor="glass.border"
                                                borderRadius="md"
                                                color="neon.blue"
                                                px={3}
                                                _placeholder={{color: 'text.placeholder'}}
                                            />
                                        </Flex>
                                    </DatePicker.Control>

                                    <DatePicker.Content
                                        width="100%"
                                        minW="0"
                                        p={0}
                                        bg="transparent"
                                        boxShadow="none"
                                        borderRadius="none"
                                        display="flex"
                                        flexDirection="column"
                                        gap={4}
                                        overflow="visible"
                                    >
                                        <Flex
                                            width="100%"
                                            border="1px solid"
                                            borderColor="glass.border"
                                            borderRadius="lg"
                                            bg="bg.layer1"
                                            p={3}
                                        >
                                            <DatePicker.View view="day" width="100%">
                                                <DatePicker.Header mb={3} />
                                                <DatePicker.DayTable width="100%" />
                                            </DatePicker.View>
                                            <DatePicker.View view="month" width="100%">
                                                <DatePicker.Header mb={3} />
                                                <DatePicker.MonthTable width="100%" />
                                            </DatePicker.View>
                                            <DatePicker.View view="year" width="100%">
                                                <DatePicker.Header mb={3} />
                                                <DatePicker.YearTable width="100%" />
                                            </DatePicker.View>
                                        </Flex>

                                        <Button
                                            bg="neon.blue"
                                            color="bg.base"
                                            _hover={{bg: 'neon.blue', opacity: 0.85}}
                                            onClick={handleApply}
                                        >
                                            Применить
                                        </Button>
                                    </DatePicker.Content>
                                </Flex>
                            </DatePicker.Root>
                        </Popover.Body>
                    </Popover.Content>
                </Popover.Positioner>
            </Portal>
        </Popover.Root>
    );
}
