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
    CloseButton,
    createListCollection,
    SelectRoot,
    SelectTrigger,
    SelectValueText,
    SelectContent,
    SelectItem,
    Spinner,
    NumberInput,
} from "@chakra-ui/react";
import {useState, useEffect, useMemo} from "react";
import type {ItemDto, ItemModification, Lookup} from "../../../../api/types.ts";
import {categoryTranslations} from "../../../../models/items.ts";
import TagsService from "../../../../api/tags.ts";

interface ItemModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (item: ItemModification) => Promise<void>;
    item?: ItemDto | null;
    mode: 'create' | 'edit';
}

const tagsService = new TagsService();

export default function ItemModal({isOpen, onClose, onSave, item, mode}: ItemModalProps) {
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        price: '',
        cashback: '',
        isActive: true,
        category: 'Не выбрано',
        tags: [] as Lookup[]
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [tagsLoading, setTagsLoading] = useState(false);
    const [tagOptions, setTagOptions] = useState<Lookup[]>([]);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            const newFormData = {
                name: item.name || '',
                description: item.description || '',
                price: item.price !== undefined && item.price !== null ? item.price.toString() : '',
                cashback: item.cashback !== undefined && item.cashback !== null ? item.cashback.toString() : '',
                isActive: item.isActive !== undefined ? item.isActive : true,
                category: item.category !== undefined ? item.category : '',
                tags: item.tags !== undefined ? item.tags : []
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
                category: 'Не выбрано',
                tags: []
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
                category: 'Не выбрано',
                tags: []
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
                price: formData.price ? parseFloat(formData.price) : 0,
                cashback: formData.cashback ? parseFloat(formData.cashback) : 0,
                isActive: formData.isActive,
                category: formData.category,
                tagIds: formData.tags ? formData.tags?.map(x => x.value ?? '').filter(Boolean) : [],
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

    useEffect(() => {
        if (isOpen) {
            const fetchTags = async () => {
                setTagsLoading(true);
                try {
                    const tags = await tagsService.getLookup({
                        page: 0,
                        pageSize: 10,
                    });
                    if (tags && tags.data) {
                        setTagOptions(tags.data);
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

    const tagsCollection = useMemo(() => createListCollection({
        items: tagOptions.map(tag => ({
            label: tag.label,
            value: tag.value,
        })),
    }), [tagOptions]);

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

                            <Field.Root required gap={0}>
                                <Field.Label color="neon.blue">Название</Field.Label>
                                <Input
                                    value={formData.name}
                                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{color: 'textMuted'}}
                                    placeholder="Введите название"
                                />
                            </Field.Root>

                            <Field.Root gap={0}>
                                <Field.Label color="neon.blue">Описание</Field.Label>
                                <Textarea
                                    value={formData.description}
                                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{color: 'textMuted'}}
                                    placeholder="Введите описание"
                                    rows={4}
                                />
                            </Field.Root>

                            <Field.Root gap={0}>
                                <Field.Label color="neon.blue">Цена (₽)</Field.Label>
                                <NumberInput.Root
                                    value={formData.price}
                                    onValueChange={(e) => setFormData({...formData, price: e.value})}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{color: 'textMuted'}}
                                    defaultValue="0.00"
                                    step={0.01}
                                    min={0}
                                    width="100%"
                                >
                                    <NumberInput.Control>
                                        <NumberInput.IncrementTrigger bg="bg.layer2" color="neon.blue"
                                                                      p={0}>+</NumberInput.IncrementTrigger>
                                        <NumberInput.DecrementTrigger bg="bg.layer2" color="neon.blue"
                                                                      p={0}>-</NumberInput.DecrementTrigger>
                                    </NumberInput.Control>
                                    <NumberInput.Input/>
                                </NumberInput.Root>/
                            </Field.Root>

                            <Field.Root gap={0}>
                                <Field.Label color="neon.blue">Кэшбэк (%)</Field.Label>
                                <NumberInput.Root
                                    value={formData.cashback}
                                    onValueChange={(e) => setFormData({...formData, cashback: e.value})}
                                    bg="bg.layer2"
                                    borderColor="glass.border"
                                    color="neon.blue"
                                    _placeholder={{color: 'textMuted'}}
                                    defaultValue="0"
                                    step={1}
                                    min={0}
                                    max={100}
                                    width="100%"
                                    gap={0}
                                >
                                    <NumberInput.Control>
                                        <NumberInput.IncrementTrigger bg="bg.layer2" color="neon.blue"
                                                                      p={0}>+</NumberInput.IncrementTrigger>
                                        <NumberInput.DecrementTrigger bg="bg.layer2" color="neon.blue"
                                                                      p={0}>-</NumberInput.DecrementTrigger>
                                    </NumberInput.Control>
                                    <NumberInput.Input/>
                                </NumberInput.Root>/
                            </Field.Root>

                            <Field.Root required gap={0}>
                                <Field.Label color="neon.blue">Категория</Field.Label>
                                <SelectRoot
                                    collection={categoryList}
                                    value={[formData.category]}
                                    onValueChange={(details) =>
                                        setFormData({...formData, category: details.value[0]})
                                    }
                                >
                                    <SelectTrigger bg="bg.layer2" borderColor="glass.border">
                                        <SelectValueText placeholder="Выберите категорию" color="neon.blue"/>
                                    </SelectTrigger>
                                    <SelectContent bg="bg.layer1" borderColor="glass.border" zIndex="popover">
                                        {categoryList.items.map((cat) => (
                                            <SelectItem
                                                item={cat}
                                                key={cat.value}
                                                _hover={{bg: "bg.layer3"}}
                                                color="white"
                                            >
                                                {cat.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </SelectRoot>
                            </Field.Root>

                            <Field.Root gap={0}>
                                <Field.Label color="neon.blue">
                                    Теги {tagsLoading && <Spinner size="xs" ml={2}/>}
                                </Field.Label>
                                <SelectRoot
                                    multiple
                                    collection={tagsCollection}
                                    value={formData.tags?.map(t => t.value ?? '')}
                                    onValueChange={(details) => {
                                        const selectedTags = tagOptions.filter(opt => details.value.includes(opt.value!));
                                        setFormData({...formData, tags: selectedTags});
                                    }}
                                >
                                    <SelectTrigger bg="bg.layer2" borderColor="glass.border">
                                        <SelectValueText placeholder="Выберите теги" color="neon.blue"/>
                                    </SelectTrigger>
                                    <SelectContent bg="bg.layer1" borderColor="glass.border" zIndex="popover">
                                        {tagsCollection.items.map((tag) => (
                                            <SelectItem
                                                item={tag}
                                                key={tag.value}
                                                _hover={{bg: "bg.layer3"}}
                                                _selected={{color: "neon.blue", fontWeight: "bold"}}
                                                color="white"
                                            >
                                                {tag.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </SelectRoot>
                            </Field.Root>

                            <Field.Root gap={0}>
                                <Flex align="center" gap={2}>
                                    <Field.Label color="neon.blue" mb={0}>
                                        Активен
                                    </Field.Label>
                                    <Switch.Root
                                        checked={formData.isActive}
                                        onCheckedChange={(details) => {
                                            const isChecked = details.checked;
                                            setFormData(prev => ({...prev, isActive: isChecked}));
                                        }}
                                    >
                                        <Switch.HiddenInput/>
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
                                            <Switch.Thumb/>
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
