import {Card, SimpleGrid, Text} from '@chakra-ui/react';
import type {DataListingColumn} from '../../dataListing/types.ts';
import DataShowcaseCard from './DataShowcaseCard.tsx';

interface DataShowcaseGridProps<T> {
    data: T[];
    columns: DataListingColumn<T>[];
    selectable: boolean;
    selectedRows: Set<string>;
    getRowId: (row: T, index: number) => string;
    onSelectRow: (rowId: string, checked: boolean) => void;
    onEdit?: (row: T) => void;
    onRowEdit: (row: T) => void;
}

export default function DataShowcaseGrid<T extends object>(props: DataShowcaseGridProps<T>) {
    const {data, columns, selectable, selectedRows, getRowId, onSelectRow, onEdit, onRowEdit} =
        props;

    if (data.length === 0) {
        return (
            <Card.Root>
                <Card.Body minH="16rem" justifyContent="center" alignItems="center">
                    <Text color="app.accentViolet" fontSize="xl">
                        Данные не найдены.
                    </Text>
                </Card.Body>
            </Card.Root>
        );
    }

    return (
        <SimpleGrid columns={{base: 1, md: 2, xl: 3, '2xl': 4}} gap={4} alignItems="stretch">
            {data.map((row, index) => {
                const rowId = getRowId(row, index);
                const isSelected = selectable && selectedRows.has(rowId);

                return (
                    <DataShowcaseCard
                        key={rowId}
                        row={row}
                        rowId={rowId}
                        columns={columns}
                        selectable={selectable}
                        isSelected={isSelected}
                        onSelectRow={onSelectRow}
                        onEdit={onEdit}
                        onRowEdit={onRowEdit}
                    />
                );
            })}
        </SimpleGrid>
    );
}
