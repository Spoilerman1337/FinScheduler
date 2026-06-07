import {Badge, Box, Button, Card, Flex, SimpleGrid, Spinner, Stack, Text} from '@chakra-ui/react';
import {CheckCircle2, Save, X} from 'lucide-react';
import {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import type {ItemDto} from '../../api/items.types.ts';
import ItemsService from '../../api/items.ts';
import TagsService from '../../api/tags.ts';
import AsyncSelectField from '../../components/formFields/AsyncSelectField.tsx';
import NumberField from '../../components/formFields/NumberField.tsx';
import SelectField from '../../components/formFields/SelectField.tsx';
import SwitchField from '../../components/formFields/SwitchField.tsx';
import TextAreaField from '../../components/formFields/TextAreaField.tsx';
import TextField from '../../components/formFields/TextField.tsx';
import Breadcrumbs from '../../components/ui/Breadcrumbs.tsx';
import {toaster} from '../../components/ui/toaster-instance.ts';
import {categoryOptions} from '../../models/items.ts';
import {mapLookupsToSelectOptions} from '../shared.ts';
import {
    buildItemModification,
    createDefaultItemFormData,
    mapItemToFormData,
    type ItemFormData,
    validateItemFormData,
} from './form.ts';
import {buildEditItemPath, itemsListPath} from '../routes.ts';

interface ItemDetailsPageProps {
    mode: 'create' | 'edit';
}

const TAGS_PAGE_SIZE = 20;
const itemsService = new ItemsService();
const tagsService = new TagsService();

export default function ItemDetailsPage({mode}: ItemDetailsPageProps) {
    const navigate = useNavigate();
    const {itemId} = useParams<{itemId: string}>();
    const [item, setItem] = useState<ItemDto | null>(null);
    const [formData, setFormData] = useState<ItemFormData>(createDefaultItemFormData);
    const [loading, setLoading] = useState(mode === 'edit');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let isMounted = true;

        async function loadItem() {
            if (mode !== 'edit' || !itemId) {
                setItem(null);
                setFormData(createDefaultItemFormData());
                setLoading(false);
                setError(null);
                return;
            }

            setLoading(true);
            setError(null);

            try {
                const loadedItem = await itemsService.getItem(itemId);

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

                toaster.create({
                    title: 'Успешно',
                    description: 'Предмет успешно добавлен',
                    type: 'success',
                });

                if (closeAfterSave) {
                    navigate(itemsListPath);
                    return;
                }

                navigate(buildEditItemPath(createdItemId), {replace: true});
                return;
            }

            if (!itemId) {
                throw new Error('Не удалось определить предмет для сохранения');
            }

            await itemsService.updateItem(itemId, payload);
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
                navigate(itemsListPath);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
        } finally {
            setSaving(false);
        }
    };

    const pageTitle = mode === 'create' ? 'Новый предмет' : 'Редактирование предмета';
    const pageSubtitle =
        formData.name.trim() || item?.name || 'Заполните данные предмета для сохранения';
    const initialTagOptions = mapLookupsToSelectOptions(item?.tags);

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
                    {label: 'Каталог', to: itemsListPath},
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
                        <Card.Description>
                            Финансовые поля предмета из текущей модели.
                        </Card.Description>
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
        </Stack>
    );
}
