import {
    Dialog,
    Button,
    Text,
    Stack,
    CloseButton,
} from "@chakra-ui/react";
import {useState, useEffect, useMemo} from "react";
import type {ItemDto, ItemModification, Lookup} from "../../../../api/types.ts";
import {categoryOptions, categoryTranslations} from "../../../../models/items.ts";
import TagsService from "../../../../api/tags.ts";
import NumberField from "../../formFields/NumberField.tsx";
import SelectField from "../../formFields/SelectField.tsx";
import SwitchField from "../../formFields/SwitchField.tsx";
import TextAreaField from "../../formFields/TextAreaField.tsx";
import TextField from "../../formFields/TextField.tsx";
import type {SelectOption} from "../../formFields/types.ts";

interface ItemModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (item: ItemModification) => Promise<void>;
    item?: ItemDto | null;
    mode: 'create' | 'edit';
}

interface ItemModalFormData {
    name: string;
    description: string;
    price: string;
    cashback: string;
    isActive: boolean;
    category: string;
    tags: Lookup[];
}

const tagsService = new TagsService();

const getDefaultFormData = (): ItemModalFormData => ({
    name: '',
    description: '',
    price: '',
    cashback: '',
    isActive: true,
    category: 'Не выбрано',
    tags: [],
});

export default function ItemModal({isOpen, onClose, onSave, item, mode}: ItemModalProps) {
    const [formData, setFormData] = useState<ItemModalFormData>(getDefaultFormData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [tagsLoading, setTagsLoading] = useState(false);
    const [tagOptions, setTagOptions] = useState<Lookup[]>([]);

    const tagSelectOptions = useMemo<SelectOption[]>(() => (tagOptions ?? [])
        .filter((tag): tag is Lookup & { value: string } => Boolean(tag.value))
        .map((tag) => ({
            label: tag.label ?? tag.value,
            value: tag.value,
        })), [tagOptions]);

    const selectedTagValues = useMemo(() => (formData.tags ?? [])
        .map((tag) => tag.value)
        .filter((value): value is string => Boolean(value)), [formData.tags]);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            const newFormData: ItemModalFormData = {
                name: item.name || '',
                description: item.description || '',
                price: item.price !== undefined && item.price !== null ? item.price.toString() : '',
                cashback: item.cashback !== undefined && item.cashback !== null ? item.cashback.toString() : '',
                isActive: item.isActive !== undefined ? item.isActive : true,
                category: item.category !== undefined ? item.category : '',
                tags: Array.isArray(item.tags) ? item.tags : []
            };
            setFormData(newFormData);
        }
    }, [isOpen, item, mode]);

    useEffect(() => {
        if (isOpen && mode === 'create') {
            setFormData(getDefaultFormData());
            setError(null);
        }
    }, [isOpen, mode]);

    useEffect(() => {
        if (!isOpen) {
            setFormData(getDefaultFormData());
            setError(null);
        }
    }, [isOpen]);

    useEffect(() => {
        if (isOpen) {
            const fetchTags = async () => {
                setTagsLoading(true);
                try {
                    const tags = await tagsService.getLookup({
                        page: 0,
                        pageSize: 10,
                    });
                    if (tags) {
                        setTagOptions(Array.isArray(tags.data) ? tags.data : []);
                    }

                } catch (e) {
                    console.error("Ошибка загрузки тегов", e);
                } finally {
                    setTagsLoading(false);
                }
            };
            fetchTags();
        }
    }, [isOpen]);

    const updateFormData = <K extends keyof ItemModalFormData>(
        field: K,
        value: ItemModalFormData[K],
    ) => {
        setFormData((prev) => ({...prev, [field]: value}));
    };

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
                price: formData.price ? parseFloat(formData.price) : 0,
                cashback: formData.cashback ? parseFloat(formData.cashback) : 0,
                isActive: formData.isActive,
                category: formData.category,
                tagIds: (formData.tags ?? []).map(x => x.value ?? '').filter(Boolean),
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
            <Dialog.Backdrop/>
            <Dialog.Positioner>
                <Dialog.Content bg="bg.layer1" border="1px solid" borderColor="glass.border" maxW="600px">
                    <Dialog.Header>
                        <Dialog.Title color="neon.blue">
                            {mode === 'create' ? 'Добавить новый элемент' : 'Редактировать элемент'}
                        </Dialog.Title>
                        <Dialog.CloseTrigger asChild bg="bg.layer1" border="1px solid" borderColor="neon.blue">
                            <CloseButton color="neon.blue" filter="drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"
                                         boxShadow="0 0 12px rgba(0, 212, 255, 0.6)"/>
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
                                onChange={(value) => updateFormData('name', value)}
                            />
                            <TextAreaField
                                label="Описание"
                                value={formData.description}
                                placeholder="Введите описание"
                                rows={4}
                                onChange={(value) => updateFormData('description', value)}
                            />
                            <NumberField
                                label="Цена (₽)"
                                value={formData.price}
                                defaultValue="0.00"
                                step={0.01}
                                min={0}
                                onChange={(value) => updateFormData('price', value)}
                            />
                            <NumberField
                                label="Кэшбэк (%)"
                                value={formData.cashback}
                                defaultValue="0"
                                step={1}
                                min={0}
                                max={100}
                                onChange={(value) => updateFormData('cashback', value)}
                            />
                            <SelectField
                                label="Категория"
                                value={formData.category}
                                options={categoryOptions}
                                placeholder="Выберите категорию"
                                required
                                onChange={(value) => updateFormData('category', value)}
                            />
                            <SelectField
                                multiple
                                label="Теги"
                                value={selectedTagValues}
                                options={tagSelectOptions}
                                placeholder="Выберите теги"
                                loading={tagsLoading}
                                onChange={(values) => {
                                    const selectedTags = tagOptions.filter((tag) =>
                                        tag.value ? values.includes(tag.value) : false,
                                    );
                                    updateFormData('tags', selectedTags);
                                }}
                            />
                            <SwitchField
                                label="Активен"
                                checked={formData.isActive}
                                onChange={(value) => updateFormData('isActive', value)}
                            />
                        </Stack>
                    </Dialog.Body>

                    <Dialog.Footer>
                        <Button
                            variant="ghost"
                            mr={3}
                            onClick={onClose}
                            color="textMuted"
                            _hover={{bg: 'bg.layer2'}}
                        >
                            Отмена
                        </Button>
                        <Button
                            bg="neon.blue"
                            color="bg.base"
                            onClick={handleSubmit}
                            loading={loading}
                            _hover={{bg: 'neon.blue', opacity: 0.8}}
                        >
                            Сохранить
                        </Button>
                    </Dialog.Footer>
                </Dialog.Content>
            </Dialog.Positioner>
        </Dialog.Root>
    );
}
