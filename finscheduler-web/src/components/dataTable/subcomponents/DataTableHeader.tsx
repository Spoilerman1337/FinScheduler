import {Flex, Text} from "@chakra-ui/react";
import {ArrowUpDownIcon} from "lucide-react";
import DataTableCell from "./DataTableCell.tsx";
import DataTableSelectionCheckbox from "./DataTableSelectionCheckbox.tsx";
import type {TableColumn} from "../models.ts";

interface DataTableHeaderProps<T> {
    columns: TableColumn<T>[];
    selectable: boolean;
    allSelected: boolean;
    onSelectAll: (checked: boolean) => void;
}

export default function DataTableHeader<T>(props: DataTableHeaderProps<T>) {
    const {columns, selectable, allSelected, onSelectAll} = props;

    return (
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
                    <DataTableSelectionCheckbox
                        checked={allSelected}
                        onCheckedChange={onSelectAll}
                        isHeader
                    />
                )}
                {columns.map((col) => (
                    <DataTableCell
                        key={col.key}
                        as="th"
                        isHeader
                        {...col.headerProps}
                    >
                        {col.key === "asset" ? (
                            <Flex align="center" cursor="pointer" _hover={{color: "neon.blue"}}>
                                <Text color="neon.blue">{col.header}</Text>
                                <ArrowUpDownIcon size={14} style={{marginLeft: "8px"}}/>
                            </Flex>
                        ) : (
                            <Text color="neon.blue">{col.header}</Text>
                        )}
                    </DataTableCell>
                ))}
            </Flex>
        </Flex>
    );
}
