import Layout from "../main/subcomponents/Layout.tsx";
import DataTable, {type TableColumn} from "../dataTable/DataTable.tsx";
import {Badge, Box, Button, Flex, Input, Text, Spinner} from "@chakra-ui/react";
import {SearchIcon} from "lucide-react";
import {useState, useEffect, useCallback} from "react";
import {itemsApi} from "../../../api/items.ts";
import type {ItemDto} from "../../../api/types.ts";

export default function Items() {
    const [items, setItems] = useState<ItemDto[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);

    const itemColumns: TableColumn<ItemDto>[] = [
        {
            header: 'ID',
            key: 'id',
            render: (row) => (
                <Text color="neon.blue" fontSize="sm">
                    {row.id ? row.id.substring(0, 8) + '...' : '-'}
                </Text>
            ),
            cellProps: {flexBasis: '120px', minWidth: '120px'}
        },
        {
            header: 'Название',
            key: 'name',
            render: (row) => (
                <Text fontWeight="semibold" color="neon.blue">
                    {row.name || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'}
        },
        {
            header: 'Описание',
            key: 'description',
            render: (row) => (
                <Text color="textMuted" fontSize="sm" overflow="hidden" textOverflow="ellipsis" whiteSpace="nowrap" maxW="200px">
                    {row.description || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'}
        },
        {
            header: 'Цена',
            key: 'price',
            render: (row) => (
                <Text color="neon.blue" fontWeight="medium">
                    {row.price !== undefined ? `₽${row.price.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}` : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-end'}
        },
        {
            header: 'Кэшбэк (%)',
            key: 'cashback',
            render: (row) => (
                <Text color="neon.green" fontWeight="bold">
                    {row.cashback !== undefined ? `${row.cashback}%` : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-end'}
        },
        {
            header: 'Статус',
            key: 'isActive',
            render: (row) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.isActive ? "neon.green" : "textMuted"}
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
            render: (row) => (
                <Text color="textMuted" fontSize="sm">
                    {row.createdAt ? new Date(row.createdAt).toLocaleDateString('ru-RU') : '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'}
        },
    ];

    // **Состояния для фильтров**
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<'All' | 'Active' | 'Inactive'>('All');
    const [priceFrom, setPriceFrom] = useState<string>('');
    const [priceTo, setPriceTo] = useState<string>('');

    // **Загрузка данных**
    const loadItems = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const filter: any = {
                page: page - 1, // API использует 0-based индексацию
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

    // **Функция сброса**
    const handleReset = () => {
        setSearchTerm('');
        setStatusFilter('All');
        setPriceFrom('');
        setPriceTo('');
        setPage(1);
    }

    // *** ОБЩАЯ НАСТРОЙКА ШИРИНЫ ***
    const filterWidthProps = {
        // На больших экранах: 1/5 ширины (5 элементов в строке)
        w: { base: '100%', md: 'calc(50% - var(--chakra-space-2))', xl: 'calc(20% - var(--chakra-space-3))' },
    };

    const statusOptions: Array<'All' | 'Active' | 'Inactive'> = ['All', 'Active', 'Inactive'];

    return (<Layout headerName={"Виды расходов"}>
        <Flex direction="column" width="100%">

            {/* START: Основной контейнер фильтров с управляемым переносом */}
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

                {/* 1. Фильтр поиска (Input) */}
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1" pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
                    </Box>
                    <Input
                        placeholder="Поиск по названию..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                setPage(1);
                                loadItems();
                            }
                        }}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="textPrimary"
                    />
                </Box>

                {/* 2. Группа кнопок для статуса (Заменяет Select) */}
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
                                bg={statusFilter === status ? 'neon.blue' : 'transparent'}
                                color={statusFilter === status ? 'bg.base' : 'textPrimary'}
                                _hover={{ bg: statusFilter === status ? 'neon.blue' : 'bg.layer2' }}
                                _active={{ bg: 'neon.blue' }}
                            >
                                {status === 'All' ? 'Все' : status === 'Active' ? 'Активные' : 'Неактивные'}
                            </Button>
                        ))}
                    </Flex>
                </Box>

                {/* 3. Фильтр цены от */}
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
                        color="textPrimary" 
                    />
                </Box>

                {/* 4. Фильтр цены до */}
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
                        color="textPrimary" 
                    />
                </Box>

                {/* 5. Кнопка сброса */}
                <Button
                    {...filterWidthProps}
                    onClick={handleReset}
                    bg="bg.layer3"
                    color="textMuted"
                    borderColor="glass.border"
                    border="1px solid"
                    _hover={{ bg: 'neon.pink', color: 'bg.base' }}
                >
                    Сброс
                </Button>

            </Flex>
            {/* END: Секция фильтров */}

            {loading ? (
                <Flex justify="center" align="center" minH="400px">
                    <Spinner size="xl" color="neon.blue" />
                </Flex>
            ) : error ? (
                <Flex justify="center" align="center" minH="400px">
                    <Text color="neon.pink" fontSize="lg">{error}</Text>
                </Flex>
            ) : (
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
                />
            )}
        </Flex>
    </Layout>)
}
