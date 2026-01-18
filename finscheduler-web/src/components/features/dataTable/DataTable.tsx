import {
    Flex,
    Text,
    Button,
    Checkbox,
} from "@chakra-ui/react";
import {ArrowUpDownIcon, PlusIcon} from "lucide-react";
import Paginator from "./subcomponents/Paginator.tsx";
import React from "react";
import DataTableActionButtons from "./subcomponents/DataTableActionButtons.tsx";

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
    page?: number;
    pageSize?: number;
    onPageChange?: (page: number) => void;
    onPageSizeChange?: (pageSize: number) => void;
    selectable?: boolean;
    selectedRows?: Set<string>;
    onSelectionChange?: (selectedIds: Set<string>) => void;
    onAdd?: () => void;
    onEdit?: (row: T) => void;
    onDelete?: (ids: string[]) => void;
    getRowId?: (row: T) => string;
}

export default function DataTable<T extends object>(props: DataTableProps<T>) {
    const {
        data = [],
        selectable = false,
        selectedRows = new Set(),
        onSelectionChange,
        onAdd,
        onEdit,
        onDelete,
        getRowId,
    } = props;

    const handleSelectAll = (checked: boolean) => {
        if (!onSelectionChange || !getRowId) return;
        if (checked) {
            const allIds = new Set(data.map(row => getRowId(row)));
            onSelectionChange(allIds);
        } else {
            onSelectionChange(new Set());
        }
    };

    const handleSelectRow = (rowId: string, checked: boolean) => {
        if (!onSelectionChange) return;
        const newSelection = new Set(selectedRows);
        if (checked) {
            newSelection.add(rowId);
        } else {
            newSelection.delete(rowId);
        }
        onSelectionChange(newSelection);
    };

    const handleEdit = (row: T) => {
        if (onEdit) {
            onEdit(row);
        }
    };

    const allSelected = selectable && data.length > 0 &&
        data.every(row => selectedRows.has(getRowId ? getRowId(row) : String(row['id' as keyof T] || '')));

    if (!data || data.length === 0) {
        return (
            <Flex direction="column" width="100%" height="100%" minHeight={0}>
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
                    <Text color="neon.purple" fontSize="xl">Данные не найдены.</Text>
                </Flex>
                <Flex
                    mt={2}
                    p={4}
                    bg="bg.layer2"
                    borderRadius="lg"
                    border="1px solid"
                    borderColor="glass.border"
                    justify="flex-end"
                    gap={2}
                >
                    {onAdd && (
                        <Button
                            onClick={onAdd}
                            color="neon.blue"
                            borderColor="neon.blue"
                            backdropFilter="blur(12px)"
                            bg="glass.bgHover"
                            border="1px solid"
                            transition="all 0.3s ease-in-out"
                            _hover={{
                                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                color: "neon.purple",
                                bg: "glass.bgHover",
                                backdropFilter: "blur(12px)",
                                borderColor: "neon.purple",
                            }}
                            focusRing="none"
                        >
                            <PlusIcon size={18} style={{marginRight: '8px'}}/>
                            Добавить
                        </Button>
                    )}
                </Flex>
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
                    {data.length > 0 ? (
                        <Flex
                            as="table"
                            width="100%"
                            minWidth="min-content"
                            borderCollapse="collapse"
                            flexDirection="column"
                            minHeight={0}
                            flex={1}
                        >
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
                                    {selectable && (
                                        <Flex
                                            as="th"
                                            color="textMuted"
                                            textTransform="uppercase"
                                            fontSize="xs"
                                            fontWeight="bold"
                                            py={4}
                                            px={5}
                                            flexBasis="50px"
                                            minWidth="50px"
                                            maxWidth="50px"
                                            align="center"
                                            justify="center"
                                        >
                                            <Checkbox.Root
                                                checked={allSelected}
                                                onCheckedChange={(details) => {
                                                    const checked = details.checked === true;
                                                    handleSelectAll(checked);
                                                }}
                                            >
                                                <Checkbox.HiddenInput/>
                                                <Checkbox.Control
                                                    filter={"drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"}
                                                    transition="all 0.3s ease-in-out"
                                                    color="neon.blue"
                                                    borderColor="neon.blue"
                                                    _hover={{
                                                        boxShadow: "0 0 12px rgba(0, 212, 255, 0.9)",
                                                        filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"
                                                    }}>
                                                    <Checkbox.Indicator/>
                                                </Checkbox.Control>
                                            </Checkbox.Root>
                                        </Flex>
                                    )}
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

                            <Flex
                                as="tbody"
                                flexDirection="column"
                                flex={1}
                                minHeight="0"
                                overflowY="auto"
                                overflowX="hidden"
                                className="custom-scrollbar"
                            >
                                {data.map((row, index) => {
                                    const rowId = getRowId ? getRowId(row) : (row['id' as keyof T] ? String(row['id' as keyof T]) : `row-${index}`);
                                    const isSelected = selectable && rowId && selectedRows.has(rowId);

                                    return (
                                        <Flex
                                            key={rowId}
                                            as="tr"
                                            bg={isSelected ? "bg.accent" : "bg.layer1"}
                                            borderBottom="1px solid"
                                            borderColor="glass.border"
                                            transition="all 0.3s ease-in-out"
                                            cursor={onEdit ? "pointer" : "default"}
                                            width="100%"
                                            boxSizing="border-box"
                                            display="flex"
                                            maxWidth="100%"
                                            onClick={() => onEdit && handleEdit(row)}
                                            _hover={{
                                                bg: "bg.accent",
                                                boxShadow: "0 0 12px 0 neon.blue4D",
                                            }}
                                        >
                                            {selectable && (
                                                <Flex
                                                    as="td"
                                                    py={3}
                                                    px={5}
                                                    flexBasis="50px"
                                                    minWidth="50px"
                                                    maxWidth="50px"
                                                    align="center"
                                                    justify="center"
                                                    onClick={(e) => e.stopPropagation()}
                                                >
                                                    <Checkbox.Root
                                                        checked={isSelected || false}
                                                        onCheckedChange={(details) => {
                                                            const checked = details.checked === true;
                                                            handleSelectRow(rowId, checked);
                                                        }}
                                                    >
                                                        <Checkbox.HiddenInput/>
                                                        <Checkbox.Control
                                                            filter={"drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"}
                                                            transition="all 0.3s ease-in-out"
                                                            color="neon.blue"
                                                            borderColor="neon.blue"
                                                            _hover={{
                                                                boxShadow: "0 0 12px rgba(0, 212, 255, 0.9)",
                                                                filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"
                                                            }}>
                                                            <Checkbox.Indicator/>
                                                        </Checkbox.Control>
                                                    </Checkbox.Root>
                                                </Flex>
                                            )}
                                            {props.columns.map((col) => {
                                                const typedRow = row as T;
                                                return (
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
                                                        {col.render(typedRow)}
                                                    </Flex>
                                                );
                                            })}
                                        </Flex>
                                    );
                                })}
                            </Flex>
                        </Flex>
                    ) : <Text color="neon.purple" fontSize="xl">Данные не найдены.</Text>}
                </Flex>

                <Paginator
                    total={props.total}
                    page={props.page || 1}
                    pageSize={props.pageSize || 10}
                    onPageChange={props.onPageChange || (() => {
                    })}
                    onPageSizeChange={props.onPageSizeChange || (() => {
                    })}
                />
                <DataTableActionButtons
                    data={data}
                    selectedRows={selectedRows}
                    onAdd={onAdd}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    getRowId={getRowId}
                    handleEdit={handleEdit}
                />
            </Flex>
        </>
    );
}