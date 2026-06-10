import {
    Badge,
    Box,
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
import {CheckCircle2, Save, X} from 'lucide-react';
import {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import type {TagDto} from '../../api/tags.types.ts';
import TagsService from '../../api/tags.ts';
import SwitchField from '../../components/formFields/SwitchField.tsx';
import TextField from '../../components/formFields/TextField.tsx';
import UnsavedChangesDialog from '../../components/unsavedChanges/UnsavedChangesDialog.tsx';
import {useUnsavedChangesGuard} from '../../hooks/useUnsavedChangesGuard.ts';
import Breadcrumbs from '../../components/ui/Breadcrumbs.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {buildEditTagPath, tagsListPath} from '../routes.ts';
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    shouldConfirmTagDeactivation,
    type TagFormData,
    validateTagFormData,
} from './form.ts';

interface TagDetailsPageProps {
    mode: 'create' | 'edit';
}

interface PendingSaveAction {
    closeAfterSave: boolean;
}

const tagsService = new TagsService();

export default function TagDetailsPage({mode}: TagDetailsPageProps) {
    const navigate = useNavigate();
    const {tagId} = useParams<{tagId: string}>();
    const [tag, setTag] = useState<TagDto | null>(null);
    const [formData, setFormData] = useState<TagFormData>(createDefaultTagFormData);
    const [loading, setLoading] = useState(mode === 'edit');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isDirty, setIsDirty] = useState(false);
    const [pendingSaveAction, setPendingSaveAction] = useState<PendingSaveAction | null>(null);
    const {isDialogOpen, leavePage, scheduleNavigation, stayOnPage} = useUnsavedChangesGuard({
        isDirty,
        isDisabled: loading || saving,
    });

    useEffect(() => {
        let isMounted = true;

        async function loadTag() {
            if (mode !== 'edit' || !tagId) {
                setTag(null);
                setFormData(createDefaultTagFormData());
                setLoading(false);
                setError(null);
                setIsDirty(false);
                return;
            }

            setLoading(true);
            setError(null);

            try {
                const loadedTag = await tagsService.getTag(tagId);

                if (!isMounted) {
                    return;
                }

                if (!loadedTag) {
                    setTag(null);
                    setError('Тег не найден');
                    return;
                }

                setTag(loadedTag);
                setFormData(mapTagToFormData(loadedTag));
                setIsDirty(false);
            } catch (err) {
                if (!isMounted) {
                    return;
                }

                setTag(null);
                setError(err instanceof Error ? err.message : 'Ошибка загрузки тега');
            } finally {
                if (isMounted) {
                    setLoading(false);
                }
            }
        }

        void loadTag();

        return () => {
            isMounted = false;
        };
    }, [mode, tagId]);

    const updateFormData = <K extends keyof TagFormData>(field: K, value: TagFormData[K]) => {
        setIsDirty(true);
        setFormData((prev) => ({...prev, [field]: value}));
    };

    const handleCancel = () => {
        navigate(tagsListPath);
    };

    const persistTag = async (closeAfterSave: boolean) => {
        setSaving(true);

        try {
            const payload = buildTagModification(formData);

            if (mode === 'create') {
                const createdTagId = await tagsService.createTag(payload);

                setIsDirty(false);
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

            await tagsService.updateTag(tagId, payload);
            setIsDirty(false);
            setTag((currentTag) =>
                currentTag
                    ? {
                          ...currentTag,
                          name: payload.name,
                          isActive: payload.isActive,
                      }
                    : currentTag,
            );

            toaster.create({
                title: 'Успешно',
                description: 'Тег успешно обновлён',
                type: 'success',
            });

            if (closeAfterSave) {
                scheduleNavigation(tagsListPath);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setSaving(false);
        }
    };

    const handleSave = async (closeAfterSave: boolean) => {
        setError(null);

        const validationError = validateTagFormData(formData);

        if (validationError) {
            setError(validationError);
            return;
        }

        if (shouldConfirmTagDeactivation(mode, tag, formData)) {
            setPendingSaveAction({closeAfterSave});
            return;
        }

        await persistTag(closeAfterSave);
    };

    const handleConfirmDeactivation = async () => {
        if (!pendingSaveAction) {
            return;
        }

        const {closeAfterSave} = pendingSaveAction;
        setPendingSaveAction(null);
        await persistTag(closeAfterSave);
    };

    const pageTitle = mode === 'create' ? 'Новый тег' : 'Редактирование тега';
    const pageSubtitle =
        formData.name.trim() || tag?.name || 'Заполните данные тега для сохранения';

    if (loading) {
        return (
            <Flex justify="center" align="center" minH="420px" width="100%">
                <Spinner size="xl" color="app.accent" />
            </Flex>
        );
    }

    return (
        <Stack width="100%" gap={6} pb={6}>
            <Breadcrumbs
                items={[
                    {label: 'Теги', to: tagsListPath},
                    {label: mode === 'create' ? 'Создание' : 'Редактирование'},
                ]}
            />

            <Card.Root overflow="visible">
                <Box
                    position="absolute"
                    inset="0"
                    pointerEvents="none"
                    borderRadius="inherit"
                    background="linear-gradient(90deg, rgba(32, 208, 255, 0.08), rgba(143, 120, 255, 0.04))"
                />
                <Card.Body position="relative" gap={6} pb={{base: 7, xl: 6}}>
                    <Flex
                        direction={{base: 'column', xl: 'row'}}
                        align="flex-start"
                        justify="space-between"
                        gap={6}
                    >
                        <Stack gap={2} flex="1" minW={0}>
                            <Text
                                textStyle="3xl"
                                color="fg"
                                fontWeight="700"
                                letterSpacing="tight"
                                lineHeight="tight"
                            >
                                {pageTitle}
                            </Text>
                            <Text color="fg.muted" textStyle="lg" lineHeight="snug">
                                {pageSubtitle}
                            </Text>
                            <Badge
                                alignSelf="flex-start"
                                px={3}
                                py={1}
                                borderRadius="full"
                                bg={formData.isActive ? 'neon.green' : 'neon.pink'}
                                color="bg.base"
                            >
                                {formData.isActive ? 'Активен' : 'Неактивен'}
                            </Badge>
                        </Stack>

                        <Flex wrap="wrap" gap={3}>
                            <Button onClick={() => void handleSave(false)} loading={saving}>
                                <Save />
                                Сохранить
                            </Button>
                            <Button
                                variant="surface"
                                borderColor="app.cardBorderActive"
                                boxShadow="app.glowViolet"
                                _hover={{
                                    bg: 'rgba(143, 120, 255, 0.18)',
                                    borderColor: 'app.cardBorderActive',
                                }}
                                onClick={() => void handleSave(true)}
                                loading={saving}
                            >
                                <CheckCircle2 />
                                Сохранить и закрыть
                            </Button>
                            <Button variant="outline" onClick={handleCancel} disabled={saving}>
                                <X />
                                Отмена
                            </Button>
                        </Flex>
                    </Flex>
                </Card.Body>
            </Card.Root>

            {error ? (
                <Card.Root borderColor="border.error">
                    <Card.Body>
                        <Text color="fg.error">{error}</Text>
                    </Card.Body>
                </Card.Root>
            ) : null}

            <Card.Root>
                <Card.Header>
                    <Card.Title>1. Основная информация</Card.Title>
                    <Card.Description>Базовые данные тега из текущей модели.</Card.Description>
                </Card.Header>
                <Card.Body gap={4}>
                    <TextField
                        label="Название"
                        value={formData.name}
                        placeholder="Введите название"
                        required
                        onChange={(value) => updateFormData('name', value)}
                    />
                    <SwitchField
                        label="Активен"
                        checked={formData.isActive}
                        onChange={(value) => updateFormData('isActive', value)}
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
        </Stack>
    );
}
