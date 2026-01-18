import {Button, Flex} from "@chakra-ui/react";
import {PencilIcon, PlusIcon, TrashIcon} from "lucide-react";

interface DataTableActionButtonsProps<T> {
    data: T[];
    selectedRows?: Set<string>;
    onAdd?: () => void;
    onEdit?: (row: T) => void;
    onDelete?: (ids: string[]) => void;
    getRowId?: (row: T) => string;
    handleEdit: (row: T) => void;
}

export default function DataTableActionButtons<T extends object>(props: DataTableActionButtonsProps<T>) {
    const {
        data = [],
        selectedRows = new Set(),
        onAdd,
        onEdit,
        onDelete,
        getRowId,
        handleEdit
    } = props;

    const handleDelete = () => {
        if (onDelete && selectedRows.size > 0) {
            onDelete(Array.from(selectedRows));
        }
    };

    return (
        <Flex
            mt={2}
            p={4}
            bg="bg.layer2"
            borderRadius="lg"
            border="1px solid"
            borderColor="glass.border"
            justify="flex-end"
            gap={2}
        >
            {onAdd && (
                <Button
                    onClick={onAdd}
                    bg="neon.blue"
                    color="bg.base"
                    _hover={{bg: 'neon.blue', opacity: 0.8}}
                >
                    <PlusIcon size={18} style={{marginRight: '8px'}}/>
                    Добавить
                </Button>
            )}
            {onEdit && selectedRows.size === 1 && (
                <Button
                    onClick={() => {
                        const selectedId = Array.from(selectedRows)[0];
                        const selectedRow = data.find(row => {
                            const rowId = getRowId ? getRowId(row) : String(row['id' as keyof T] || '');
                            return rowId === selectedId;
                        });
                        if (selectedRow) {
                            handleEdit(selectedRow);
                        }
                    }}
                    bg="neon.green"
                    color="bg.base"
                    _hover={{bg: 'neon.green', opacity: 0.8}}
                    disabled={selectedRows.size !== 1}
                >
                    <PencilIcon size={18} style={{marginRight: '8px'}}/>
                    Редактировать
                </Button>
            )}
            {onDelete && (
                <Button
                    onClick={handleDelete}
                    color="neon.blue"
                    borderColor="neon.blue"
                    backdropFilter="blur(12px)"
                    bg="glass.bgHover"
                    border="1px solid"
                    transition="all 0.3s ease-in-out"
                    disabled={selectedRows.size === 0}
                    opacity={selectedRows.size === 0 ? 0.5 : 1}
                    _hover={selectedRows.size > 0 ? {
                        filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                        boxShadow: "0 0 20px rgba(212, 0,255,1)",
                        color: "neon.purple",
                        bg: "glass.bgHover",
                        backdropFilter: "blur(12px)",
                        borderColor: "neon.purple",
                    } : {}}
                    focusRing="none"
                >
                    <TrashIcon size={18} style={{marginRight: '8px'}}/>
                    Удалить {selectedRows.size > 0 ? `(${selectedRows.size})` : ''}
                </Button>
            )}
        </Flex>
    )
}