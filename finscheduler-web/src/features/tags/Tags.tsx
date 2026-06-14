import DataTable, {type DataListingColumn} from '../../components/dataTable/DataTable.tsx';
import {Badge, Flex, Spinner, Text} from '@chakra-ui/react';
import {useMemo, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import type {TagFilter, TagListingDto, TagStatusFilter} from '../../api/tags.types.ts';
import {buildTagFilter} from '../../api/tags.ts';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {buildEditTagPath, newTagPath} from '../routes.ts';
import {useTagsListQuery} from './queries.ts';
import TagsFilters from './subcomponents/TagsFilters.tsx';

export default function Tags() {
    const navigate = useNavigate();
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<TagStatusFilter>('Active');

    const filter: TagFilter = useMemo(
        () =>
            buildTagFilter({
                page,
                pageSize,
                searchTerm,
                statusFilter,
            }),
        [page, pageSize, searchTerm, statusFilter],
    );
    const tagsQuery = useTagsListQuery(filter);
    const tags = tagsQuery.data?.data ?? [];
    const total = tagsQuery.data?.count ?? 0;
    const loading = tagsQuery.isPending;
    const error =
        tagsQuery.isError && tagsQuery.error instanceof Error
            ? tagsQuery.error.message
            : tagsQuery.isError
              ? 'Ошибка загрузки данных'
              : null;

    const tagColumns: DataListingColumn<TagListingDto>[] = [
        {
            header: 'Название',
            key: 'name',
            render: (row: TagListingDto) => (
                <Text fontWeight="semibold" color="neon.blue">
                    {row.name || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Статус',
            key: 'isActive',
            render: (row: TagListingDto) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.isActive ? 'neon.green' : 'neon.pink'}
                    color="bg.base"
                >
                    {row.isActive ? 'Активен' : 'Неактивен'}
                </Badge>
            ),
            headerProps: {textAlign: 'left'},
        },
    ];

    const handleReset = () => {
        setSearchTerm('');
        setStatusFilter('All');
        setPage(1);
    };

    const handleOpenCreatePage = () => {
        navigate(newTagPath);
    };

    const handleOpenEditPage = (tag: TagListingDto) => {
        if (!tag.id) {
            toaster.create({
                title: 'Ошибка',
                description: 'Не удалось открыть карточку тега',
                type: 'error',
            });
            return;
        }

        navigate(buildEditTagPath(tag.id));
    };

    return (
        <Flex direction="column" width="100%">
            <TagsFilters
                searchTerm={searchTerm}
                statusFilter={statusFilter}
                onSearchTermChange={(value) => {
                    setPage(1);
                    setSearchTerm(value);
                }}
                onStatusFilterChange={(value) => {
                    setStatusFilter(value);
                    setPage(1);
                }}
                onApply={() => {
                    setPage(1);
                }}
                onReset={handleReset}
            />

            {loading ? (
                <Flex justify="center" align="center" minH="400px">
                    <Spinner size="xl" color="neon.blue" />
                </Flex>
            ) : error ? (
                <Flex justify="center" align="center" minH="400px">
                    <Text color="neon.pink" fontSize="lg">
                        {error}
                    </Text>
                </Flex>
            ) : (
                <DataTable
                    data={tags}
                    columns={tagColumns}
                    total={total}
                    page={page}
                    pageSize={pageSize}
                    onPageChange={setPage}
                    onPageSizeChange={(newSize) => {
                        setPageSize(newSize);
                        setPage(1);
                    }}
                    selectable={true}
                    selectedRows={selectedRows}
                    onSelectionChange={setSelectedRows}
                    onAdd={handleOpenCreatePage}
                    onEdit={handleOpenEditPage}
                    getRowId={(row) => row.id}
                />
            )}
        </Flex>
    );
}
