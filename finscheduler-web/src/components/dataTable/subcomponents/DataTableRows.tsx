import {Flex, Text} from "@chakra-ui/react";
import DataTableCell from "./DataTableCell.tsx";
import DataTableRow from "./DataTableRow.tsx";
import type {TableColumn} from "../models.ts";

interface DataTableRowsProps<T> {
    data: T[];
    columns: TableColumn<T>[];
    selectable: boolean;
    selectedRows: Set<string>;
    getRowId: (row: T, index: number) => string;
    onSelectRow: (rowId: string, checked: boolean) => void;
    onEdit?: (row: T) => void;
    onRowEdit: (row: T) => void;
}

export default function DataTableRows<T extends object>(props: DataTableRowsProps<T>) {
    const {
        data,
        columns,
        selectable,
        selectedRows,
        getRowId,
        onSelectRow,
        onEdit,
        onRowEdit,
    } = props;

    return (
        <Flex
            as="tbody"
            flexDirection="column"
            overflowY="visible"
            overflowX="hidden"
        >
            {data.length === 0 && (
                <Flex
                    as="tr"
                    width="100%"
                    bg="bg.layer1"
                    borderBottom="1px solid"
                    borderColor="glass.border"
                >
                    <DataTableCell
                        as="td"
                        flex="1 0 100%"
                        minWidth="100%"
                        justify="center"
                        py={10}
                    >
                        <Text color="neon.purple" fontSize="xl">Данные не найдены.</Text>
                    </DataTableCell>
                </Flex>
            )}
            {data.map((row, index) => {
                const rowId = getRowId(row, index);
                const isSelected = selectable && selectedRows.has(rowId);

                return (
                    <DataTableRow
                        key={rowId}
                        row={row}
                        rowId={rowId}
                        columns={columns}
                        selectable={selectable}
                        isSelected={isSelected}
                        onSelectRow={onSelectRow}
                        onEdit={onEdit}
                        onRowEdit={onRowEdit}
                    />
                );
            })}
        </Flex>
    );
}
