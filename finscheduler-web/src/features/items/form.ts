import type {ItemDto, ItemModification} from '../../api/types.ts';
import {categoryTranslations} from '../../models/items.ts';

export interface ItemModalFormData {
    name: string;
    description: string;
    price: string;
    cashback: string;
    isActive: boolean;
    category: string;
    tagIds: string[];
}

export function createDefaultItemFormData(): ItemModalFormData {
    return {
        name: '',
        description: '',
        price: '',
        cashback: '',
        isActive: true,
        category: '',
        tagIds: [],
    };
}

export function mapItemToFormData(item?: ItemDto | null): ItemModalFormData {
    if (!item) {
        return createDefaultItemFormData();
    }

    return {
        name: item.name || '',
        description: item.description || '',
        price: item.price !== undefined && item.price !== null ? item.price.toString() : '',
        cashback:
            item.cashback !== undefined && item.cashback !== null ? item.cashback.toString() : '',
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
    };
}

export function validateItemFormData(formData: ItemModalFormData): string | null {
    if (!formData.name.trim()) {
        return 'Название обязательно для заполнения';
    }

    if (formData.price && Number.isNaN(parseFloat(formData.price))) {
        return 'Цена должна быть числом';
    }

    if (
        formData.cashback &&
        (Number.isNaN(parseFloat(formData.cashback)) ||
            parseFloat(formData.cashback) < 0 ||
            parseFloat(formData.cashback) > 100)
    ) {
        return 'Кэшбэк должен быть числом от 0 до 100';
    }

    if (!formData.category.trim() || !categoryTranslations[formData.category]) {
        return 'Выберите категорию';
    }

    return null;
}

export function buildItemModification(formData: ItemModalFormData): ItemModification {
    return {
        name: formData.name.trim(),
        description: formData.description.trim() || undefined,
        price: formData.price ? parseFloat(formData.price) : 0,
        cashback: formData.cashback ? parseFloat(formData.cashback) : 0,
        isActive: formData.isActive,
        category: formData.category,
        tagIds: formData.tagIds,
    };
}
