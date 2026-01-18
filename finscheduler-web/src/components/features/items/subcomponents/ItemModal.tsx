import {
    Dialog,
    Button,
    Input,
    Textarea,
    Flex,
    Text,
    Field,
    Switch,
    Stack,
    CloseButton, createListCollection, SelectRoot, SelectTrigger, SelectValueText, SelectContent, SelectItem,
} from "@chakra-ui/react";
import { useState, useEffect } from "react";
import type { ItemDto } from "../../../../api/types.ts";
import {categoryTranslations} from "../../../../models/items.ts";

interface ItemModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (item: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
    item?: ItemDto | null;
    mode: 'create' | 'edit';
}

export default function ItemModal({ isOpen, onClose, onSave, item, mode }: ItemModalProps) {
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        price: '',
        cashback: '',
        isActive: true,
        category: 'Не выбрано'
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            const newFormData = {
                name: item.name || '',
                description: item.description || '',
                price: item.price !== undefined && item.price !== null ? item.price.toString() : '',
                cashback: item.cashback !== undefined && item.cashback !== null ? item.cashback.toString() : '',
                isActive: item.isActive !== undefined ? item.isActive : true,
                category: item.category !== undefined ? item.category : ''
            };
            setFormData(newFormData);
        }
    }, [isOpen, item, mode]);

    useEffect(() => {
        if (isOpen && mode === 'create') {
            setFormData({
                name: '',
                description: '',
                price: '',
                cashback: '',
                isActive: true,
                category: 'Не выбрано'
            });
            setError(null);
        }
    }, [isOpen, mode]);

    useEffect(() => {
        if (!isOpen) {
            setFormData({
                name: '',
                description: '',
                price: '',
                cashback: '',
                isActive: true,
                category: 'Не выбрано'
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

        if (formData.price && isNaN(parseFloat(formData.price))) {
            setError('Цена должна быть числом');
            return;
        }

        if (formData.cashback && (isNaN(parseFloat(formData.cashback)) || parseFloat(formData.cashback) < 0 || parseFloat(formData.cashback) > 100)) {
            setError('Кэшбэк должен быть числом от 0 до 100');
            return;
        }

        if (!formData.category.trim() || !categoryTranslations[formData.category]) {
            setError('Выберите категорию');
            return;
        }

        setLoading(true);
        try {
            await onSave({
                name: formData.name.trim(),
                description: formData.description.trim() || undefined,
                price: formData.price ? parseFloat(formData.price) : undefined,
                cashback: formData.cashback ? parseFloat(formData.cashback) : undefined,
                isActive: formData.isActive,
                category: formData.category
            });
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setLoading(false);
        }
    };

    const categoryList = createListCollection({
        items: Object.entries(categoryTranslations).map(([value, label]) => ({
            label,
            value,
        })),
    });

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

                            <Field.Root required>
                                <Field.Label color="neon.blue">Название</Field.Label>
                                <Input
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{ color: 'textMuted' }}
                                    placeholder="Введите название"
                                />
                            </Field.Root>

                            <Field.Root>
                                <Field.Label color="neon.blue">Описание</Field.Label>
                                <Textarea
                                    value={formData.description}
                                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{ color: 'textMuted' }}
                                    placeholder="Введите описание"
                                    rows={4}
                                />
                            </Field.Root>

                            <Field.Root>
                                <Field.Label color="neon.blue">Цена (₽)</Field.Label>
                                <Input
                                    type="number"
                                    value={formData.price}
                                    onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{ color: 'textMuted' }}
                                    placeholder="0.00"
                                    step="0.01"
                                    min="0"
                                />
                            </Field.Root>

                            <Field.Root>
                                <Field.Label color="neon.blue">Кэшбэк (%)</Field.Label>
                                <Input
                                    type="number"
                                    value={formData.cashback}
                                    onChange={(e) => setFormData({ ...formData, cashback: e.target.value })}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{ color: 'textMuted' }}
                                    placeholder="0"
                                    step="0.1"
                                    min="0"
                                    max="100"
                                />
                            </Field.Root>

                            <Field.Root required>
                                <Field.Label color="neon.blue">Категория</Field.Label>
                                <SelectRoot
                                    collection={categoryList}
                                    value={[formData.category]}
                                    onValueChange={(details) =>
                                        setFormData({ ...formData, category: details.value[0] })
                                    }
                                >
                                    <SelectTrigger bg="bg.layer2" borderColor="glass.border">
                                        <SelectValueText placeholder="Выберите категорию" color="neon.blue" />
                                    </SelectTrigger>
                                    <SelectContent bg="bg.layer1" borderColor="glass.border" zIndex="popover">
                                        {categoryList.items.map((cat) => (
                                            <SelectItem
                                                item={cat}
                                                key={cat.value}
                                                _hover={{ bg: "bg.layer3" }}
                                                color="white"
                                            >
                                                {cat.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </SelectRoot>
                            </Field.Root>

                            <Field.Root>
                                <Flex align="center" gap={2}>
                                    <Field.Label color="neon.blue" mb={0}>
                                        Активен
                                    </Field.Label>
                                    <Switch.Root
                                        checked={formData.isActive}
                                        onCheckedChange={(details) => {
                                            const isChecked = details.checked;
                                            setFormData(prev => ({ ...prev, isActive: isChecked }));
                                        }}
                                    >
                                        <Switch.HiddenInput />
                                        <Switch.Control
                                            bg={formData.isActive ? "neon.blue" : "neon.purple"}
                                            filter={formData.isActive 
                                                ? "drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))" 
                                                : "drop-shadow(0 0 8px rgba(212, 0, 255, 0.9))"}
                                            boxShadow={formData.isActive
                                                ? "0 0 12px rgba(0, 212, 255, 0.6)"
                                                : "0 0 12px rgba(212, 0, 255, 0.6)"}
                                            transition="all 0.3s ease-in-out"
                                        >
                                            <Switch.Thumb />
                                        </Switch.Control>
                                    </Switch.Root>
                                </Flex>
                            </Field.Root>
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
