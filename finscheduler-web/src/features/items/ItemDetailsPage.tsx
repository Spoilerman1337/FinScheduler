import {Card, Flex, SimpleGrid, Spinner} from '@chakra-ui/react';
import {useEffect, useState} from 'react';
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
import DetailsPageLayout from '../../layout/details/DetailsPageLayout.tsx';
import {categoryOptions} from '../../models/items.ts';
import {buildEditItemPath, itemsListPath} from '../routes.ts';
import {mapLookupsToSelectOptions} from '../shared.ts';
import {
    buildItemModification,
    createDefaultItemFormData,
    mapItemToFormData,
    type ItemFormData,
    validateItemFormData,
} from './form.ts';

interface ItemDetailsPageProps {
    mode: 'create' | 'edit';
}

const TAGS_PAGE_SIZE = 20;
const itemsService = new ItemsService();
const tagsService = new TagsService();

export default function ItemDetailsPage({mode}: ItemDetailsPageProps) {
    const navigate = useNavigate();
    const {itemId} = useParams<{itemId: string}>();
    const [item, setItem] = useState<ItemDetailedDto | null>(null);
    const [formData, setFormData] = useState<ItemFormData>(createDefaultItemFormData);
    const [loading, setLoading] = useState(mode === 'edit');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isDirty, setIsDirty] = useState(false);
    const {isDialogOpen, leavePage, scheduleNavigation, stayOnPage} = useUnsavedChangesGuard({
        isDirty,
        isDisabled: loading || saving,
    });

    useEffect(() => {
        let isMounted = true;

        async function loadItem() {
            if (mode !== 'edit' || !itemId) {
                setItem(null);
                setFormData(createDefaultItemFormData());
                setLoading(false);
                setError(null);
                setIsDirty(false);
                return;
            }

            setLoading(true);
            setError(null);

            try {
                const loadedItem = await itemsService.getDetailedInfo(itemId);

                if (!isMounted) {
                    return;
                }

                if (!loadedItem) {
                    setItem(null);
                    setError('Предмет не найден');
                    return;
                }

                setItem(loadedItem);
                setFormData(mapItemToFormData(loadedItem));
                setIsDirty(false);
            } catch (err) {
                if (!isMounted) {
                    return;
                }

                setItem(null);
                setError(err instanceof Error ? err.message : 'Ошибка загрузки предмета');
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
    }, [itemId, mode]);

    const updateFormData = <K extends keyof ItemFormData>(field: K, value: ItemFormData[K]) => {
        setIsDirty(true);
        setFormData((prev) => ({...prev, [field]: value}));
    };

    const handleCancel = () => {
        navigate(itemsListPath);
    };

    const handleSave = async (closeAfterSave: boolean) => {
        setError(null);

        const validationError = validateItemFormData(formData);

        if (validationError) {
            setError(validationError);
            return;
        }

        setSaving(true);

        try {
            const payload = buildItemModification(formData);

            if (mode === 'create') {
                const createdItemId = await itemsService.createItem(payload);

                setIsDirty(false);
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
            setIsDirty(false);
            setItem((currentItem) =>
                currentItem
                    ? {
                          ...currentItem,
                          name: payload.name,
                          description: payload.description,
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
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setSaving(false);
        }
    };

    const pageSubtitle = formData.name.trim() || item?.name || 'Заполните данные предмета для сохранения';
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
            isActive={formData.isActive}
            isDirty={isDirty}
            saving={saving}
            error={error}
            onSave={() => void handleSave(false)}
            onSaveAndClose={() => void handleSave(true)}
            onBack={handleCancel}
        >
            <SimpleGrid columns={{base: 1, xl: 2}} gap={6}>
                <Card.Root>
                    <Card.Header>
                        <Card.Title>1. Основная информация</Card.Title>
                        <Card.Description>Базовые данные предмета для каталога.</Card.Description>
                    </Card.Header>
                    <Card.Body gap={4}>
                        <TextField
                            label="Название"
                            value={formData.name}
                            placeholder="Введите название"
                            required
                            onChange={(value) => updateFormData('name', value)}
                        />
                        <TextAreaField
                            label="Описание"
                            value={formData.description}
                            placeholder="Введите описание"
                            rows={6}
                            onChange={(value) => updateFormData('description', value)}
                        />
                    </Card.Body>
                </Card.Root>

                <Card.Root>
                    <Card.Header>
                        <Card.Title>2. Классификация</Card.Title>
                        <Card.Description>Категория, теги и актуальный статус.</Card.Description>
                    </Card.Header>
                    <Card.Body gap={4}>
                        <SelectField
                            label="Категория"
                            value={formData.category}
                            options={categoryOptions}
                            placeholder="Выберите категорию"
                            required
                            onChange={(value) => updateFormData('category', value)}
                        />
                        <AsyncSelectField
                            multiple
                            label="Теги"
                            value={formData.tagIds}
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
                            onChange={(value) => updateFormData('tagIds', value)}
                        />
                        <SwitchField
                            label="Активен"
                            checked={formData.isActive}
                            onChange={(value) => updateFormData('isActive', value)}
                        />
                    </Card.Body>
                </Card.Root>

                <Card.Root gridColumn={{xl: '1 / -1'}}>
                    <Card.Header>
                        <Card.Title>3. Цена и бонусы</Card.Title>
                        <Card.Description>Финансовые поля предмета из текущей модели.</Card.Description>
                    </Card.Header>
                    <Card.Body>
                        <SimpleGrid columns={{base: 1, md: 2}} gap={4}>
                            <NumberField
                                label="Цена (₽)"
                                value={formData.price}
                                defaultValue="0.00"
                                step={0.01}
                                min={0}
                                onChange={(value) => updateFormData('price', value)}
                            />
                            <NumberField
                                label="Кэшбэк (%)"
                                value={formData.cashback}
                                defaultValue="0"
                                step={1}
                                min={0}
                                max={100}
                                onChange={(value) => updateFormData('cashback', value)}
                            />
                        </SimpleGrid>
                    </Card.Body>
                </Card.Root>
            </SimpleGrid>

            <UnsavedChangesDialog open={isDialogOpen} onStay={stayOnPage} onLeave={leavePage} />
        </DetailsPageLayout>
    );
}
