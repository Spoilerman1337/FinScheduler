import type {TagDetailedDto, TagModification} from '../../api/tags.types.ts';
import {type FieldValidator, runFieldValidators} from '../shared.ts';

export interface TagFormData {
    name: string;
    isActive: boolean;
}

export const tagNameValidators: FieldValidator<string>[] = [
    {
        validate: (value) => value.trim().length > 0,
        errorMessage: 'Название обязательно для заполнения',
    },
];

export const tagFormValidators = {
    name: (value: string) => runFieldValidators(value, tagNameValidators),
} as const;

export function createDefaultTagFormData(): TagFormData {
    return {
        name: '',
        isActive: true,
    };
}

export function mapTagToFormData(item?: TagDetailedDto | null): TagFormData {
    if (!item) {
        return createDefaultTagFormData();
    }

    return {
        name: item.name || '',
        isActive: typeof item.isActive === 'boolean' ? item.isActive : true,
    };
}

export function normalizeTagFormData(formData: TagFormData): TagFormData {
    return {
        name: formData.name.trim(),
        isActive: formData.isActive,
    };
}

export function buildTagModification(formData: TagFormData): TagModification {
    return {
        name: formData.name.trim(),
        isActive: formData.isActive,
    };
}

export function shouldConfirmTagDeactivation(
    mode: 'create' | 'edit',
    tag: TagDetailedDto | null,
    formData: TagFormData,
): boolean {
    return mode === 'edit' && tag?.isActive === true && formData.isActive === false;
}
