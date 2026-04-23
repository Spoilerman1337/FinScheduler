import {Flex} from "@chakra-ui/react";
import Paginator from "./subcomponents/Paginator.tsx";
import DataTableActionButtons from "./subcomponents/DataTableActionButtons.tsx";
import DataTableHeader from "./subcomponents/DataTableHeader.tsx";
import DataTableRows from "./subcomponents/DataTableRows.tsx";
import type {TableColumn} from "./models.ts";

export type {TableColumn} from "./models.ts";

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

    const resolveRowId = (row: T, index: number) => {
        if (getRowId) {
            return getRowId(row);
        }

        const fallbackId = row["id" as keyof T];
        return fallbackId ? String(fallbackId) : `row-${index}`;
    };

    const handleSelectAll = (checked: boolean) => {
        if (!onSelectionChange) return;

        if (checked) {
            const allIds = new Set(data.map((row, index) => resolveRowId(row, index)));
            onSelectionChange(allIds);
            return;
        }

        onSelectionChange(new Set());
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

    const selectedIds = Array.from(selectedRows);
    const selectedRow = selectedIds.length === 1
        ? data.find((row, index) => resolveRowId(row, index) === selectedIds[0])
        : undefined;

    const handleEditSelected = () => {
        if (selectedRow) {
            handleEdit(selectedRow);
        }
    };

    const handleDeleteSelected = () => {
        if (onDelete && selectedIds.length > 0) {
            onDelete(selectedIds);
        }
    };

    const allSelected = selectable && data.length > 0 &&
        data.every((row, index) => selectedRows.has(resolveRowId(row, index)));

    return (
        <Flex direction="column" width="100%">
            <Flex
                width="100%"
                bg="bg.layer1"
                borderColor="borderStrong"
                boxShadow="card"
                p={0}
                background="gradients.cosmic"
                flexDirection="column"
                borderRadius="xl"
                border="1px solid rgba(255,255,255,0.2)"
                overflow="hidden"
            >
                <Flex
                    as="table"
                    width="100%"
                    minWidth="min-content"
                    borderCollapse="collapse"
                    flexDirection="column"
                >
                    <DataTableHeader
                        columns={props.columns}
                        selectable={selectable}
                        allSelected={allSelected}
                        onSelectAll={handleSelectAll}
                    />
                    <DataTableRows
                        data={data}
                        columns={props.columns}
                        selectable={selectable}
                        selectedRows={selectedRows}
                        getRowId={resolveRowId}
                        onSelectRow={handleSelectRow}
                        onEdit={onEdit}
                        onRowEdit={handleEdit}
                    />
                </Flex>
            </Flex>

            {props.total > 0 && (
                <Paginator
                    total={props.total}
                    page={props.page || 1}
                    pageSize={props.pageSize || 10}
                    onPageChange={props.onPageChange || (() => {
                    })}
                    onPageSizeChange={props.onPageSizeChange || (() => {
                    })}
                />
            )}
            <DataTableActionButtons
                selectedCount={selectedIds.length}
                onAdd={onAdd}
                onEditSelected={onEdit && selectedRow ? handleEditSelected : undefined}
                onDeleteSelected={onDelete ? handleDeleteSelected : undefined}
            />
        </Flex>
    );
}
