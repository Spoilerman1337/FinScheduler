import {useEffect, useMemo, useState} from "react";
import type {ItemDto, ItemModification} from "../../../api/types.ts";
import TagsService from "../../../api/tags.ts";
import AsyncSelectField from "../../../components/formFields/AsyncSelectField.tsx";
import NumberField from "../../../components/formFields/NumberField.tsx";
import SelectField from "../../../components/formFields/SelectField.tsx";
import SwitchField from "../../../components/formFields/SwitchField.tsx";
import TextAreaField from "../../../components/formFields/TextAreaField.tsx";
import TextField from "../../../components/formFields/TextField.tsx";
import type {SelectOption} from "../../../components/formFields/types.ts";
import FormModal from "../../../components/ui/FormModal.tsx";
import {categoryOptions, categoryTranslations} from "../../../models/items.ts";

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
    tagIds: string[];
}

const TAGS_PAGE_SIZE = 20;
const tagsService = new TagsService();

function mapItemTagsToOptions(item?: ItemDto | null): SelectOption[] {
    if (!item || !Array.isArray(item.tags)) {
        return [];
    }

    return item.tags.reduce<SelectOption[]>((result, tag) => {
        if (!tag?.value) {
            return result;
        }

        result.push({
            label: tag.label ?? tag.value,
            value: tag.value,
        });

        return result;
    }, []);
}

export default function ItemModal({isOpen, onClose, onSave, item, mode}: ItemModalProps) {
    const getDefaultFormData = (): ItemModalFormData => ({
        name: '',
        description: '',
        price: '',
        cashback: '',
        isActive: true,
        category: '',
        tagIds: [],
    });

    const [formData, setFormData] = useState<ItemModalFormData>(getDefaultFormData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const initialTagOptions = useMemo(() => mapItemTagsToOptions(item), [item]);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            setFormData({
                name: item.name || '',
                description: item.description || '',
                price: item.price !== undefined && item.price !== null ? item.price.toString() : '',
                cashback: item.cashback !== undefined && item.cashback !== null ? item.cashback.toString() : '',
                isActive: typeof item.isActive === 'boolean' ? item.isActive : true,
                category: typeof item.category === 'string' ? item.category : '',
                tagIds: Array.isArray(item.tags)
                    ? item.tags.reduce<string[]>((result, tag) => {
                        if (tag?.value) {
                            result.push(tag.value);
                        }

                        return result;
                    }, [])
                    : [],
            });
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

        if (formData.price && Number.isNaN(parseFloat(formData.price))) {
            setError('Цена должна быть числом');
            return;
        }

        if (
            formData.cashback &&
            (Number.isNaN(parseFloat(formData.cashback)) ||
                parseFloat(formData.cashback) < 0 ||
                parseFloat(formData.cashback) > 100)
        ) {
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
                tagIds: formData.tagIds,
            });
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setLoading(false);
        }
    };

    return (
        <FormModal
            isOpen={isOpen}
            onClose={onClose}
            onSubmit={handleSubmit}
            title={mode === 'create' ? 'Добавить новый элемент' : 'Редактировать элемент'}
            error={error}
            loading={loading}
        >
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
            <AsyncSelectField
                multiple
                label="Теги"
                value={formData.tagIds}
                initialOptions={initialTagOptions}
                cacheKey="item-tags"
                placeholder="Выберите теги"
                emptyText="Теги не найдены"
                collapseThreshold={4}
                loadOptions={async ({page, search}) => {
                    const tags = await tagsService.getLookup({
                        page,
                        pageSize: TAGS_PAGE_SIZE,
                        name: search || undefined,
                    });

                    return {
                        options: (tags.data ?? []).reduce<SelectOption[]>((result, tag) => {
                            if (!tag?.value) {
                                return result;
                            }

                            result.push({
                                label: tag.label ?? tag.value,
                                value: tag.value,
                            });

                            return result;
                        }, []),
                        hasMore: (page + 1) * TAGS_PAGE_SIZE < tags.count,
                    };
                }}
                onChange={(value) => updateFormData('tagIds', value)}
            />
            <SwitchField
                label="Активен"
                checked={formData.isActive}
                onChange={(value) => updateFormData('isActive', value)}
            />
        </FormModal>
    );
}
