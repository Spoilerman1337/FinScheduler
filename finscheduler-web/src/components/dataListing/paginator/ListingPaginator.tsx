import {Flex, Text} from '@chakra-ui/react';
import ListingPaginatorPages from './ListingPaginatorPages.tsx';
import ListingPaginatorPageSizeSelector from './ListingPaginatorPageSizeSelector.tsx';

export interface ListingPaginatorProps {
    total: number;
    page: number;
    pageSize: number;
    pageSizeOptions?: readonly number[];
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
}

export default function ListingPaginator(props: ListingPaginatorProps) {
    const totalPages = Math.ceil(props.total / props.pageSize);

    return (
        <Flex justify="space-between" align="center" mt={2} py={5}>
            <ListingPaginatorPageSizeSelector
                pageSize={props.pageSize}
                pageSizeOptions={props.pageSizeOptions}
                onPageSizeChange={props.onPageSizeChange}
            />

            <Text fontSize="lg" color="white" fontFamily="body">
                Элементов всего: {props.total}
            </Text>

            <ListingPaginatorPages
                totalPages={totalPages}
                page={props.page}
                onPageChange={props.onPageChange}
            />
        </Flex>
    );
}
