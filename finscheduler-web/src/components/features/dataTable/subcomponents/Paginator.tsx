import {Flex, Text} from "@chakra-ui/react";
import PaginatorPageSizeSelector from "./PaginatorPageSizeSelector.tsx";
import PaginatorPages from "./PaginatorPages.tsx";

export interface PaginatorProps {
    total: number;
    page: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
}

export default function Paginator(props: PaginatorProps) {
    const totalPages = Math.ceil(props.total / props.pageSize);

    return (
        <Flex justify="space-between" align="center" mt={2} py={5}>
            <PaginatorPageSizeSelector
                pageSize={props.pageSize}
                onPageSizeChange={props.onPageSizeChange}
            />

            <Text fontSize="lg" color="white" fontFamily="body">Элементов всего: {props.total}</Text>

            <PaginatorPages
                totalPages={totalPages}
                page={props.page}
                onPageChange={props.onPageChange}
            />
        </Flex>
    )
}