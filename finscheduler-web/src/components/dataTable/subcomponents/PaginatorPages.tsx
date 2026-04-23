import {ButtonGroup, Pagination} from "@chakra-ui/react";
import PaginatorArrowButton from "./PaginatorArrowButton.tsx";
import PaginatorEllipsisButton from "./PaginatorEllipsisButton.tsx";
import PaginatorPageButton from "./PaginatorPageButton.tsx";

interface PaginatorPagesProps {
    totalPages?: number;
    page?: number;
    onPageChange: (page: number) => void;
}

export default function PaginatorPages(props: PaginatorPagesProps) {
    return (
        <Pagination.Root
            count={props.totalPages}
            pageSize={1}
            page={props.page}
            onPageChange={(e) => props.onPageChange(e.page)}
            key={props.totalPages}
        >
            <ButtonGroup variant="ghost" size="lg" my={-5}>
                <Pagination.PrevTrigger asChild>
                    <PaginatorArrowButton direction="previous"/>
                </Pagination.PrevTrigger>

                <Pagination.Items
                    render={(page) => <PaginatorPageButton value={page.value}/>}
                    ellipsis={<PaginatorEllipsisButton/>}
                />

                <Pagination.NextTrigger asChild>
                    <PaginatorArrowButton direction="next"/>
                </Pagination.NextTrigger>
            </ButtonGroup>
        </Pagination.Root>
    );
}
