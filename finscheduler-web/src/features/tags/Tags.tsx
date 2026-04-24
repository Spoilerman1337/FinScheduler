import DataTable, {type TableColumn} from "../../components/dataTable/DataTable.tsx";
import {Badge, Box, Button, Flex, Input, Text, Spinner} from "@chakra-ui/react";
import {SearchIcon} from "lucide-react";
import {useState, useEffect, useCallback} from "react";
import type {TagDto, TagFilter} from "../../api/types.ts";
import {toaster} from "../../components/ui/toaster.tsx";
import TagModal from "./subcomponents/TagModal.tsx";
import TagsService from "../../api/tags.ts";

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
    const [statusFilter, setStatusFilter] = useState<'All' | 'Active' | 'Inactive'>('Active');

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
            headerProps: {textAlign: 'left'}
        },
    ];

    const loadTags = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const filter: TagFilter = {
                page: page - 1,
                pageSize: pageSize,
            };

            if (searchTerm) {
                filter.name = searchTerm;
            }

            if (statusFilter !== 'All') {
                filter.isActive = statusFilter === 'Active';
            }

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
    }

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

    const handleSaveItem = async (tagData: Omit<TagDto, 'id'>) => {
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
            console.error('Error saving item:', err);
            throw err;
        }
    };

    const getRowId = (row: TagDto): string => {
        return row.id ?? '';
    };

    const filterWidthProps = {
        w: {base: '100%', md: 'calc(50% - var(--chakra-space-2))', xl: 'calc(20% - var(--chakra-space-3))'},
    };

    const statusOptions: Array<'All' | 'Active' | 'Inactive'> = ['All', 'Active', 'Inactive'];

    return <Flex direction="column" width="100%">
            <Flex
                mb={4}
                p={4}
                bg="bg.layer2"
                borderRadius="lg"
                border="1px solid"
                borderColor="glass.border"
                width="100%"
                align="center"
                gap={4}
                flexWrap="wrap"
                justifyContent="flex-start"
            >
                <Box {...filterWidthProps} position="relative">
                    <Box position="absolute" left="3" top="50%" transform="translateY(-50%)" zIndex="1"
                         pointerEvents="none">
                        <SearchIcon size={18} color="rgba(255,255,255,0.6)"/>
                    </Box>
                    <Input
                        placeholder="Поиск по названию..."
                        value={searchTerm}
                        onChange={(e) => {
                            setPage(1)
                            setSearchTerm(e.target.value)
                        }}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                setPage(1);
                                loadTags();
                            }
                        }}
                        pl="10"
                        bg="bg.layer1"
                        borderColor="glass.border"
                        color="neon.blue"
                        _placeholder={{color: 'textMuted'}}
                    />
                </Box>

                <Box {...filterWidthProps}>
                    <Flex gap={1} borderRadius="md" p={1} bg="bg.layer1" border="1px solid" borderColor="glass.border">
                        {statusOptions.map(status => (
                            <Button
                                key={status}
                                size="sm"
                                flex={1}
                                onClick={() => {
                                    setStatusFilter(status);
                                    setPage(1);
                                }}
                                color={statusFilter === status ? "neon.blue" : "neon.blue"}
                                borderColor={statusFilter === status ? "neon.blue" : "glass.border"}
                                backdropFilter="blur(12px)"
                                bg={statusFilter === status ? "glass.bgHover" : "transparent"}
                                border="1px solid"
                                transition="all 0.3s ease-in-out"
                                filter={statusFilter === status ? "drop-shadow(0 0 16px rgba(0,212,255,0.9))" : "none"}
                                boxShadow={statusFilter === status ? "0 0 20px rgba(0,212,255,1)" : "none"}
                                _hover={{
                                    filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                    boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                    backdropFilter: "blur(12px)",
                                    borderColor: "neon.purple",
                                }}
                                _active={{
                                    filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                                    boxShadow: "0 0 20px rgba(0,212,255,1)",
                                    color: "neon.blue",
                                    bg: "glass.bgHover",
                                    backdropFilter: "blur(12px)",
                                    borderColor: "neon.blue",
                                }}
                                focusRing="none"
                            >
                                {status === 'All' ? 'Все' : status === 'Active' ? 'Активные' : 'Неактивные'}
                            </Button>
                        ))}
                    </Flex>
                </Box>

                <Button
                    {...filterWidthProps}
                    onClick={handleReset}
                    bg="bg.layer3"
                    color="textMuted"
                    borderColor="glass.border"
                    border="1px solid"
                    _hover={{bg: 'neon.pink', color: 'bg.base'}}
                >
                    Сброс
                </Button>
            </Flex>

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
                        onSave={handleSaveItem}
                        item={editingTag}
                        mode={modalMode}
                    />
                </>
            )}
        </Flex>
}
