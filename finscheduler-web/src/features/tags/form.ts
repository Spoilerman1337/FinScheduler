import type {TagDto, TagModification} from '../../api/types.ts';

export interface TagModalFormData {
    name: string;
    isActive: boolean;
}

export function createDefaultTagFormData(): TagModalFormData {
    return {
        name: '',
        isActive: true,
    };
}

export function mapTagToFormData(item?: TagDto | null): TagModalFormData {
    if (!item) {
        return createDefaultTagFormData();
    }

    return {
        name: item.name || '',
        isActive: typeof item.isActive === 'boolean' ? item.isActive : true,
    };
}

export function validateTagFormData(formData: TagModalFormData): string | null {
    if (!formData.name.trim()) {
        return 'Название обязательно для заполнения';
    }

    return null;
}

export function buildTagModification(formData: TagModalFormData): TagModification {
    return {
        name: formData.name.trim(),
        isActive: formData.isActive,
    };
}
