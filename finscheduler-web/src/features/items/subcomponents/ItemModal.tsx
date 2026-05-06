import {useEffect, useMemo, useState} from "react";
import type {ItemDto, ItemModification} from "../../../api/types.ts";
import TagsService from "../../../api/tags.ts";
import AsyncSelectField from "../../../components/formFields/AsyncSelectField.tsx";
import NumberField from "../../../components/formFields/NumberField.tsx";
import SelectField from "../../../components/formFields/SelectField.tsx";
import SwitchField from "../../../components/formFields/SwitchField.tsx";
import TextAreaField from "../../../components/formFields/TextAreaField.tsx";
import TextField from "../../../components/formFields/TextField.tsx";
import FormModal from "../../../components/formModal/FormModal.tsx";
import {categoryOptions} from "../../../models/items.ts";
import {mapLookupsToSelectOptions} from "../../shared.ts";
import {
    buildItemModification,
    createDefaultItemFormData,
    mapItemToFormData,
    type ItemModalFormData,
    validateItemFormData,
} from "../form.ts";

interface ItemModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (item: ItemModification) => Promise<void>;
    item?: ItemDto | null;
    mode: 'create' | 'edit';
}

const TAGS_PAGE_SIZE = 20;
const tagsService = new TagsService();

export default function ItemModal({isOpen, onClose, onSave, item, mode}: ItemModalProps) {
    const [formData, setFormData] = useState<ItemModalFormData>(createDefaultItemFormData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const initialTagOptions = useMemo(() => mapLookupsToSelectOptions(item?.tags), [item]);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            setFormData(mapItemToFormData(item));
        }
    }, [isOpen, item, mode]);

    useEffect(() => {
        if (isOpen && mode === 'create') {
            setFormData(createDefaultItemFormData());
            setError(null);
        }
    }, [isOpen, mode]);

    useEffect(() => {
        if (!isOpen) {
            setFormData(createDefaultItemFormData());
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

        const validationError = validateItemFormData(formData);

        if (validationError) {
            setError(validationError);
            return;
        }

        setLoading(true);
        try {
            await onSave(buildItemModification(formData));
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
                        options: mapLookupsToSelectOptions(tags.data),
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
