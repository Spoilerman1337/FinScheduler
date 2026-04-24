import DataTable, {type TableColumn} from "../../components/dataTable/DataTable.tsx";
import {Badge, Flex, Spinner, Text} from "@chakra-ui/react";
import {useState, useEffect, useCallback} from "react";
import type {ItemDto, ItemFilter, ItemModification} from "../../api/types.ts";
import ItemModal from "./subcomponents/ItemModal.tsx";
import {toaster} from "../../components/ui/toaster.tsx";
import ItemsService from "../../api/items.ts";
import ItemsFilters from "./subcomponents/ItemsFilters.tsx";

const itemsService = new ItemsService();

export default function Items() {
    const [items, setItems] = useState<ItemDto[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingItem, setEditingItem] = useState<ItemDto | null>(null);
    const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<'All' | 'Active' | 'Inactive'>('Active');
    const [priceFrom, setPriceFrom] = useState<string>('');
    const [priceTo, setPriceTo] = useState<string>('');

    const itemColumns: TableColumn<ItemDto>[] = [
        {
            header: 'Название',
            key: 'name',
            render: (row: ItemDto) => (
                <Text fontWeight="semibold" color="neon.blue">
                    {row.name || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Цена',
            key: 'price',
            render: (row: ItemDto) => (
                <Text color="neon.blue" fontWeight="medium">
                    {row.price !== undefined ? `₽${row.price.toLocaleString(undefined, {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                    })}` : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-start'},
        },
        {
            header: 'Кэшбэк (%)',
            key: 'cashback',
            render: (row: ItemDto) => (
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
            render: (row: ItemDto) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.isActive ? "neon.green" : "neon.pink"}
                    color="bg.base"
                >
                    {row.isActive ? "Активен" : "Неактивен"}
                </Badge>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Создан',
            key: 'createdAt',
            render: (row: ItemDto) => (
                <Text color="neon.blue" fontSize="sm">
                    {row.createdAt ? new Date(row.createdAt).toLocaleDateString('ru-RU') : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
    ];

    const getCashbackColor = (cashback: number | undefined) => {
        let color = "neon.pink";

        if (cashback) {
            color = cashback > 1 ? "neon.green" : "neon.yellow";
        }

        return color;
    };

    const loadItems = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const filter: ItemFilter = {
                page: page - 1,
                pageSize,
            };

            if (searchTerm) {
                filter.name = searchTerm;
            }

            if (statusFilter !== 'All') {
                filter.isActive = statusFilter === 'Active';
            }

            if (priceFrom) {
                const price = parseFloat(priceFrom);
                if (!isNaN(price)) {
                    filter.priceFrom = price;
                }
            }

            if (priceTo) {
                const price = parseFloat(priceTo);
                if (!isNaN(price)) {
                    filter.priceTo = price;
                }
            }

            const result = await itemsService.getItems(filter);
            setItems(result.data);
            setTotal(result.count);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка загрузки данных');
            console.error('Failed to load items:', err);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, searchTerm, statusFilter, priceFrom, priceTo]);

    useEffect(() => {
        loadItems();
    }, [loadItems]);

    const handleReset = () => {
        setSearchTerm('');
        setStatusFilter('All');
        setPriceFrom('');
        setPriceTo('');
        setPage(1);
    };

    const handleOpenAddModal = () => {
        setEditingItem(null);
        setModalMode('create');
        setIsModalOpen(true);
    };

    const handleOpenEditModal = (item: ItemDto) => {
        setEditingItem(item);
        setModalMode('edit');
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEditingItem(null);
        setSelectedRows(new Set());
    };

    const handleSaveItem = async (itemData: ItemModification) => {
        try {
            if (modalMode === 'create') {
                await itemsService.createItem(itemData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Элемент успешно добавлен',
                    type: 'success',
                });
            } else if (modalMode === 'edit' && editingItem?.id) {
                await itemsService.updateItem(editingItem.id, itemData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Элемент успешно обновлен',
                    type: 'success',
                });
            } else {
                console.error('Cannot save: missing id for edit mode', {modalMode, editingItem});
                throw new Error('Не удалось определить режим сохранения');
            }

            await loadItems();
        } catch (err) {
            console.error('Error saving item:', err);
            throw err;
        }
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
            await loadItems();
        } catch (err) {
            toaster.create({
                title: 'Ошибка',
                description: err instanceof Error ? err.message : 'Не удалось удалить элементы',
                type: 'error',
            });
        }
    };

    const getRowId = (row: ItemDto): string => {
        return row.id ?? '';
    };

    return (
        <Flex direction="column" width="100%">
            <ItemsFilters
                searchTerm={searchTerm}
                statusFilter={statusFilter}
                priceFrom={priceFrom}
                priceTo={priceTo}
                onSearchTermChange={(value) => {
                    setPage(1);
                    setSearchTerm(value);
                }}
                onStatusFilterChange={(value) => {
                    setStatusFilter(value);
                    setPage(1);
                }}
                onPriceFromChange={setPriceFrom}
                onPriceToChange={setPriceTo}
                onApply={() => {
                    setPage(1);
                    loadItems();
                }}
                onReset={handleReset}
            />

            {loading ? (
                <Flex justify="center" align="center" minH="400px">
                    <Spinner size="xl" color="neon.blue"/>
                </Flex>
            ) : error ? (
                <Flex justify="center" align="center" minH="400px">
                    <Text color="neon.pink" fontSize="lg">{error}</Text>
                </Flex>
            ) : (
                <>
                    <DataTable
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
                        onSelectionChange={setSelectedRows}
                        onAdd={handleOpenAddModal}
                        onEdit={handleOpenEditModal}
                        onDelete={handleDeleteItems}
                        getRowId={getRowId}
                    />
                    <ItemModal
                        key={`${modalMode}-${editingItem?.id || 'new'}`}
                        isOpen={isModalOpen}
                        onClose={handleCloseModal}
                        onSave={handleSaveItem}
                        item={editingItem}
                        mode={modalMode}
                    />
                </>
            )}
        </Flex>
    );
}
