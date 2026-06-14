import type {ItemDetailedDto, ItemModification} from '../../api/items.types.ts';
import {categoryTranslations} from '../../models/items.ts';
import {type FieldValidator, runFieldValidators} from '../shared.ts';

export interface ItemFormData {
    name: string;
    description: string;
    price: string;
    cashback: string;
    isActive: boolean;
    category: string;
    tagIds: string[];
}

export const itemNameValidators: FieldValidator<string>[] = [
    {
        validate: (value) => value.trim().length > 0,
        errorMessage: 'Название обязательно для заполнения',
    },
];

export const itemPriceValidators: FieldValidator<string>[] = [
    {
        validate: (value) => value.trim() === '' || !Number.isNaN(parseFloat(value)),
        errorMessage: 'Цена должна быть числом',
    },
];

export const itemCashbackValidators: FieldValidator<string>[] = [
    {
        validate: (value) => {
            if (!value.trim()) {
                return true;
            }

            const parsedValue = parseFloat(value);

            return !Number.isNaN(parsedValue) && parsedValue >= 0 && parsedValue <= 100;
        },
        errorMessage: 'Кэшбэк должен быть числом от 0 до 100',
    },
];

export const itemCategoryValidators: FieldValidator<string>[] = [
    {
        validate: (value) => value.trim().length > 0 && Boolean(categoryTranslations[value]),
        errorMessage: 'Выберите категорию',
    },
];

export const itemFormValidators = {
    name: (value: string) => runFieldValidators(value, itemNameValidators),
    price: (value: string) => runFieldValidators(value, itemPriceValidators),
    cashback: (value: string) => runFieldValidators(value, itemCashbackValidators),
    category: (value: string) => runFieldValidators(value, itemCategoryValidators),
} as const;

export function createDefaultItemFormData(): ItemFormData {
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

export function mapItemToFormData(item?: ItemDetailedDto | null): ItemFormData {
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

export function normalizeItemFormData(formData: ItemFormData): ItemFormData {
    const payload = buildItemModification(formData);

    return {
        name: payload.name,
        description: payload.description ?? '',
        price: formData.price.trim() ? payload.price.toString() : '',
        cashback: formData.cashback.trim() ? payload.cashback.toString() : '',
        isActive: payload.isActive,
        category: payload.category,
        tagIds: formData.tagIds,
    };
}

export function buildItemModification(formData: ItemFormData): ItemModification {
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
