import {Flex} from "@chakra-ui/react";
import type {DataListingColumn} from "../../dataListing/types.ts";
import ListingSelectionCheckbox from "../../dataListing/selectionCheckbox/ListingSelectionCheckbox.tsx";
import DataTableCell from "./DataTableCell.tsx";

interface DataTableRowProps<T> {
    row: T;
    rowId: string;
    columns: DataListingColumn<T>[];
    selectable: boolean;
    isSelected: boolean;
    onSelectRow: (rowId: string, checked: boolean) => void;
    onEdit?: (row: T) => void;
    onRowEdit: (row: T) => void;
}

export default function DataTableRow<T extends object>(props: DataTableRowProps<T>) {
    const {row, rowId, columns, selectable, isSelected, onSelectRow, onEdit, onRowEdit} = props;

    return (
        <Flex
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
            onClick={() => onEdit && onRowEdit(row)}
            _hover={{
                bg: "bg.accent",
                boxShadow: "0 0 12px 0 neon.blue4D",
            }}
        >
            {selectable && (
                <ListingSelectionCheckbox
                    checked={isSelected}
                    onCheckedChange={(checked) => onSelectRow(rowId, checked)}
                    onClick={(event) => event.stopPropagation()}
                />
            )}
            {columns.map((col) => (
                <DataTableCell
                    key={col.key}
                    as="td"
                    {...col.cellProps}
                >
                    {col.render(row)}
                </DataTableCell>
            ))}
        </Flex>
    );
}
