import {
    Flex,
    Text,
} from "@chakra-ui/react";
import {ArrowUpDownIcon} from "lucide-react";
import Paginator from "./subcomponents/Paginator.tsx";

export interface TableColumn<TData> {
    header: string | React.ReactNode;
    key: string;
    render: (row: TData) => React.ReactNode;
    headerProps?: React.ComponentProps<typeof Flex>;
    cellProps?: React.ComponentProps<typeof Flex>;
}

interface DataTableProps<T> {
    data: T[];
    columns: TableColumn<T>[];
    rowKey?: keyof T;
    total: number;
}

export default function DataTable<T extends object>(props: DataTableProps<T>) {

    // Гарантируем, что компонент не сломается, если данных нет
    if (!props.data || props.data.length === 0) {
        return (
            <Flex
                flex={1}
                height="100%"
                width="100%"
                bg="bg.layer1"
                borderRadius="xl"
                border="1px solid"
                borderColor="borderStrong"
                boxShadow="card"
                align="center"
                justify="center"
            >
                <Text color="textMuted" fontSize="xl">Данные не найдены.</Text>
            </Flex>
        );
    }

    return (<>
            <Flex direction="column" width="100%" height="100%" minHeight={0}>
                <Flex
                    flex={1}
                    height="100%"
                    width="100%"
                    bg="bg.layer1"
                    borderColor="borderStrong"
                    boxShadow="card"
                    p={0}
                    background="gradients.cosmic"
                    flexDirection="column"
                    borderRadius="xl"
                    border="1px solid rgba(255,255,255,0.2)"
                    minHeight={0}
                    overflow="hidden"
                >
                    {/* 2. Таблица (Flex as="table") */}
                    <Flex
                        as="table"
                        width="100%"
                        minWidth="min-content" // Для горизонтального скролла
                        borderCollapse="collapse"
                        flexDirection="column"
                        minHeight={0}
                        flex={1} // Растягиваем таблицу на 100% высоты контейнера
                    >

                        {/* 3. Заголовки таблицы (Thead) - С фиксированной позицией */}
                        <Flex
                            as="thead"
                            position="sticky"
                            top="0"
                            bg="bg.layer2"
                            zIndex={10}
                            boxShadow="sm"
                        >
                            <Flex
                                as="tr"
                                borderBottom="1px solid"
                                borderColor="glass.borderStrong"
                                width="100%"
                                display="flex"
                            >
                                {props.columns.map((col) => (
                                    <Flex
                                        key={col.key}
                                        as="th"
                                        color="textMuted"
                                        textTransform="uppercase"
                                        fontSize="xs"
                                        fontWeight="bold"
                                        py={4}
                                        px={5}
                                        flex={1}
                                        minWidth="100px"
                                        {...col.headerProps}
                                    >
                                        {col.key === 'asset' ? (
                                            <Flex align="center" cursor="pointer" _hover={{color: "neon.blue"}}>
                                                <Text color="neon.blue">{col.header}</Text>
                                                <ArrowUpDownIcon size={14} style={{marginLeft: '8px'}}/>
                                            </Flex>
                                        ) : (
                                            <Text color="neon.blue">{col.header}</Text>
                                        )}
                                    </Flex>
                                ))}
                            </Flex>
                        </Flex>

                        {/* 4. Тело таблицы (Tbody) - Основная область скролла */}
                        <Flex
                            as="tbody"
                            flexDirection="column"
                            flex={1}
                            minHeight="0"
                            overflowY="auto"
                            overflowX="hidden"
                            className="custom-scrollbar"
                        >
                            {props.data.map((row, index) => {
                                const key = row['id' as keyof T] ? String(row['id' as keyof T]) : index;

                                return (
                                    <Flex
                                        key={key}
                                        as="tr"
                                        bg="bg.layer1"
                                        borderBottom="1px solid"
                                        borderColor="glass.border"
                                        transition="all 0.3s ease-in-out"
                                        cursor="pointer"
                                        width="100%"
                                        boxSizing="border-box"
                                        display="flex"
                                        maxWidth="100%"
                                        _hover={{
                                            bg: "bg.accent",
                                            boxShadow: `0 0 12px 0 ${'neon.blue'}4D`,
                                        }}
                                    >
                                        {props.columns.map((col) => (
                                            <Flex
                                                key={col.key}
                                                as="td"
                                                py={3}
                                                px={5}
                                                flex={1}
                                                minWidth="0"
                                                align="center"
                                                overflow="hidden"
                                                textOverflow="ellipsis"
                                                whiteSpace="nowrap"
                                                flexShrink={1}
                                                {...col.cellProps}
                                            >
                                                {col.render(row)}
                                            </Flex>
                                        ))}
                                    </Flex>
                                );
                            })}
                        </Flex>
                    </Flex>
                </Flex>
                <Paginator total={props.total} />
            </Flex>
        </>
    );
}