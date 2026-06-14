import type {ReactNode} from 'react';
import ListingActionButtons from './actionButtons/ListingActionButtons.tsx';
import ListingPaginator from './paginator/ListingPaginator.tsx';
import {defaultPageSizeOptions} from './types.ts';

export interface DataListingCommonProps<T> {
    data: T[];
    total: number;
    page?: number;
    pageSize?: number;
    pageSizeOptions?: readonly number[];
    onPageChange?: (page: number) => void;
    onPageSizeChange?: (pageSize: number) => void;
    selectable?: boolean;
    selectedRows?: Set<string>;
    onSelectionChange?: (selectedIds: Set<string>) => void;
    onAdd?: () => void;
    onEdit?: (row: T) => void;
    onDelete?: (ids: string[]) => void;
    getRowId?: (row: T) => string;
    footerActions?: ReactNode;
    showEditSelectedAction?: boolean;
}

export interface DataListingContainerRenderProps<T> {
    data: T[];
    selectable: boolean;
    selectedRows: Set<string>;
    selectedIds: string[];
    selectedRow?: T;
    allSelected: boolean;
    resolveRowId: (row: T, index: number) => string;
    handleSelectAll: (checked: boolean) => void;
    handleSelectRow: (rowId: string, checked: boolean) => void;
    handleEdit: (row: T) => void;
    handleEditSelected: () => void;
    handleDeleteSelected: () => void;
}

interface DataListingContainerProps<T> extends DataListingCommonProps<T> {
    children: (props: DataListingContainerRenderProps<T>) => ReactNode;
}

export default function DataListingContainer<T extends object>(
    props: DataListingContainerProps<T>,
) {
    const {
        data = [],
        selectable = false,
        selectedRows = new Set(),
        onSelectionChange,
        onAdd,
        onEdit,
        onDelete,
        getRowId,
        footerActions,
        showEditSelectedAction = true,
    } = props;

    const resolveRowId = (row: T, index: number) => {
        if (getRowId) {
            return getRowId(row);
        }

        const fallbackId = row['id' as keyof T];
        return fallbackId ? String(fallbackId) : `row-${index}`;
    };

    const handleSelectAll = (checked: boolean) => {
        if (!onSelectionChange) {
            return;
        }

        if (checked) {
            const allIds = new Set(data.map((row, index) => resolveRowId(row, index)));
            onSelectionChange(allIds);
            return;
        }

        onSelectionChange(new Set());
    };

    const handleSelectRow = (rowId: string, checked: boolean) => {
        if (!onSelectionChange) {
            return;
        }

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
    const selectedRow =
        selectedIds.length === 1
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

    const allSelected =
        selectable &&
        data.length > 0 &&
        data.every((row, index) => selectedRows.has(resolveRowId(row, index)));

    const page = props.page ?? 1;
    const pageSize = props.pageSize ?? props.pageSizeOptions?.[0] ?? defaultPageSizeOptions[0];

    return (
        <>
            {props.children({
                data,
                selectable,
                selectedRows,
                selectedIds,
                selectedRow,
                allSelected,
                resolveRowId,
                handleSelectAll,
                handleSelectRow,
                handleEdit,
                handleEditSelected,
                handleDeleteSelected,
            })}

            {props.total > 0 && (
                <ListingPaginator
                    total={props.total}
                    page={page}
                    pageSize={pageSize}
                    pageSizeOptions={props.pageSizeOptions}
                    onPageChange={props.onPageChange ?? (() => {})}
                    onPageSizeChange={props.onPageSizeChange ?? (() => {})}
                />
            )}
            <ListingActionButtons
                selectedCount={selectedIds.length}
                onAdd={onAdd}
                onEditSelected={
                    showEditSelectedAction && onEdit && selectedRow ? handleEditSelected : undefined
                }
                onDeleteSelected={onDelete ? handleDeleteSelected : undefined}
                additionalActions={footerActions}
            />
        </>
    );
}
