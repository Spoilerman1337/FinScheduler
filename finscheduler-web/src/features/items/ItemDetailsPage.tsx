import {Card, Flex, SimpleGrid, Spinner} from '@chakra-ui/react';
import {useEffect, useMemo, useState} from 'react';
import {Controller, useForm, useWatch} from 'react-hook-form';
import {useNavigate, useParams} from 'react-router-dom';
import AsyncSelectField from '../../components/formFields/AsyncSelectField.tsx';
import NumberField from '../../components/formFields/NumberField.tsx';
import SelectField from '../../components/formFields/SelectField.tsx';
import SwitchField from '../../components/formFields/SwitchField.tsx';
import TextAreaField from '../../components/formFields/TextAreaField.tsx';
import TextField from '../../components/formFields/TextField.tsx';
import UnsavedChangesDialog from '../../components/unsavedChanges/UnsavedChangesDialog.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {useUnsavedChangesGuard} from '../../hooks/useUnsavedChangesGuard.ts';
import DetailsPageLayout, {
    type DetailsPageStatus,
} from '../../layout/details/DetailsPageLayout.tsx';
import {categoryOptions} from '../../models/items.ts';
import {buildEditItemPath, itemsListPath} from '../routes.ts';
import {mapLookupsToSelectOptions} from '../shared.ts';
import {useTagLookupOptionsLoader} from '../tags/queries.ts';
import {useCreateItemMutation, useItemDetailsQuery, useUpdateItemMutation} from './queries.ts';
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
const validationStatus: DetailsPageStatus = 'Ошибка валидации';

export default function ItemDetailsPage({mode}: ItemDetailsPageProps) {
    const navigate = useNavigate();
    const {itemId} = useParams<{itemId: string}>();
    const [status, setStatus] = useState<DetailsPageStatus | null>(null);
    const itemQuery = useItemDetailsQuery(mode === 'edit' ? itemId : undefined);
    const createItemMutation = useCreateItemMutation();
    const updateItemMutation = useUpdateItemMutation();
    const loadTagOptions = useTagLookupOptionsLoader(TAGS_PAGE_SIZE);
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
    const saving = createItemMutation.isPending || updateItemMutation.isPending;
    const loading = mode === 'edit' ? itemQuery.isPending : false;
    const item = itemQuery.data ?? null;
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

        if (mode !== 'edit' || !itemId) {
            reset(createDefaultItemFormData());
            return;
        }

        reset(createDefaultItemFormData());
    }, [itemId, mode, reset]);

    useEffect(() => {
        if (!item) {
            return;
        }

        reset(mapItemToFormData(item));
    }, [item, reset]);

    const persistItem = async (formData: ItemFormData, closeAfterSave: boolean) => {
        try {
            const normalizedFormData = normalizeItemFormData(formData);
            const payload = buildItemModification(normalizedFormData);

            if (mode === 'create') {
                const createdItemId = await createItemMutation.mutateAsync(payload);

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

            await updateItemMutation.mutateAsync({itemId, item: payload});
            reset(normalizedFormData);

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
    const loadStatus = useMemo(() => {
        if (mode !== 'edit') {
            return null;
        }

        if (itemQuery.isError) {
            return itemQuery.error instanceof Error
                ? itemQuery.error.message
                : 'Ошибка загрузки предмета';
        }

        if (itemQuery.isSuccess && item === null) {
            return 'Предмет не найден';
        }

        return null;
    }, [item, itemQuery.error, itemQuery.isError, itemQuery.isSuccess, mode]);
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
                {label: 'Каталог', to: itemsListPath},
                {label: mode === 'create' ? 'Создание' : 'Редактирование'},
            ]}
            title={mode === 'create' ? 'Новый предмет' : 'Редактирование предмета'}
            subtitle={pageSubtitle}
            isActive={watchedIsActive}
            isDirty={isDirty}
            saving={saving}
            status={displayStatus}
            onSave={() => handleSave(false)}
            onSaveAndClose={() => handleSave(true)}
            onBack={() => navigate(itemsListPath)}
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
                                    placeholder="Выберите теги"
                                    emptyText="Теги не найдены"
                                    collapseThreshold={4}
                                    loadOptions={loadTagOptions}
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
                                        label="Кешбек (%)"
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
