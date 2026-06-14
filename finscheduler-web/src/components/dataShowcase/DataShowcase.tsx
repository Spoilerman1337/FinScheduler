import {Card, Flex, Stack, Text} from '@chakra-ui/react';
import DataListingContainer, {
    type DataListingCommonProps,
} from '../dataListing/DataListingContainer.tsx';
import {showcasePageSizeOptions, type DataListingColumn} from '../dataListing/types.ts';
import ListingSelectionCheckbox from '../dataListing/selectionCheckbox/ListingSelectionCheckbox.tsx';
import DataShowcaseGrid from './subcomponents/DataShowcaseGrid.tsx';

export type {DataListingColumn} from '../dataListing/types.ts';

interface DataShowcaseProps<T> extends DataListingCommonProps<T> {
    columns: DataListingColumn<T>[];
}

export default function DataShowcase<T extends object>(props: DataShowcaseProps<T>) {
    return (
        <Stack width="100%" gap={4}>
            <DataListingContainer
                {...props}
                pageSizeOptions={props.pageSizeOptions ?? showcasePageSizeOptions}
                showEditSelectedAction={false}
            >
                {({
                    data,
                    selectable,
                    selectedRows,
                    selectedIds,
                    allSelected,
                    resolveRowId,
                    handleSelectAll,
                    handleSelectRow,
                    handleEdit,
                }) => (
                    <>
                        {selectable && data.length > 0 && (
                            <Card.Root>
                                <Card.Body py="4">
                                    <Flex align="center" gap={3}>
                                        <ListingSelectionCheckbox
                                            checked={allSelected}
                                            onCheckedChange={handleSelectAll}
                                            containerProps={{
                                                as: 'div',
                                                py: 0,
                                                px: 0,
                                                flexBasis: 'auto',
                                                minWidth: 'auto',
                                                maxWidth: 'none',
                                            }}
                                        />
                                        <Stack gap={0}>
                                            <Text color="app.accent" fontWeight="semibold">
                                                {allSelected
                                                    ? 'Все карточки на странице выбраны'
                                                    : 'Выделить все карточки на странице'}
                                            </Text>
                                            <Text color="fg.muted" fontSize="sm">
                                                Выбрано: {selectedIds.length} из {data.length}
                                            </Text>
                                        </Stack>
                                    </Flex>
                                </Card.Body>
                            </Card.Root>
                        )}

                        <DataShowcaseGrid
                            data={data}
                            columns={props.columns}
                            selectable={selectable}
                            selectedRows={selectedRows}
                            getRowId={resolveRowId}
                            onSelectRow={handleSelectRow}
                            onEdit={props.onEdit}
                            onRowEdit={handleEdit}
                        />
                    </>
                )}
            </DataListingContainer>
        </Stack>
    );
}
