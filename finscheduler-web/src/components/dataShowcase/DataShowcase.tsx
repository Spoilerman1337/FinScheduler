import {Card, Flex, Stack, Text} from "@chakra-ui/react";
import ListingActionButtons from "../dataListing/actionButtons/ListingActionButtons.tsx";
import ListingPaginator from "../dataListing/paginator/ListingPaginator.tsx";
import ListingSelectionCheckbox from "../dataListing/selectionCheckbox/ListingSelectionCheckbox.tsx";
import {showcasePageSizeOptions, type DataListingColumn} from "../dataListing/types.ts";
import DataShowcaseGrid from "./subcomponents/DataShowcaseGrid.tsx";

export type {DataListingColumn} from "../dataListing/types.ts";

interface DataShowcaseProps<T> {
    data: T[];
    columns: DataListingColumn<T>[];
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

export default function DataShowcase<T extends object>(props: DataShowcaseProps<T>) {
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
    const handleDeleteSelected = () => {
        if (onDelete && selectedIds.length > 0) {
            onDelete(selectedIds);
        }
    };

    const allSelected = selectable && data.length > 0 &&
        data.every((row, index) => selectedRows.has(resolveRowId(row, index)));

    return (
        <Stack width="100%" gap={4}>
            {selectable && data.length > 0 && (
                <Card.Root>
                    <Card.Body py="4">
                        <Flex align="center" gap={3}>
                            <ListingSelectionCheckbox
                                checked={allSelected}
                                onCheckedChange={handleSelectAll}
                                containerProps={{
                                    as: "div",
                                    py: 0,
                                    px: 0,
                                    flexBasis: "auto",
                                    minWidth: "auto",
                                    maxWidth: "none",
                                }}
                            />
                            <Stack gap={0}>
                                <Text color="app.accent" fontWeight="semibold">
                                    {allSelected ? "Все карточки на странице выбраны" : "Выделить все карточки на странице"}
                                </Text>
                                <Text color="fg.muted" fontSize="sm">
                                    Выбрано: {selectedIds.length} из {data.length}
                                </Text>
                            </Stack>
                        </Flex>
                    </Card.Body>
                </Card.Root>
            )}

            <DataShowcaseGrid
                data={data}
                columns={props.columns}
                selectable={selectable}
                selectedRows={selectedRows}
                getRowId={resolveRowId}
                onSelectRow={handleSelectRow}
                onEdit={onEdit}
                onRowEdit={handleEdit}
            />

            {props.total > 0 && (
                <ListingPaginator
                    total={props.total}
                    page={props.page || 1}
                    pageSize={props.pageSize || showcasePageSizeOptions[0]}
                    pageSizeOptions={showcasePageSizeOptions}
                    onPageChange={props.onPageChange || (() => {})}
                    onPageSizeChange={props.onPageSizeChange || (() => {})}
                />
            )}
            <ListingActionButtons
                selectedCount={selectedIds.length}
                onAdd={onAdd}
                onDeleteSelected={onDelete ? handleDeleteSelected : undefined}
            />
        </Stack>
    );
}
