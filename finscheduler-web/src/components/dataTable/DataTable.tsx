import {Flex} from '@chakra-ui/react';
import DataListingContainer, {
    type DataListingCommonProps,
} from '../dataListing/DataListingContainer.tsx';
import type {DataListingColumn} from '../dataListing/types.ts';
import DataTableHeader from './subcomponents/DataTableHeader.tsx';
import DataTableRows from './subcomponents/DataTableRows.tsx';

export type {DataListingColumn} from '../dataListing/types.ts';

interface DataTableProps<T> extends DataListingCommonProps<T> {
    columns: DataListingColumn<T>[];
}

export default function DataTable<T extends object>(props: DataTableProps<T>) {
    return (
        <Flex direction="column" width="100%">
            <DataListingContainer {...props}>
                {({
                    data,
                    selectable,
                    selectedRows,
                    allSelected,
                    resolveRowId,
                    handleSelectAll,
                    handleSelectRow,
                    handleEdit,
                }) => (
                    <Flex
                        width="100%"
                        bg="bg.layer1"
                        borderColor="borderStrong"
                        boxShadow="card"
                        p={0}
                        background="gradients.cosmic"
                        flexDirection="column"
                        borderRadius="xl"
                        border="1px solid rgba(255,255,255,0.2)"
                        overflow="hidden"
                    >
                        <Flex
                            as="table"
                            width="100%"
                            minWidth="min-content"
                            borderCollapse="collapse"
                            flexDirection="column"
                        >
                            <DataTableHeader
                                columns={props.columns}
                                selectable={selectable}
                                allSelected={allSelected}
                                onSelectAll={handleSelectAll}
                            />
                            <DataTableRows
                                data={data}
                                columns={props.columns}
                                selectable={selectable}
                                selectedRows={selectedRows}
                                getRowId={resolveRowId}
                                onSelectRow={handleSelectRow}
                                onEdit={props.onEdit}
                                onRowEdit={handleEdit}
                            />
                        </Flex>
                    </Flex>
                )}
            </DataListingContainer>
        </Flex>
    );
}
