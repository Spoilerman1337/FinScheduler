import {createListCollection} from "@chakra-ui/react";
import type {Flex} from "@chakra-ui/react";
import type React from "react";

export const defaultPageSizeOptions = [10, 25, 50, 100] as const;
export const showcasePageSizeOptions = [12, 24, 48, 96] as const;

export function createPageSizeCollection(pageSizeOptions: readonly number[]) {
    return createListCollection({
        items: pageSizeOptions.map((pageSize) => ({
            label: pageSize.toString(),
            value: pageSize,
        })),
    });
}

export interface DataListingColumn<TData> {
    header: string | React.ReactNode;
    key: string;
    render: (row: TData) => React.ReactNode;
    headerProps?: React.ComponentProps<typeof Flex>;
    cellProps?: React.ComponentProps<typeof Flex>;
}
