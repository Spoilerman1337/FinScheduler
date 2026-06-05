import {Flex, Stack, Text} from "@chakra-ui/react";
import type {ReactNode} from "react";
import type {DataListingColumn} from "../../dataListing/types.ts";

interface DataShowcaseFieldProps<T> {
    row: T;
    column: DataListingColumn<T>;
    isPrimary?: boolean;
}

function renderFieldLabel(header: string | ReactNode) {
    if (typeof header === "string") {
        return (
            <Text
                textStyle="eyebrow"
                color="fg.muted"
            >
                {header}
            </Text>
        );
    }

    return (
        <Flex color="fg.muted" fontSize="xs" textTransform="uppercase">
            {header}
        </Flex>
    );
}

export default function DataShowcaseField<T>(props: DataShowcaseFieldProps<T>) {
    const {row, column, isPrimary = false} = props;

    return (
        <Stack gap={isPrimary ? 2 : 1.5} align="stretch" minW={0}>
            {renderFieldLabel(column.header)}
            <Flex
                minW={0}
                align="flex-start"
                justify="flex-start"
                wrap="wrap"
                {...column.cellProps}
            >
                {column.render(row)}
            </Flex>
        </Stack>
    );
}
