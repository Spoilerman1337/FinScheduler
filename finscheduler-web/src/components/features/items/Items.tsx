import Layout from "../main/subcomponents/Layout.tsx";
import DataTable, {type TableColumn} from "../dataTable/DataTable.tsx";
import {Badge, Box, Button, Flex, Input, Text, Spinner} from "@chakra-ui/react";
import {SearchIcon} from "lucide-react";
import {useState, useEffect, useCallback} from "react";
import {itemsApi} from "../../../api/items.ts";
import type {ItemDto} from "../../../api/types.ts";
import ItemModal from "./subcomponents/ItemModal.tsx";
import {toaster} from "../../ui/toaster.tsx";

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

    const itemColumns: TableColumn<ItemDto>[] = [
        {
            header: 'ID',
            key: 'Id',
            render: (row: ItemDto) => (
                <Text color="neon.blue" fontSize="sm">
                    {row.id ? row.id.substring(0, 4) + '...' : '-'}
                </Text>
            ),
            cellProps: {flexBasis: '80px', minWidth: '80px', maxWidth: '80px', flexShrink: 0},
            headerProps: {flexBasis: '80px', minWidth: '80px', maxWidth: '80px', flexShrink: 0}
        },
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
                        maximumFractionDigits: 2
                    })}` : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-start'}
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
            cellProps: {justifyContent: 'flex-start'}
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
            headerProps: {textAlign: 'left'}
        },
        {
            header: 'Создан',
            key: 'createdAt',
            render: (row: ItemDto) => (
                <Text color="neon.blue" fontSize="sm">
                    {row.createdAt ? new Date(row.createdAt).toLocaleDateString('ru-RU') : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'}
        },
    ];

    const getCashbackColor = (cashback: number | undefined) => {
        let color = "neon.pink";

        if (cashback) {
            color = cashback > 1 ? "neon.green" : "neon.yellow";
        }

        return color
    }

    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<'All' | 'Active' | 'Inactive'>('Active');
    const [priceFrom, setPriceFrom] = useState<string>('');
    const [priceTo, setPriceTo] = useState<string>('');

    const loadItems = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const filter: any = {
                page: page - 1,
                pageSize: pageSize,
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

            const result = await itemsApi.getItems(filter);
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
    }

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

    const handleSaveItem = async (itemData: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>) => {
        try {
            if (modalMode === 'create') {
                await itemsApi.createItem(itemData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Элемент успешно добавлен',
                    type: 'success',
                });
            } else if (modalMode === 'edit' && editingItem?.id) {
                await itemsApi.updateItem(editingItem.id, itemData);
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
            await Promise.all(ids.map(id => itemsApi.deleteItem(id)));
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
        return row.id || `item-${row.name || 'unknown'}-${row.price || 0}`;
    };

    const filterWidthProps = {
        w: {base: '100%', md: 'calc(50% - var(--chakra-space-2))', xl: 'calc(20% - var(--chakra-space-3))'},
    };

    const statusOptions: Array<'All' | 'Active' | 'Inactive'> = ['All', 'Active', 'Inactive'];

    return (<Layout headerName={"Виды расходов"}>
        <Flex direction="column" width="100%">

            <Flex
                mb={4}
                p={4}
                bg="bg.layer2"
                borderRadius="lg"
                border="1px solid"
                borderColor="glass.border"
                width="100%"
                align="center"
                gap={4}
                flexWrap="wrap"
                justifyContent="flex-start"
            >

                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1"
                         pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)"/>
                    </Box>
                    <Input
                        placeholder="Поиск по названию..."
                        value={searchTerm}
                        onChange={(e) => {
                            setPage(1)
                            setSearchTerm(e.target.value)
                        }}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                setPage(1);
                                loadItems();
                            }
                        }}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="neon.blue"
                        _placeholder={{color: 'textMuted'}}
                    />
                </Box>

                <Box {...filterWidthProps}>
                    <Flex gap={1} borderRadius="md" p={1} bg="bg.layer1" border="1px solid" borderColor="glass.border">
                        {statusOptions.map(status => (
                            <Button
                                key={status}
                                size="sm"
                                flex={1}
                                onClick={() => {
                                    setStatusFilter(status);
                                    setPage(1);
                                }}
                                color={statusFilter === status ? "neon.blue" : "neon.blue"}
                                borderColor={statusFilter === status ? "neon.blue" : "glass.border"}
                                backdropFilter="blur(12px)"
                                bg={statusFilter === status ? "glass.bgHover" : "transparent"}
                                border="1px solid"
                                transition="all 0.3s ease-in-out"
                                filter={statusFilter === status ? "drop-shadow(0 0 16px rgba(0,212,255,0.9))" : "none"}
                                boxShadow={statusFilter === status ? "0 0 20px rgba(0,212,255,1)" : "none"}
                                _hover={{
                                    filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                    boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                    backdropFilter: "blur(12px)",
                                    borderColor: "neon.purple",
                                }}
                                _active={{
                                    filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                                    boxShadow: "0 0 20px rgba(0,212,255,1)",
                                    color: "neon.blue",
                                    bg: "glass.bgHover",
                                    backdropFilter: "blur(12px)",
                                    borderColor: "neon.blue",
                                }}
                                focusRing="none"
                            >
                                {status === 'All' ? 'Все' : status === 'Active' ? 'Активные' : 'Неактивные'}
                            </Button>
                        ))}
                    </Flex>
                </Box>

                <Box {...filterWidthProps}>
                    <Input
                        placeholder="Цена от"
                        type="number"
                        value={priceFrom}
                        onChange={(e) => setPriceFrom(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                setPage(1);
                                loadItems();
                            }
                        }}
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="neon.blue"
                        _placeholder={{color: 'textMuted'}}
                    />
                </Box>

                <Box {...filterWidthProps}>
                    <Input
                        placeholder="Цена до"
                        type="number"
                        value={priceTo}
                        onChange={(e) => setPriceTo(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                setPage(1);
                                loadItems();
                            }
                        }}
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="neon.blue"
                        _placeholder={{color: 'textMuted'}}
                    />
                </Box>

                <Button
                    {...filterWidthProps}
                    onClick={handleReset}
                    bg="bg.layer3"
                    color="textMuted"
                    borderColor="glass.border"
                    border="1px solid"
                    _hover={{bg: 'neon.pink', color: 'bg.base'}}
                >
                    Сброс
                </Button>

            </Flex>

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
    </Layout>)
}
