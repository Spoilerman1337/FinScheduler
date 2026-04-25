import {useEffect, useState} from "react";
import type {TagDto} from "../../../api/types.ts";
import SwitchField from "../../../components/formFields/SwitchField.tsx";
import TextField from "../../../components/formFields/TextField.tsx";
import FormModal from "../../../components/formModal/FormModal.tsx";

interface TagModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (tag: Omit<TagDto, 'id'>) => Promise<void>;
    item?: TagDto | null;
    mode: 'create' | 'edit';
}

export default function TagModal({isOpen, onClose, onSave, item, mode}: TagModalProps) {
    const getDefaultFormData = () => ({
        name: '',
        isActive: true,
    });

    const [formData, setFormData] = useState(getDefaultFormData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (isOpen && mode === 'edit' && item) {
            setFormData({
                name: item.name || '',
                isActive: typeof item.isActive === 'boolean' ? item.isActive : true,
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

    const handleSubmit = async () => {
        setError(null);

        if (!formData.name.trim()) {
            setError('Название обязательно для заполнения');
            return;
        }

        setLoading(true);
        try {
            await onSave({
                name: formData.name.trim(),
                isActive: formData.isActive,
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
