import {useEffect, useState} from 'react';
import type {TagDto, TagModification} from '../../../api/types.ts';
import SwitchField from '../../../components/formFields/SwitchField.tsx';
import TextField from '../../../components/formFields/TextField.tsx';
import FormModal from '../../../components/formModal/FormModal.tsx';
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    type TagModalFormData,
    validateTagFormData,
} from '../form.ts';

interface TagModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (tag: TagModification) => Promise<void>;
    item?: TagDto | null;
    mode: 'create' | 'edit';
}

export default function TagModal({isOpen, onClose, onSave, item, mode}: TagModalProps) {
    const [formData, setFormData] = useState<TagModalFormData>(createDefaultTagFormData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            setFormData(mapTagToFormData(item));
        }
    }, [isOpen, item, mode]);

    useEffect(() => {
        if (isOpen && mode === 'create') {
            setFormData(createDefaultTagFormData());
            setError(null);
        }
    }, [isOpen, mode]);

    useEffect(() => {
        if (!isOpen) {
            setFormData(createDefaultTagFormData());
            setError(null);
        }
    }, [isOpen]);

    const handleSubmit = async () => {
        setError(null);

        const validationError = validateTagFormData(formData);

        if (validationError) {
            setError(validationError);
            return;
        }

        setLoading(true);
        try {
            await onSave(buildTagModification(formData));
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
                onChange={(value) => setFormData((prev) => ({...prev, name: value}))}
            />

            <SwitchField
                label="Активен"
                checked={formData.isActive}
                onChange={(value) => setFormData((prev) => ({...prev, isActive: value}))}
            />
        </FormModal>
    );
}
