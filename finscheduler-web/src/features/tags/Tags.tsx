import DataTable, {type TableColumn} from "../../components/dataTable/DataTable.tsx";
import {Badge, Flex, Spinner, Text} from "@chakra-ui/react";
import {useState, useEffect, useCallback} from "react";
import type {TagDto, TagFilter} from "../../api/types.ts";
import {toaster} from "../../components/ui/toaster.tsx";
import TagModal from "./subcomponents/TagModal.tsx";
import TagsService, {buildTagFilter, type TagStatusFilter} from "../../api/tags.ts";
import TagsFilters from "./subcomponents/TagsFilters.tsx";

const tagsService = new TagsService();

export default function Tags() {
    const [tags, setTags] = useState<TagDto[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingTag, setEditingTag] = useState<TagDto | null>(null);
    const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState<TagStatusFilter>('Active');

    const itemColumns: TableColumn<TagDto>[] = [
        {
            header: 'Название',
            key: 'name',
            render: (row: TagDto) => (
                <Text fontWeight="semibold" color="neon.blue">
                    {row.name || '-'}
                </Text>
            ),
            headerProps: {textAlign: 'left'},
        },
        {
            header: 'Статус',
            key: 'isActive',
            render: (row: TagDto) => (
                <Badge
                    fontSize="sm"
                    px={3}
                    py={1}
                    borderRadius="full"
                    bg={row.isActive ? "neon.green" : "neon.pink"}
                    color="bg.base"
                >
                    {row.isActive ? "Активен" : "Неактивен"}
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
            const result = await tagsService.getTags(filter);
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

    const handleOpenAddModal = () => {
        setEditingTag(null);
        setModalMode('create');
        setIsModalOpen(true);
    };

    const handleOpenEditModal = (tag: TagDto) => {
        setEditingTag(tag);
        setModalMode('edit');
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEditingTag(null);
        setSelectedRows(new Set());
    };

    const handleSaveTag = async (tagData: Omit<TagDto, 'id'>) => {
        try {
            if (modalMode === 'create') {
                await tagsService.createTag(tagData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Элемент успешно добавлен',
                    type: 'success',
                });
            } else if (modalMode === 'edit' && editingTag?.id) {
                await tagsService.updateTag(editingTag.id, tagData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Элемент успешно обновлен',
                    type: 'success',
                });
            } else {
                console.error('Cannot save: missing id for edit mode', {modalMode, editingItem: editingTag});
                throw new Error('Не удалось определить режим сохранения');
            }

            await loadTags();
        } catch (err) {
            console.error('Error saving tag:', err);
            throw err;
        }
    };

    const getRowId = (row: TagDto): string => {
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
                    <Spinner size="xl" color="neon.blue"/>
                </Flex>
            ) : error ? (
                <Flex justify="center" align="center" minH="400px">
                    <Text color="neon.pink" fontSize="lg">{error}</Text>
                </Flex>
            ) : (
                <>
                    <DataTable
                        data={tags}
                        columns={itemColumns}
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
                        onAdd={handleOpenAddModal}
                        onEdit={handleOpenEditModal}
                        getRowId={getRowId}
                    />
                    <TagModal
                        key={`${modalMode}-${editingTag?.id || 'new'}`}
                        isOpen={isModalOpen}
                        onClose={handleCloseModal}
                        onSave={handleSaveTag}
                        item={editingTag}
                        mode={modalMode}
                    />
                </>
            )}
        </Flex>
    );
}
