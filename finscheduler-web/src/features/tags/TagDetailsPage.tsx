import {
    Button,
    Card,
    CloseButton,
    Dialog,
    Flex,
    Portal,
    Spinner,
    Stack,
    Text,
} from '@chakra-ui/react';
import {useEffect, useMemo, useState} from 'react';
import {Controller, useForm, useWatch} from 'react-hook-form';
import {useNavigate, useParams} from 'react-router-dom';
import SwitchField from '../../components/formFields/SwitchField.tsx';
import TextField from '../../components/formFields/TextField.tsx';
import UnsavedChangesDialog from '../../components/unsavedChanges/UnsavedChangesDialog.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {useUnsavedChangesGuard} from '../../hooks/useUnsavedChangesGuard.ts';
import DetailsPageLayout, {
    type DetailsPageStatus,
} from '../../layout/details/DetailsPageLayout.tsx';
import {buildEditTagPath, tagsListPath} from '../routes.ts';
import {useCreateTagMutation, useTagDetailsQuery, useUpdateTagMutation} from './queries.ts';
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    normalizeTagFormData,
    shouldConfirmTagDeactivation,
    tagFormValidators,
    type TagFormData,
} from './form.ts';

interface TagDetailsPageProps {
    mode: 'create' | 'edit';
}

interface PendingSaveAction {
    closeAfterSave: boolean;
    values: TagFormData;
}

const validationStatus: DetailsPageStatus = 'Ошибка валидации';

export default function TagDetailsPage({mode}: TagDetailsPageProps) {
    const navigate = useNavigate();
    const {tagId} = useParams<{tagId: string}>();
    const [status, setStatus] = useState<DetailsPageStatus | null>(null);
    const [pendingSaveAction, setPendingSaveAction] = useState<PendingSaveAction | null>(null);
    const tagQuery = useTagDetailsQuery(mode === 'edit' ? tagId : undefined);
    const createTagMutation = useCreateTagMutation();
    const updateTagMutation = useUpdateTagMutation();
    const {
        control,
        formState: {isDirty},
        handleSubmit,
        reset,
    } = useForm<TagFormData>({
        defaultValues: createDefaultTagFormData(),
        mode: 'onSubmit',
        reValidateMode: 'onChange',
    });
    const saving = createTagMutation.isPending || updateTagMutation.isPending;
    const loading = mode === 'edit' ? tagQuery.isPending : false;
    const tag = tagQuery.data ?? null;
    const {isDialogOpen, leavePage, scheduleNavigation, stayOnPage} = useUnsavedChangesGuard({
        isDirty,
        isDisabled: loading || saving,
    });

    const watchedName = useWatch({
        control,
        name: 'name',
        defaultValue: '',
    });
    const watchedIsActive = useWatch({
        control,
        name: 'isActive',
        defaultValue: true,
    });

    useEffect(() => {
        setStatus(null);

        if (mode !== 'edit' || !tagId) {
            reset(createDefaultTagFormData());
            return;
        }

        reset(createDefaultTagFormData());
    }, [mode, reset, tagId]);

    useEffect(() => {
        if (!tag) {
            return;
        }

        reset(mapTagToFormData(tag));
    }, [reset, tag]);

    const persistTag = async (formData: TagFormData, closeAfterSave: boolean) => {
        try {
            const normalizedFormData = normalizeTagFormData(formData);
            const payload = buildTagModification(normalizedFormData);

            if (mode === 'create') {
                const createdTagId = await createTagMutation.mutateAsync(payload);

                reset(normalizedFormData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Тег успешно добавлен',
                    type: 'success',
                });

                if (closeAfterSave) {
                    scheduleNavigation(tagsListPath);
                    return;
                }

                scheduleNavigation(buildEditTagPath(createdTagId), {replace: true});
                return;
            }

            if (!tagId) {
                throw new Error('Не удалось определить тег для сохранения');
            }

            await updateTagMutation.mutateAsync({tagId, tag: payload});
            reset(normalizedFormData);

            toaster.create({
                title: 'Успешно',
                description: 'Тег успешно обновлён',
                type: 'success',
            });

            if (closeAfterSave) {
                scheduleNavigation(tagsListPath);
            }
        } catch (err) {
            setStatus(err instanceof Error ? err.message : 'Ошибка при сохранении');
        }
    };

    const handleSave = (closeAfterSave: boolean) => {
        setStatus(null);

        void handleSubmit(
            async (values) => {
                if (shouldConfirmTagDeactivation(mode, tag, values)) {
                    setPendingSaveAction({closeAfterSave, values});
                    return;
                }

                await persistTag(values, closeAfterSave);
            },
            () => {
                setStatus(validationStatus);
            },
        )();
    };

    const handleConfirmDeactivation = async () => {
        if (!pendingSaveAction) {
            return;
        }

        const {closeAfterSave, values} = pendingSaveAction;
        setPendingSaveAction(null);
        await persistTag(values, closeAfterSave);
    };

    const pageSubtitle = watchedName.trim() || tag?.name || 'Заполните данные тега для сохранения';
    const loadStatus = useMemo(() => {
        if (mode !== 'edit') {
            return null;
        }

        if (tagQuery.isError) {
            return tagQuery.error instanceof Error
                ? tagQuery.error.message
                : 'Ошибка загрузки тега';
        }

        if (tagQuery.isSuccess && tag === null) {
            return 'Тег не найден';
        }

        return null;
    }, [mode, tag, tagQuery.error, tagQuery.isError, tagQuery.isSuccess]);
    const displayStatus = status ?? loadStatus;

    if (loading) {
        return (
            <Flex justify="center" align="center" minH="420px" width="100%">
                <Spinner size="xl" color="app.accent" />
            </Flex>
        );
    }

    return (
        <DetailsPageLayout
            breadcrumbItems={[
                {label: 'Теги', to: tagsListPath},
                {label: mode === 'create' ? 'Создание' : 'Редактирование'},
            ]}
            title={mode === 'create' ? 'Новый тег' : 'Редактирование тега'}
            subtitle={pageSubtitle}
            isActive={watchedIsActive}
            isDirty={isDirty}
            saving={saving}
            status={displayStatus}
            onSave={() => handleSave(false)}
            onSaveAndClose={() => handleSave(true)}
            onBack={() => navigate(tagsListPath)}
        >
            <Card.Root>
                <Card.Header>
                    <Card.Title>1. Основная информация</Card.Title>
                    <Card.Description>Базовые данные тега из текущей модели.</Card.Description>
                </Card.Header>
                <Card.Body gap={4}>
                    <Controller
                        name="name"
                        control={control}
                        rules={{validate: tagFormValidators.name}}
                        render={({field, fieldState}) => (
                            <TextField
                                label="Название"
                                value={field.value}
                                placeholder="Введите название"
                                required
                                invalid={fieldState.invalid}
                                errorText={fieldState.error?.message}
                                onChange={field.onChange}
                            />
                        )}
                    />
                    <Controller
                        name="isActive"
                        control={control}
                        render={({field}) => (
                            <SwitchField
                                label="Активен"
                                checked={field.value}
                                onChange={field.onChange}
                            />
                        )}
                    />
                </Card.Body>
            </Card.Root>

            <Dialog.Root
                open={pendingSaveAction !== null}
                onOpenChange={(details) => {
                    if (!details.open) {
                        setPendingSaveAction(null);
                    }
                }}
                placement="center"
                lazyMount
                unmountOnExit
            >
                <Portal>
                    <Dialog.Backdrop />
                    <Dialog.Positioner>
                        <Dialog.Content
                            bg="bg.layer1"
                            border="1px solid"
                            borderColor="border.error"
                            boxShadow="0 0 24px rgba(255, 74, 122, 0.18)"
                            maxW="560px"
                        >
                            <Dialog.Header>
                                <Dialog.Title color="fg">Подтвердить деактивацию тега</Dialog.Title>
                                <Dialog.CloseTrigger asChild>
                                    <CloseButton
                                        color="fg.error"
                                        _hover={{bg: 'rgba(255, 74, 122, 0.12)'}}
                                    />
                                </Dialog.CloseTrigger>
                            </Dialog.Header>
                            <Dialog.Body>
                                <Stack gap={3}>
                                    <Text color="fg">
                                        При деактивации тег будет отвязан от всех элементов
                                        каталога.
                                    </Text>
                                    <Text color="fg.muted">Подтвердите сохранение изменений.</Text>
                                </Stack>
                            </Dialog.Body>
                            <Dialog.Footer>
                                <Button variant="ghost" onClick={() => setPendingSaveAction(null)}>
                                    Отменить
                                </Button>
                                <Button
                                    bg="neon.pink"
                                    color="white"
                                    _hover={{bg: 'neon.pink', opacity: 0.88}}
                                    onClick={() => void handleConfirmDeactivation()}
                                >
                                    Подтвердить
                                </Button>
                            </Dialog.Footer>
                        </Dialog.Content>
                    </Dialog.Positioner>
                </Portal>
            </Dialog.Root>

            <UnsavedChangesDialog open={isDialogOpen} onStay={stayOnPage} onLeave={leavePage} />
        </DetailsPageLayout>
    );
}
