import {ButtonGroup, Pagination} from '@chakra-ui/react';
import ListingPaginatorArrowButton from './ListingPaginatorArrowButton.tsx';
import ListingPaginatorEllipsisButton from './ListingPaginatorEllipsisButton.tsx';
import ListingPaginatorPageButton from './ListingPaginatorPageButton.tsx';

interface ListingPaginatorPagesProps {
    totalPages?: number;
    page?: number;
    onPageChange: (page: number) => void;
}

export default function ListingPaginatorPages(props: ListingPaginatorPagesProps) {
    return (
        <Pagination.Root
            count={props.totalPages}
            pageSize={1}
            page={props.page}
            onPageChange={(event) => props.onPageChange(event.page)}
            key={props.totalPages}
        >
            <ButtonGroup variant="ghost" size="lg" my={-5}>
                <Pagination.PrevTrigger asChild>
                    <ListingPaginatorArrowButton direction="previous" />
                </Pagination.PrevTrigger>

                <Pagination.Items
                    render={(page) => <ListingPaginatorPageButton value={page.value} />}
                    ellipsis={<ListingPaginatorEllipsisButton />}
                />

                <Pagination.NextTrigger asChild>
                    <ListingPaginatorArrowButton direction="next" />
                </Pagination.NextTrigger>
            </ButtonGroup>
        </Pagination.Root>
    );
}
