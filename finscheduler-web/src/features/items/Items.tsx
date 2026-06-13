import DataShowcase, {type DataListingColumn} from '../../components/dataShowcase/DataShowcase.tsx';
import {Badge, Flex, Spinner, Text} from '@chakra-ui/react';
import {useCallback, useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import type {
    ItemDateFilterValue,
    ItemFilter,
    ItemListingDto,
    ItemStatusFilter,
} from '../../api/items.types.ts';
import ItemsService, {buildItemFilter} from '../../api/items.ts';
import type {NumberRangeValue} from '../../components/listingFilters/NumberRangeFilter.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import ListingBulkCashbackButton from '../../components/dataListing/actionButtons/ListingBulkCashbackButton.tsx';
import {buildEditItemPath, newItemPath} from '../routes.ts';
import {createDefaultItemDateFilter} from './types.ts';
import BulkCashbackModal, {
    type BulkCashbackSubmitPayload,
} from './subcomponents/BulkCashbackModal.tsx';
import ItemsFilters from './subcomponents/ItemsFilters.tsx';
import {getCashbackColor} from './types.ts';

const itemsService = new ItemsService();

export default function Items() {
    const navigate = useNavigate();
    const [items, setItems] = useState<ItemListingDto[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(12);
    const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
    const [selectedItemLabels, setSelectedItemLabels] = useState<Record<string, string>>({});
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<ItemStatusFilter>('Active');
    const [dateFilter, setDateFilter] = useState<ItemDateFilterValue>(createDefaultItemDateFilter);
    const [priceRange, setPriceRange] = useState<NumberRangeValue>({from: '', to: ''});
    const [cashbackRange, setCashbackRange] = useState<NumberRangeValue>({from: '', to: ''});
    const [isBulkCashbackModalOpen, setIsBulkCashbackModalOpen] = useState(false);
    const [bulkCashbackLoading, setBulkCashbackLoading] = useState(false);
    const [bulkCashbackError, setBulkCashbackError] = useState<string | null>(null);

    const itemColumns: DataListingColumn<ItemListingDto>[] = [
        {
            header: 'Название',
            key: 'name',
            render: (row: ItemListingDto) => (
                <Text fontWeight="semibold" color="neon.blue">
                    {row.name || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Цена',
            key: 'price',
            render: (row: ItemListingDto) => (
                <Text color="neon.blue" fontWeight="medium">
                    {row.price !== undefined
                        ? `${row.price.toLocaleString('ru-RU', {
                              minimumFractionDigits: 2,
                              maximumFractionDigits: 2,
                          })} ₽`
                        : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-start'},
        },
        {
            header: 'Кэшбэк (%)',
            key: 'cashback',
            render: (row: ItemListingDto) => (
                <Text color={getCashbackColor(row.cashback)} fontWeight="bold">
                    {row.cashback !== undefined ? `${row.cashback}%` : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-start'},
        },
        {
            header: 'Статус',
            key: 'isActive',
            render: (row: ItemListingDto) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.isActive ? 'neon.green' : 'neon.pink'}
                    color="bg.base"
                >
                    {row.isActive ? 'Активен' : 'Неактивен'}
                </Badge>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Обновлён',
            key: 'updatedAt',
            render: (row: ItemListingDto) => (
                <Text color="neon.blue" fontSize="sm">
                    {row.updatedAt ? new Date(row.updatedAt).toLocaleDateString('ru-RU') : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
    ];

    const loadItems = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const filter: ItemFilter = buildItemFilter({
                page,
                pageSize,
                searchTerm,
                statusFilter,
                dateFilter,
                priceFrom: priceRange.from,
                priceTo: priceRange.to,
                cashbackFrom: cashbackRange.from,
                cashbackTo: cashbackRange.to,
            });
            const result = await itemsService.getListingInfo(filter);
            setItems(result.data);
            setTotal(result.count);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка загрузки данных');
            console.error('Failed to load items:', err);
        } finally {
            setLoading(false);
        }
    }, [
        page,
        pageSize,
        searchTerm,
        statusFilter,
        dateFilter,
        priceRange.from,
        priceRange.to,
        cashbackRange.from,
        cashbackRange.to,
    ]);

    useEffect(() => {
        loadItems();
    }, [loadItems]);

    const handleReset = () => {
        setSearchTerm('');
        setStatusFilter('Active');
        setDateFilter(createDefaultItemDateFilter());
        setPriceRange({from: '', to: ''});
        setCashbackRange({from: '', to: ''});
        setPage(1);
    };

    const handleSelectionChange = useCallback(
        (nextSelectedRows: Set<string>) => {
            setSelectedRows(nextSelectedRows);
            setSelectedItemLabels((prev) => {
                const next: Record<string, string> = {};

                nextSelectedRows.forEach((id) => {
                    const currentItem = items.find((item) => item.id === id);

                    if (currentItem?.name) {
                        next[id] = currentItem.name;
                        return;
                    }

                    if (prev[id]) {
                        next[id] = prev[id];
                    }
                });

                const prevKeys = Object.keys(prev);
                const nextKeys = Object.keys(next);

                if (
                    prevKeys.length === nextKeys.length &&
                    nextKeys.every((key) => prev[key] === next[key])
                ) {
                    return prev;
                }

                return next;
            });
        },
        [items],
    );

    const handleOpenCreatePage = () => {
        navigate(newItemPath);
    };

    const handleOpenBulkCashbackModal = () => {
        setBulkCashbackError(null);
        setIsBulkCashbackModalOpen(true);
    };

    const handleCloseBulkCashbackModal = () => {
        if (bulkCashbackLoading) {
            return;
        }

        setBulkCashbackError(null);
        setIsBulkCashbackModalOpen(false);
    };

    const handleOpenEditPage = (item: ItemListingDto) => {
        if (!item.id) {
            toaster.create({
                title: 'Ошибка',
                description: 'Не удалось открыть карточку предмета',
                type: 'error',
            });
            return;
        }

        navigate(buildEditItemPath(item.id));
    };

    const handleDeleteItems = async (ids: string[]) => {
        try {
            await Promise.all(ids.map((id) => itemsService.deleteItem(id)));
            toaster.create({
                title: 'Успешно',
                description: `Удалено элементов: ${ids.length}`,
                type: 'success',
            });
            setSelectedRows(new Set());
            setSelectedItemLabels({});
            await loadItems();
        } catch (err) {
            toaster.create({
                title: 'Ошибка',
                description: err instanceof Error ? err.message : 'Не удалось удалить элементы',
                type: 'error',
            });
        }
    };

    const handleBulkCashbackSubmit = async (payload: BulkCashbackSubmitPayload) => {
        setBulkCashbackLoading(true);
        setBulkCashbackError(null);

        try {
            let successDescription = 'Кешбек успешно обновлен';

            if (payload.itemIds && payload.itemIds.length > 0) {
                await itemsService.updateCashbackByItems(payload.itemIds, payload.cashback);
                successDescription = `Кешбек обновлен у выбранных элементов: ${payload.itemIds.length}`;
                setSelectedRows(new Set());
                setSelectedItemLabels({});
            } else if (payload.tagId) {
                await itemsService.updateCashbackByTag(payload.tagId, payload.cashback);
                successDescription = 'Кешбек обновлен у элементов выбранного тега';
            } else {
                setBulkCashbackError('Не удалось определить сценарий массового обновления');
                return;
            }

            setIsBulkCashbackModalOpen(false);
            await loadItems();
            toaster.create({
                title: 'Успешно',
                description: successDescription,
                type: 'success',
            });
        } catch (err) {
            setBulkCashbackError(
                err instanceof Error ? err.message : 'Не удалось массово обновить кешбек',
            );
        } finally {
            setBulkCashbackLoading(false);
        }
    };

    const getRowId = (row: ItemListingDto): string => {
        return row.id ?? '';
    };

    const selectedItems = Array.from(selectedRows).map((id) => ({
        id,
        label: selectedItemLabels[id] ?? id,
    }));

    return (
        <Flex direction="column" width="100%">
            <ItemsFilters
                searchTerm={searchTerm}
                statusFilter={statusFilter}
                dateFilter={dateFilter}
                priceRange={priceRange}
                cashbackRange={cashbackRange}
                onSearchTermChange={(value) => {
                    setPage(1);
                    setSearchTerm(value);
                }}
                onStatusFilterChange={(value) => {
                    setStatusFilter(value);
                    setPage(1);
                }}
                onDateFilterChange={(value) => {
                    setDateFilter(value);
                    setPage(1);
                }}
                onPriceRangeChange={(value) => {
                    setPriceRange(value);
                    setPage(1);
                }}
                onCashbackRangeChange={(value) => {
                    setCashbackRange(value);
                    setPage(1);
                }}
                onApply={() => {
                    setPage(1);
                    loadItems();
                }}
                onReset={handleReset}
            />

            {loading ? (
                <Flex justify="center" align="center" minH="400px">
                    <Spinner size="xl" color="neon.blue" />
                </Flex>
            ) : error ? (
                <Flex justify="center" align="center" minH="400px">
                    <Text color="neon.pink" fontSize="lg">
                        {error}
                    </Text>
                </Flex>
            ) : (
                <>
                    <DataShowcase
                        data={items}
                        columns={itemColumns}
                        total={total}
                        page={page}
                        pageSize={pageSize}
                        onPageChange={setPage}
                        onPageSizeChange={(newSize) => {
                            setPageSize(newSize);
                            setPage(1);
                        }}
                        selectable={true}
                        selectedRows={selectedRows}
                        onSelectionChange={handleSelectionChange}
                        onAdd={handleOpenCreatePage}
                        onEdit={handleOpenEditPage}
                        onDelete={handleDeleteItems}
                        getRowId={getRowId}
                        footerActions={
                            <ListingBulkCashbackButton onClick={handleOpenBulkCashbackModal} />
                        }
                    />
                </>
            )}
            <BulkCashbackModal
                isOpen={isBulkCashbackModalOpen}
                selectedItems={selectedItems}
                loading={bulkCashbackLoading}
                error={bulkCashbackError}
                onClose={handleCloseBulkCashbackModal}
                onSubmit={handleBulkCashbackSubmit}
            />
        </Flex>
    );
}
