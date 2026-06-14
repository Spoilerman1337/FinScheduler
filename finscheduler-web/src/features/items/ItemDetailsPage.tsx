import {Card, Flex, SimpleGrid, Spinner} from '@chakra-ui/react';
import {useEffect, useState} from 'react';
import {Controller, useForm, useWatch} from 'react-hook-form';
import {useNavigate, useParams} from 'react-router-dom';
import type {ItemDetailedDto} from '../../api/items.types.ts';
import ItemsService from '../../api/items.ts';
import TagsService from '../../api/tags.ts';
import AsyncSelectField from '../../components/formFields/AsyncSelectField.tsx';
import NumberField from '../../components/formFields/NumberField.tsx';
import SelectField from '../../components/formFields/SelectField.tsx';
import SwitchField from '../../components/formFields/SwitchField.tsx';
import TextAreaField from '../../components/formFields/TextAreaField.tsx';
import TextField from '../../components/formFields/TextField.tsx';
import UnsavedChangesDialog from '../../components/unsavedChanges/UnsavedChangesDialog.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {useUnsavedChangesGuard} from '../../hooks/useUnsavedChangesGuard.ts';
import DetailsPageLayout, {type DetailsPageStatus} from '../../layout/details/DetailsPageLayout.tsx';
import {categoryOptions} from '../../models/items.ts';
import {buildEditItemPath, itemsListPath} from '../routes.ts';
import {mapLookupsToSelectOptions} from '../shared.ts';
import {
    buildItemModification,
    createDefaultItemFormData,
    itemFormValidators,
    mapItemToFormData,
    normalizeItemFormData,
    type ItemFormData,
} from './form.ts';

interface ItemDetailsPageProps {
    mode: 'create' | 'edit';
}

const TAGS_PAGE_SIZE = 20;
const itemsService = new ItemsService();
const tagsService = new TagsService();
const validationStatus: DetailsPageStatus = 'Ошибка валидации';

export default function ItemDetailsPage({mode}: ItemDetailsPageProps) {
    const navigate = useNavigate();
    const {itemId} = useParams<{itemId: string}>();
    const [item, setItem] = useState<ItemDetailedDto | null>(null);
    const [loading, setLoading] = useState(mode === 'edit');
    const [saving, setSaving] = useState(false);
    const [status, setStatus] = useState<DetailsPageStatus | null>(null);
    const {
        control,
        formState: {isDirty},
        handleSubmit,
        reset,
    } = useForm<ItemFormData>({
        defaultValues: createDefaultItemFormData(),
        mode: 'onSubmit',
        reValidateMode: 'onChange',
    });
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
        let isMounted = true;

        async function loadItem() {
            if (mode !== 'edit' || !itemId) {
                setItem(null);
                reset(createDefaultItemFormData());
                setLoading(false);
                setStatus(null);
                return;
            }

            setLoading(true);
            setStatus(null);

            try {
                const loadedItem = await itemsService.getDetailedInfo(itemId);

                if (!isMounted) {
                    return;
                }

                if (!loadedItem) {
                    setItem(null);
                    setStatus('Предмет не найден');
                    return;
                }

                setItem(loadedItem);
                reset(mapItemToFormData(loadedItem));
            } catch (err) {
                if (!isMounted) {
                    return;
                }

                setItem(null);
                setStatus(err instanceof Error ? err.message : 'Ошибка загрузки предмета');
            } finally {
                if (isMounted) {
                    setLoading(false);
                }
            }
        }

        void loadItem();

        return () => {
            isMounted = false;
        };
    }, [itemId, mode, reset]);

    const handleCancel = () => {
        navigate(itemsListPath);
    };

    const persistItem = async (formData: ItemFormData, closeAfterSave: boolean) => {
        setSaving(true);

        try {
            const normalizedFormData = normalizeItemFormData(formData);
            const payload = buildItemModification(normalizedFormData);

            if (mode === 'create') {
                const createdItemId = await itemsService.createItem(payload);

                reset(normalizedFormData);
                toaster.create({
                    title: 'Успешно',
                    description: 'Предмет успешно добавлен',
                    type: 'success',
                });

                if (closeAfterSave) {
                    scheduleNavigation(itemsListPath);
                    return;
                }

                scheduleNavigation(buildEditItemPath(createdItemId), {replace: true});
                return;
            }

            if (!itemId) {
                throw new Error('Не удалось определить предмет для сохранения');
            }

            await itemsService.updateItem(itemId, payload);
            reset(normalizedFormData);
            setItem((currentItem) =>
                currentItem
                    ? {
                          ...currentItem,
                          name: payload.name,
                          description: normalizedFormData.description,
                          price: payload.price,
                          cashback: payload.cashback,
                          isActive: payload.isActive,
                          category: payload.category,
                          tags: currentItem.tags,
                      }
                    : currentItem,
            );

            toaster.create({
                title: 'Успешно',
                description: 'Предмет успешно обновлён',
                type: 'success',
            });

            if (closeAfterSave) {
                scheduleNavigation(itemsListPath);
            }
        } catch (err) {
            setStatus(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setSaving(false);
        }
    };

    const handleSave = (closeAfterSave: boolean) => {
        setStatus(null);

        void handleSubmit(
            async (values) => {
                await persistItem(values, closeAfterSave);
            },
            () => {
                setStatus(validationStatus);
            },
        )();
    };

    const pageSubtitle =
        watchedName.trim() || item?.name || 'Заполните данные предмета для сохранения';
    const initialTagOptions = mapLookupsToSelectOptions(item?.tags);

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
                {label: 'Каталог', to: itemsListPath},
                {label: mode === 'create' ? 'Создание' : 'Редактирование'},
            ]}
            title={mode === 'create' ? 'Новый предмет' : 'Редактирование предмета'}
            subtitle={pageSubtitle}
            isActive={watchedIsActive}
            isDirty={isDirty}
            saving={saving}
            status={status}
            onSave={() => handleSave(false)}
            onSaveAndClose={() => handleSave(true)}
            onBack={handleCancel}
        >
            <SimpleGrid columns={{base: 1, xl: 2}} gap={6}>
                <Card.Root>
                    <Card.Header>
                        <Card.Title>1. Основная информация</Card.Title>
                        <Card.Description>Базовые данные предмета для каталога.</Card.Description>
                    </Card.Header>
                    <Card.Body gap={4}>
                        <Controller
                            name="name"
                            control={control}
                            rules={{validate: itemFormValidators.name}}
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
                            name="description"
                            control={control}
                            render={({field}) => (
                                <TextAreaField
                                    label="Описание"
                                    value={field.value}
                                    placeholder="Введите описание"
                                    rows={6}
                                    onChange={field.onChange}
                                />
                            )}
                        />
                    </Card.Body>
                </Card.Root>

                <Card.Root>
                    <Card.Header>
                        <Card.Title>2. Классификация</Card.Title>
                        <Card.Description>Категория, теги и актуальный статус.</Card.Description>
                    </Card.Header>
                    <Card.Body gap={4}>
                        <Controller
                            name="category"
                            control={control}
                            rules={{validate: itemFormValidators.category}}
                            render={({field, fieldState}) => (
                                <SelectField
                                    label="Категория"
                                    value={field.value}
                                    options={categoryOptions}
                                    placeholder="Выберите категорию"
                                    required
                                    invalid={fieldState.invalid}
                                    errorText={fieldState.error?.message}
                                    onChange={field.onChange}
                                />
                            )}
                        />
                        <Controller
                            name="tagIds"
                            control={control}
                            render={({field}) => (
                                <AsyncSelectField
                                    multiple
                                    label="Теги"
                                    value={field.value}
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

                <Card.Root gridColumn={{xl: '1 / -1'}}>
                    <Card.Header>
                        <Card.Title>3. Цена и бонусы</Card.Title>
                        <Card.Description>
                            Финансовые поля предмета из текущей модели.
                        </Card.Description>
                    </Card.Header>
                    <Card.Body>
                        <SimpleGrid columns={{base: 1, md: 2}} gap={4}>
                            <Controller
                                name="price"
                                control={control}
                                rules={{validate: itemFormValidators.price}}
                                render={({field, fieldState}) => (
                                    <NumberField
                                        label="Цена (₽)"
                                        value={field.value}
                                        defaultValue="0.00"
                                        step={0.01}
                                        min={0}
                                        invalid={fieldState.invalid}
                                        errorText={fieldState.error?.message}
                                        onChange={field.onChange}
                                    />
                                )}
                            />
                            <Controller
                                name="cashback"
                                control={control}
                                rules={{validate: itemFormValidators.cashback}}
                                render={({field, fieldState}) => (
                                    <NumberField
                                        label="Кэшбэк (%)"
                                        value={field.value}
                                        defaultValue="0"
                                        step={1}
                                        min={0}
                                        max={100}
                                        invalid={fieldState.invalid}
                                        errorText={fieldState.error?.message}
                                        onChange={field.onChange}
                                    />
                                )}
                            />
                        </SimpleGrid>
                    </Card.Body>
                </Card.Root>
            </SimpleGrid>

            <UnsavedChangesDialog open={isDialogOpen} onStay={stayOnPage} onLeave={leavePage} />
        </DetailsPageLayout>
    );
}
