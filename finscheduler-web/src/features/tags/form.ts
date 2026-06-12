import type {TagDetailedDto, TagModification} from '../../api/tags.types.ts';

export interface TagFormData {
    name: string;
    isActive: boolean;
}

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

export function validateTagFormData(formData: TagFormData): string | null {
    if (!formData.name.trim()) {
        return 'Название обязательно для заполнения';
    }

    return null;
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
