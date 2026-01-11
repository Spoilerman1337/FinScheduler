import Layout from "../main/subcomponents/Layout.tsx";
import DataTable, {type TableColumn} from "../dataTable/DataTable.tsx";
import {Badge, Box, Button, Flex, Input, Text} from "@chakra-ui/react";
import {SearchIcon} from "lucide-react";
import {useState} from "react";

export default function Items() {
    interface CryptoData {
        id: number;
        asset: string;
        symbol: string;
        price: number;
        change24h: number;
        status: "Active" | "Inactive" | "Trending";
        colorToken: 'btc' | 'eth' | 'ltc' | 'neon.purple' | 'neon.pink';
    }

    const mockData: CryptoData[] = [
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
        {
            id: 1,
            asset: "Bitcoin",
            symbol: "BTC",
            price: 65123.45,
            change24h: 3.45,
            status: "Trending",
            colorToken: 'btc'
        },
        {
            id: 2,
            asset: "Ethereum",
            symbol: "ETH",
            price: 3456.78,
            change24h: -1.89,
            status: "Active",
            colorToken: 'eth'
        },
        {id: 3, asset: "Litecoin", symbol: "LTC", price: 120.50, change24h: 5.01, status: "Active", colorToken: 'ltc'},
    ];

    const cryptoColumns: TableColumn<CryptoData>[] = [
        {
            header: '#',
            key: 'id',
            render: (row) => (<Text color="neon.blue">{row.id}</Text>),
            cellProps: {flexBasis: '50px', minWidth: '50px'} // Фиксированная ширина для ID
        },
        {
            header: 'Актив',
            key: 'asset',
            render: (row) => (
                <Flex flexDirection="column">
                    <Text fontWeight="semibold" color="neon.blue">{row.asset}</Text>
                    <Text fontSize="sm" color="textMuted">{row.symbol}</Text>
                </Flex>
            ),
            headerProps: {textAlign: 'left'}
        },
        {
            header: 'Цена (USD)',
            key: 'price',
            render: (row) => (
                <Text color="neon.blue" fontWeight="medium">
                    ${row.price.toLocaleString(undefined, {minimumFractionDigits: 2})}
                </Text>
            ),
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-end'}
        },
        {
            header: 'Изменение 24ч (%)',
            key: 'change24h',
            render: (row) => {
                const color = row.change24h > 0 ? "neon.green" : row.change24h < 0 ? "neon.pink" : "textMuted";
                const sign = row.change24h > 0 ? '+' : '';
                return <Text color={color} fontWeight="bold">{sign}{row.change24h.toFixed(2)}%</Text>;
            },
            headerProps: {textAlign: 'right'},
            cellProps: {justifyContent: 'flex-end'}
        },
        {
            header: 'Статус',
            key: 'status',
            render: (row) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.status === "Trending" ? "neon.purple" : row.status === "Active" ? "neon.green" : "textMuted"}
                    color="bg.base"
                >
                    {row.status}
                </Badge>
            ),
            headerProps: {textAlign: 'left'}
        },
    ];

    // **Состояния для фильтров**
    const [searchTerm, setSearchTerm] = useState('');
    // Используем строку для статуса
    const [statusFilter, setStatusFilter] = useState('All');

    // // **Логика фильтрации**
    // const filteredData = useMemo(() => {
    //     let currentData = mockData;
    //
    //     // Фильтрация по поисковому запросу
    //     if (searchTerm) {
    //         const lowerCaseSearch = searchTerm.toLowerCase();
    //         currentData = currentData.filter(item =>
    //             item.asset.toLowerCase().includes(lowerCaseSearch) ||
    //             item.symbol.toLowerCase().includes(lowerCaseSearch)
    //         );
    //     }
    //
    //     // Фильтрация по статусу
    //     if (statusFilter !== 'All') {
    //         // TypeScript будет жаловаться, если статус не соответствует CryptoData["status"],
    //         // но для простоты примера мы используем приведение типов.
    //         currentData = currentData.filter(item => item.status === statusFilter as CryptoData["status"]);
    //     }
    //
    //     return currentData;
    // }, [searchTerm, statusFilter, mockData]);

    // **Функция сброса**
    const handleReset = () => {
        setSearchTerm('');
        setStatusFilter('All');
    }

    // *** ОБЩАЯ НАСТРОЙКА ШИРИНЫ ***
    const filterWidthProps = {
        // На больших экранах: 1/5 ширины (5 элементов в строке)
        w: { base: '100%', md: 'calc(50% - var(--chakra-space-2))', xl: 'calc(20% - var(--chakra-space-3))' },
    };

    const statusOptions = ['All', 'Active', 'Trending', 'Inactive'];

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
                        placeholder="Поиск по активу или символу..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="textPrimary"
                    />
                </Box>
                {/* 1. Фильтр поиска (Input) */}
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1" pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
                    </Box>
                    <Input
                        placeholder="Поиск по активу или символу..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="textPrimary"
                    />
                </Box>
                {/* 1. Фильтр поиска (Input) */}
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1" pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
                    </Box>
                    <Input
                        placeholder="Поиск по активу или символу..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="textPrimary"
                    />
                </Box>
                {/* 1. Фильтр поиска (Input) */}
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1" pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
                    </Box>
                    <Input
                        placeholder="Поиск по активу или символу..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="textPrimary"
                    />
                </Box>
                {/* 1. Фильтр поиска (Input) */}
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1" pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
                    </Box>
                    <Input
                        placeholder="Поиск по активу или символу..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
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
                                onClick={() => setStatusFilter(status)}
                                // Подсветка активной кнопки
                                bg={statusFilter === status ? 'neon.blue' : 'transparent'}
                                color={statusFilter === status ? 'bg.base' : 'textPrimary'}
                                _hover={{ bg: statusFilter === status ? 'neon.blue' : 'bg.layer2' }}
                                _active={{ bg: 'neon.blue' }}
                            >
                                {status === 'All' ? 'Все' : status}
                            </Button>
                        ))}
                    </Flex>
                </Box>

                {/* 3. Произвольный фильтр 1 (Категория - Input) */}
                <Box {...filterWidthProps}>
                    <Input placeholder="Категория" bg="bg.layer1" borderColor="glass.border" color="textPrimary" />
                </Box>

                {/* 4. Произвольный фильтр 2 (Диапазон цен - Input) */}
                <Box {...filterWidthProps}>
                    <Input placeholder="Диапазон цен" bg="bg.layer1" borderColor="glass.border" color="textPrimary" />
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

            <DataTable
                data={mockData}
                columns={cryptoColumns}
                total={0}
            />
        </Flex>
    </Layout>)
}