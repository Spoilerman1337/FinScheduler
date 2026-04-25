import {createListCollection} from "@chakra-ui/react";
import type {Flex} from "@chakra-ui/react";
import type React from "react";

export const pageSizes = createListCollection({
    items: [
        {label: "10", value: 10},
        {label: "25", value: 25},
        {label: "50", value: 50},
        {label: "100", value: 100},
    ],
});

export interface TableColumn<TData> {
    header: string | React.ReactNode;
    key: string;
    render: (row: TData) => React.ReactNode;
    headerProps?: React.ComponentProps<typeof Flex>;
    cellProps?: React.ComponentProps<typeof Flex>;
}
