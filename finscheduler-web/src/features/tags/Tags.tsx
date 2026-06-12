import DataTable, {type DataListingColumn} from '../../components/dataTable/DataTable.tsx';
import {Badge, Flex, Spinner, Text} from '@chakra-ui/react';
import {useCallback, useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import type {TagFilter, TagListingDto, TagStatusFilter} from '../../api/tags.types.ts';
import TagsService, {buildTagFilter} from '../../api/tags.ts';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {buildEditTagPath, newTagPath} from '../routes.ts';
import TagsFilters from './subcomponents/TagsFilters.tsx';

const tagsService = new TagsService();

export default function Tags() {
    const navigate = useNavigate();
    const [tags, setTags] = useState<TagListingDto[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<TagStatusFilter>('Active');

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

    const loadTags = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const filter: TagFilter = buildTagFilter({
                page,
                pageSize,
                searchTerm,
                statusFilter,
            });
            const result = await tagsService.getListingInfo(filter);
            setTags(result.data);
            setTotal(result.count);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка загрузки данных');
            console.error('Failed to load tags:', err);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, searchTerm, statusFilter]);

    useEffect(() => {
        loadTags();
    }, [loadTags]);

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

    const getRowId = (row: TagListingDto): string => {
        return row.id ?? '';
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
                    loadTags();
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
                <>
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
                        getRowId={getRowId}
                    />
                </>
            )}
        </Flex>
    );
}
