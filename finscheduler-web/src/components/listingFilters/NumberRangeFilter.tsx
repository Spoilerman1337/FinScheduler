import {Button, Flex, NumberInput, Popover, Text} from '@chakra-ui/react';
import {ChevronDown} from 'lucide-react';
import {useEffect, useMemo, useState} from 'react';
import {filterWidthProps} from './shared.ts';

export type NumberRangeValue = {
    from: string;
    to: string;
};

type NumberRangeFilterProps = {
    label: string;
    value: NumberRangeValue;
    onChange: (value: NumberRangeValue) => void;
    description?: string;
    suffix?: string;
    min?: number;
    formatValue?: (value: number) => string;
    fromPlaceholder?: string;
    toPlaceholder?: string;
    fromInputLabel?: string;
    toInputLabel?: string;
    clearLabel?: string;
    applyLabel?: string;
};

function sanitizeNumberValue(value: string, min?: number): string {
    const normalized = value.replace(',', '.').trim();

    if (!normalized) {
        return '';
    }

    const parsed = Number(normalized);

    if (!Number.isFinite(parsed)) {
        return '';
    }

    if (min !== undefined && parsed < min) {
        return '';
    }

    return normalized;
}

function normalizeNumberRange(value: NumberRangeValue, min?: number): NumberRangeValue {
    const from = sanitizeNumberValue(value.from, min);
    const to = sanitizeNumberValue(value.to, min);

    if (!from || !to) {
        return {from, to};
    }

    return parseFloat(from) <= parseFloat(to)
        ? {from, to}
        : {from: to, to: from};
}

function formatNumberValue(
    value: string,
    formatValue?: (value: number) => string,
): string | null {
    if (!value) {
        return null;
    }

    const parsed = Number(value);

    if (!Number.isFinite(parsed)) {
        return null;
    }

    if (formatValue) {
        return formatValue(parsed);
    }

    return parsed.toLocaleString('ru-RU', {maximumFractionDigits: 2});
}

function withSuffix(value: string, suffix?: string): string {
    return suffix ? `${value} ${suffix}` : value;
}

function getNumberRangeSummary(props: {
    label: string;
    value: NumberRangeValue;
    suffix?: string;
    formatValue?: (value: number) => string;
}): string {
    const from = formatNumberValue(props.value.from, props.formatValue);
    const to = formatNumberValue(props.value.to, props.formatValue);

    if (from && to) {
        return `${props.label}: ${from} - ${withSuffix(to, props.suffix)}`;
    }

    if (from) {
        return `${props.label} от ${withSuffix(from, props.suffix)}`;
    }

    if (to) {
        return `${props.label} до ${withSuffix(to, props.suffix)}`;
    }

    return props.label;
}

export default function NumberRangeFilter(props: NumberRangeFilterProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [draftValue, setDraftValue] = useState<NumberRangeValue>(props.value);

    const fromInputLabel = props.fromInputLabel ?? `${props.label} от`;
    const toInputLabel = props.toInputLabel ?? `${props.label} до`;
    const fromPlaceholder = props.fromPlaceholder ?? 'От';
    const toPlaceholder = props.toPlaceholder ?? 'До';
    const clearLabel = props.clearLabel ?? 'Очистить';
    const applyLabel = props.applyLabel ?? 'Применить';

    useEffect(() => {
        setDraftValue(props.value);
    }, [props.value]);

    const summary = useMemo(
        () =>
            getNumberRangeSummary({
                label: props.label,
                value: props.value,
                suffix: props.suffix,
                formatValue: props.formatValue,
            }),
        [props.formatValue, props.label, props.suffix, props.value],
    );
    const isActive = Boolean(props.value.from || props.value.to);

    const handleApply = () => {
        props.onChange(normalizeNumberRange(draftValue, props.min));
        setIsOpen(false);
    };

    const handleResetDraft = () => {
        setDraftValue({from: '', to: ''});
    };

    const handleOpenChange = (details: {open: boolean}) => {
        setIsOpen(details.open);

        if (!details.open) {
            setDraftValue(props.value);
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

            <Popover.Positioner>
                <Popover.Content
                    bg="bg.layer1"
                    border="1px solid"
                    borderColor="glass.border"
                    boxShadow="0 0 24px rgba(0, 0, 0, 0.35)"
                    maxW="360px"
                    zIndex="popover"
                >
                    <Popover.Arrow>
                        <Popover.ArrowTip />
                    </Popover.Arrow>

                    <Popover.Body p={4}>
                        <Flex direction="column" gap={4}>
                            <Flex justify="space-between" align="center" gap={3}>
                                <Flex direction="column" gap={1}>
                                    <Text color="neon.blue" fontWeight="semibold">
                                        {props.label}
                                    </Text>
                                    {props.description ? (
                                        <Text color="textMuted" fontSize="sm">
                                            {props.description}
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
                                    {clearLabel}
                                </Button>
                            </Flex>

                            <Flex gap={3} direction={{base: 'column', sm: 'row'}}>
                                <NumberInput.Root
                                    value={draftValue.from}
                                    onValueChange={(details) => {
                                        setDraftValue((current) => ({...current, from: details.value}));
                                    }}
                                    min={props.min}
                                    width="100%"
                                    gap={0}
                                >
                                    <NumberInput.Control>
                                        <NumberInput.IncrementTrigger
                                            bg="bg.layer1"
                                            color="neon.blue"
                                            p={0}
                                        >
                                            +
                                        </NumberInput.IncrementTrigger>
                                        <NumberInput.DecrementTrigger
                                            bg="bg.layer1"
                                            color="neon.blue"
                                            p={0}
                                        >
                                            -
                                        </NumberInput.DecrementTrigger>
                                    </NumberInput.Control>
                                    <NumberInput.Input
                                        aria-label={fromInputLabel}
                                        placeholder={fromPlaceholder}
                                        bg="bg.layer1"
                                        borderColor="glass.border"
                                        color="neon.blue"
                                        _placeholder={{color: 'text.placeholder'}}
                                        onKeyDown={(event) => {
                                            if (event.key === 'Enter') {
                                                handleApply();
                                            }
                                        }}
                                    />
                                </NumberInput.Root>

                                <NumberInput.Root
                                    value={draftValue.to}
                                    onValueChange={(details) => {
                                        setDraftValue((current) => ({...current, to: details.value}));
                                    }}
                                    min={props.min}
                                    width="100%"
                                    gap={0}
                                >
                                    <NumberInput.Control>
                                        <NumberInput.IncrementTrigger
                                            bg="bg.layer1"
                                            color="neon.blue"
                                            p={0}
                                        >
                                            +
                                        </NumberInput.IncrementTrigger>
                                        <NumberInput.DecrementTrigger
                                            bg="bg.layer1"
                                            color="neon.blue"
                                            p={0}
                                        >
                                            -
                                        </NumberInput.DecrementTrigger>
                                    </NumberInput.Control>
                                    <NumberInput.Input
                                        aria-label={toInputLabel}
                                        placeholder={toPlaceholder}
                                        bg="bg.layer1"
                                        borderColor="glass.border"
                                        color="neon.blue"
                                        _placeholder={{color: 'text.placeholder'}}
                                        onKeyDown={(event) => {
                                            if (event.key === 'Enter') {
                                                handleApply();
                                            }
                                        }}
                                    />
                                </NumberInput.Root>
                            </Flex>

                            <Button
                                bg="neon.blue"
                                color="bg.base"
                                _hover={{bg: 'neon.blue', opacity: 0.85}}
                                onClick={handleApply}
                            >
                                {applyLabel}
                            </Button>
                        </Flex>
                    </Popover.Body>
                </Popover.Content>
            </Popover.Positioner>
        </Popover.Root>
    );
}
