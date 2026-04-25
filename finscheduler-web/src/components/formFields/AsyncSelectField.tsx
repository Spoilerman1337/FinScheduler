import {
    Combobox,
    createListCollection,
    Field,
    Flex,
    Portal,
    Spinner,
    Text,
} from "@chakra-ui/react";
import {CheckIcon, ChevronDownIcon, XIcon} from "lucide-react";
import {useCallback, useEffect, useMemo, useRef, useState} from "react";
import type {SelectOption} from "./types.ts";

interface AsyncLoadResult {
    options: SelectOption[];
    hasMore: boolean;
}

interface BaseAsyncSelectFieldProps {
    label: string;
    placeholder?: string;
    required?: boolean;
    emptyText?: string;
    initialOptions?: SelectOption[];
    collapseThreshold?: number;
    cacheKey?: string;
    preloadOnMount?: boolean;
    loadOptions: (params: { search: string; page: number }) => Promise<AsyncLoadResult>;
}

type SingleAsyncSelectFieldProps = BaseAsyncSelectFieldProps & {
    multiple?: false;
    value: string;
    onChange: (value: string) => void;
};

type MultipleAsyncSelectFieldProps = BaseAsyncSelectFieldProps & {
    multiple: true;
    value: string[];
    onChange: (value: string[]) => void;
};

type AsyncSelectFieldProps = SingleAsyncSelectFieldProps | MultipleAsyncSelectFieldProps;

interface CacheEntry {
    options: SelectOption[];
    pages: Set<number>;
    hasMore: boolean;
}

const globalSearchCache = new Map<string, CacheEntry>();
const globalRequestCache = new Map<string, Promise<AsyncLoadResult>>();

function mergeOptions(existing: SelectOption[], incoming: SelectOption[]) {
    const map = new Map<string, SelectOption>();

    existing.forEach((option) => map.set(option.value, option));
    incoming.forEach((option) => map.set(option.value, option));

    return Array.from(map.values());
}

function getSearchCacheKey(namespace: string, search: string) {
    return `${namespace}::${search}`;
}

function getRequestCacheKey(namespace: string, search: string, page: number) {
    return `${namespace}::${search}::${page}`;
}

export default function AsyncSelectField(props: AsyncSelectFieldProps) {
    const {
        label,
        placeholder,
        required = false,
        emptyText = 'Ничего не найдено',
        initialOptions = [],
        collapseThreshold = 4,
        cacheKey,
        preloadOnMount = false,
        loadOptions,
    } = props;

    const [inputValue, setInputValue] = useState('');
    const [open, setOpen] = useState(false);
    const [options, setOptions] = useState<SelectOption[]>([]);
    const [currentPage, setCurrentPage] = useState(0);
    const [hasMore, setHasMore] = useState(false);
    const [initialLoading, setInitialLoading] = useState(false);
    const [loadingMore, setLoadingMore] = useState(false);
    const searchRef = useRef('');

    const cacheNamespace = cacheKey ?? label;
    const values = useMemo(
        () => (props.multiple ? props.value : props.value ? [props.value] : []),
        [props.multiple, props.value],
    );
    const mergedOptions = useMemo(
        () => mergeOptions(initialOptions, options),
        [initialOptions, options],
    );
    const selectedOptions = useMemo(
        () => values.map((value) => mergedOptions.find((option) => option.value === value) ?? {
            value,
            label: value,
        }),
        [mergedOptions, values],
    );
    const shouldCollapseSelection = props.multiple && selectedOptions.length > collapseThreshold;
    const collection = useMemo(() => createListCollection({items: mergedOptions}), [mergedOptions]);

    const syncFromCache = useCallback((search: string) => {
        const cachedEntry = globalSearchCache.get(getSearchCacheKey(cacheNamespace, search));

        if (!cachedEntry) {
            return false;
        }

        setOptions(cachedEntry.options);
        setCurrentPage(cachedEntry.pages.size > 0 ? Math.max(...cachedEntry.pages) : 0);
        setHasMore(cachedEntry.hasMore);
        return true;
    }, [cacheNamespace]);

    const loadPage = useCallback(async (search: string, page: number, append: boolean) => {
        const searchCacheKey = getSearchCacheKey(cacheNamespace, search);
        const requestCacheKey = getRequestCacheKey(cacheNamespace, search, page);

        if (append) {
            setLoadingMore(true);
        } else {
            setInitialLoading(true);
        }

        try {
            let request = globalRequestCache.get(requestCacheKey);

            if (!request) {
                request = loadOptions({search, page}).finally(() => {
                    globalRequestCache.delete(requestCacheKey);
                });
                globalRequestCache.set(requestCacheKey, request);
            }

            const result = await request;
            const cachedEntry = globalSearchCache.get(searchCacheKey);
            const previousOptions = append ? cachedEntry?.options ?? [] : [];
            const nextOptions = mergeOptions(previousOptions, result.options);
            const nextPages = new Set(append ? cachedEntry?.pages ?? [] : []);

            nextPages.add(page);

            const nextEntry: CacheEntry = {
                options: nextOptions,
                pages: nextPages,
                hasMore: result.hasMore,
            };

            globalSearchCache.set(searchCacheKey, nextEntry);
            setOptions(nextEntry.options);
            setCurrentPage(page);
            setHasMore(nextEntry.hasMore);
        } finally {
            if (append) {
                setLoadingMore(false);
            } else {
                setInitialLoading(false);
            }
        }
    }, [cacheNamespace, loadOptions]);

    useEffect(() => {
        if (!preloadOnMount) {
            return;
        }

        if (syncFromCache('')) {
            return;
        }

        void loadPage('', 0, false);
    }, [loadPage, preloadOnMount, syncFromCache]);

    useEffect(() => {
        if (!open) {
            return;
        }

        const search = inputValue.trim();
        searchRef.current = search;

        if (syncFromCache(search)) {
            return;
        }

        void loadPage(search, 0, false);
    }, [inputValue, loadPage, open, syncFromCache]);

    const handleValueChange = (nextValues: string[]) => {
        if (props.multiple) {
            props.onChange(nextValues);
        } else {
            props.onChange(nextValues[0] ?? '');
        }
    };

    const handleRemove = (valueToRemove: string) => {
        if (!props.multiple) {
            return;
        }

        props.onChange(props.value.filter((value) => value !== valueToRemove));
    };

    const handleOpenChange = (details: { open: boolean }) => {
        setOpen(details.open);

        if (!details.open) {
            setInputValue('');
        }
    };

    const handleContentScroll = async (event: React.UIEvent<HTMLDivElement>) => {
        if (!hasMore || loadingMore || initialLoading) {
            return;
        }

        const element = event.currentTarget;
        const isNearBottom =
            element.scrollHeight - element.scrollTop - element.clientHeight <= 24;

        if (!isNearBottom) {
            return;
        }

        await loadPage(searchRef.current, currentPage + 1, true);
    };

    return (
        <Field.Root required={required} gap={0}>
            <Field.Label color="neon.blue">
                {label} {(initialLoading || loadingMore) && <Spinner size="xs" ml={2}/>}
            </Field.Label>
            <Combobox.Root
                multiple={props.multiple}
                collection={collection}
                value={values}
                inputValue={inputValue}
                open={open}
                closeOnSelect={!props.multiple}
                openOnClick
                openOnChange
                selectionBehavior={props.multiple ? "clear" : "replace"}
                onOpenChange={handleOpenChange}
                onInputValueChange={(details) => setInputValue(details.inputValue)}
                onValueChange={(details) => handleValueChange(details.value)}
            >
                <Combobox.Control
                    h="40px"
                    px={3}
                    py="0"
                    border="1px solid"
                    borderRadius="sm"
                    borderColor="glass.border"
                    bg="bg.layer2"
                    color="neon.blue"
                    overflow="hidden"
                >
                    <Flex align="center" gap={2} h="100%" flexWrap="nowrap">
                        <Flex
                            align="center"
                            gap="1.5"
                            flex="1"
                            minW="0"
                            overflow="hidden"
                            flexWrap="nowrap"
                        >
                            {props.multiple ? (
                                shouldCollapseSelection ? (
                                    <Text
                                        minW="0"
                                        color="textMuted"
                                        fontSize="sm"
                                        overflow="hidden"
                                        textOverflow="ellipsis"
                                        whiteSpace="nowrap"
                                    >
                                        Выбрано: {selectedOptions.length}
                                    </Text>
                                ) : (
                                    selectedOptions.map((option) => (
                                        <Flex
                                            key={option.value}
                                            align="center"
                                            gap="1"
                                            flexShrink={0}
                                            maxW="140px"
                                            h="20px"
                                            px={2}
                                            borderRadius="sm"
                                            bg="glass.bgHover"
                                            color="neon.blue"
                                            border="1px solid"
                                            borderColor="glass.borderStrong"
                                        >
                                            <Text
                                                fontSize="xs"
                                                lineHeight="1"
                                                overflow="hidden"
                                                textOverflow="ellipsis"
                                                whiteSpace="nowrap"
                                            >
                                                {option.label}
                                            </Text>
                                            <Flex
                                                as="button"
                                                appearance="none"
                                                bg="transparent"
                                                border="none"
                                                outline="none"
                                                p="0"
                                                m="0"
                                                align="center"
                                                justify="center"
                                                w="10px"
                                                h="10px"
                                                flexShrink={0}
                                                color="neon.blue"
                                                opacity={0.8}
                                                transition="opacity 0.2s ease"
                                                _hover={{opacity: 1}}
                                                onPointerDown={(event) => {
                                                    event.preventDefault();
                                                    event.stopPropagation();
                                                }}
                                                onClick={(event) => {
                                                    event.preventDefault();
                                                    event.stopPropagation();
                                                    handleRemove(option.value);
                                                }}
                                            >
                                                <XIcon
                                                    size={8}
                                                    strokeWidth={2}
                                                    style={{
                                                        color: "var(--chakra-colors-neon-blue)",
                                                        stroke: "var(--chakra-colors-neon-blue)",
                                                        display: "block",
                                                    }}
                                                />
                                            </Flex>
                                        </Flex>
                                    ))
                                )
                            ) : (
                                !open &&
                                selectedOptions[0] && (
                                    <Text
                                        minW="0"
                                        color="neon.blue"
                                        fontSize="sm"
                                        overflow="hidden"
                                        textOverflow="ellipsis"
                                        whiteSpace="nowrap"
                                    >
                                        {selectedOptions[0].label}
                                    </Text>
                                )
                            )}
                            <Combobox.Input
                                placeholder={values.length === 0 ? placeholder : ''}
                                flex="1"
                                minW={values.length === 0 ? "0" : "56px"}
                                h="24px"
                                bg="transparent"
                                border="none"
                                outline="none"
                                color="neon.blue"
                                fontSize="sm"
                                _placeholder={{color: 'text.placeholder'}}
                                _focusVisible={{outline: 'none', boxShadow: 'none'}}
                            />
                        </Flex>
                        <Combobox.Trigger
                            appearance="none"
                            bg="transparent"
                            border="none"
                            outline="none"
                            p="0"
                            m="0"
                            display="flex"
                            alignItems="center"
                            justifyContent="center"
                            w="14px"
                            h="14px"
                            flexShrink={0}
                        >
                            <ChevronDownIcon
                                size={14}
                                strokeWidth={1.8}
                                style={{
                                    color: open
                                        ? "var(--chakra-colors-neon-blue)"
                                        : "var(--chakra-colors-textMuted)",
                                    stroke: open
                                        ? "var(--chakra-colors-neon-blue)"
                                        : "var(--chakra-colors-textMuted)",
                                    display: "block",
                                }}
                            />
                        </Combobox.Trigger>
                    </Flex>
                </Combobox.Control>
                <Portal>
                    <Combobox.Positioner>
                        <Combobox.Content
                            p="1"
                            mt="1"
                            border="1px solid"
                            borderRadius="sm"
                            bg="bg.layer1"
                            borderColor="glass.borderStrong"
                            backdropFilter="blur(16px)"
                            boxShadow="lg"
                            zIndex="popover"
                            width="var(--reference-width)"
                            maxH="200px"
                            overflowY="auto"
                            overscrollBehavior="contain"
                            className="custom-scrollbar"
                            onScroll={handleContentScroll}
                        >
                            {mergedOptions.map((option) => (
                                <Combobox.Item
                                    key={option.value}
                                    item={option}
                                    px="3"
                                    py="2"
                                    borderRadius="sm"
                                    color="neon.blue"
                                    display="flex"
                                    alignItems="center"
                                    justifyContent="space-between"
                                    _highlighted={{
                                        color: "neon.purple",
                                        bg: "glass.bgHover",
                                        filter: "drop-shadow(0 0 12px rgba(212, 0, 255, 0.55))",
                                        boxShadow: "0 0 12px rgba(212, 0, 255, 0.35)",
                                    }}
                                    _selected={{
                                        color: "neon.blue",
                                        bg: "glass.bgHover",
                                        fontWeight: "semibold",
                                        filter: "drop-shadow(0 0 10px rgba(0, 212, 255, 0.5))",
                                        boxShadow: "0 0 10px rgba(0, 212, 255, 0.3)",
                                    }}
                                >
                                    <Combobox.ItemText>{option.label}</Combobox.ItemText>
                                    <Combobox.ItemIndicator>
                                        <CheckIcon size={14}/>
                                    </Combobox.ItemIndicator>
                                </Combobox.Item>
                            ))}
                            {!initialLoading && mergedOptions.length === 0 && (
                                <Combobox.Empty px="3" py="2" color="textMuted">
                                    {emptyText}
                                </Combobox.Empty>
                            )}
                            {loadingMore && (
                                <Flex align="center" justify="center" py="2">
                                    <Spinner size="sm"/>
                                </Flex>
                            )}
                            {initialLoading && mergedOptions.length === 0 && (
                                <Flex align="center" justify="center" py="3">
                                    <Spinner size="sm"/>
                                </Flex>
                            )}
                        </Combobox.Content>
                    </Combobox.Positioner>
                </Portal>
            </Combobox.Root>
        </Field.Root>
    );
}
