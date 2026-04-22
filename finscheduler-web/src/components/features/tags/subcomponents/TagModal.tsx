import {
    Dialog,
    Button,
    Text,
    Stack,
    CloseButton,
} from "@chakra-ui/react";
import { useState, useEffect } from "react";
import type { ItemDto } from "../../../../api/types.ts";
import SwitchField from "../../formFields/SwitchField.tsx";
import TextField from "../../formFields/TextField.tsx";

interface TagModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (tag: Omit<ItemDto, 'id'>) => Promise<void>;
    item?: ItemDto | null;
    mode: 'create' | 'edit';
}

export default function TagModal({ isOpen, onClose, onSave, item, mode }: TagModalProps) {
    const [formData, setFormData] = useState({
        name: '',
        isActive: true
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            const newFormData = {
                name: item.name || '',
                isActive: item.isActive !== undefined ? item.isActive : true,
            };
            setFormData(newFormData);
        }
    }, [isOpen, item, mode]);

    useEffect(() => {
        if (isOpen && mode === 'create') {
            setFormData({
                name: '',
                isActive: true,
            });
            setError(null);
        }
    }, [isOpen, mode]);

    useEffect(() => {
        if (!isOpen) {
            setFormData({
                name: '',
                isActive: true,
            });
            setError(null);
        }
    }, [isOpen]);

    const handleSubmit = async () => {
        setError(null);

        if (!formData.name.trim()) {
            setError('Название обязательно для заполнения');
            return;
        }

        setLoading(true);
        try {
            await onSave({
                name: formData.name.trim(),
                isActive: formData.isActive,
            });
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog.Root open={isOpen} onOpenChange={(details) => !details.open && onClose()} placement="center">
            <Dialog.Backdrop />
            <Dialog.Positioner>
                <Dialog.Content bg="bg.layer1" border="1px solid" borderColor="glass.border" maxW="600px">
                    <Dialog.Header>
                        <Dialog.Title color="neon.blue">
                            {mode === 'create' ? 'Добавить новый элемент' : 'Редактировать элемент'}
                        </Dialog.Title>
                        <Dialog.CloseTrigger asChild>
                            <CloseButton />
                        </Dialog.CloseTrigger>
                    </Dialog.Header>
                    <Dialog.Body>
                        <Stack gap={4}>
                            {error && (
                                <Text color="neon.pink" fontSize="sm">
                                    {error}
                                </Text>
                            )}

                            <TextField
                                label="Название"
                                value={formData.name}
                                placeholder="Введите название"
                                required
                                onChange={(value) => setFormData({...formData, name: value})}
                            />

                            <SwitchField
                                label="Активен"
                                checked={formData.isActive}
                                onChange={(value) => setFormData(prev => ({...prev, isActive: value}))}
                            />
                        </Stack>
                    </Dialog.Body>

                    <Dialog.Footer>
                        <Button
                            variant="ghost"
                            mr={3}
                            onClick={onClose}
                            color="textMuted"
                            _hover={{ bg: 'bg.layer2' }}
                        >
                            Отмена
                        </Button>
                        <Button
                            bg="neon.blue"
                            color="bg.base"
                            onClick={handleSubmit}
                            loading={loading}
                            _hover={{ bg: 'neon.blue', opacity: 0.8 }}
                        >
                            Сохранить
                        </Button>
                    </Dialog.Footer>
                </Dialog.Content>
            </Dialog.Positioner>
        </Dialog.Root>
    );
}
