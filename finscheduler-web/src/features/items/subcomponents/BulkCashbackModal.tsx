import {Box, Flex, Stack, Text} from '@chakra-ui/react';
import {useEffect, useState} from 'react';
import AsyncSelectField from '../../../components/formFields/AsyncSelectField.tsx';
import NumberField from '../../../components/formFields/NumberField.tsx';
import FormModal from '../../../components/formModal/FormModal.tsx';
import {useTagLookupOptionsLoader} from '../../tags/queries.ts';

const TAGS_PAGE_SIZE = 20;

export interface SelectedItemSummary {
    id: string;
    label: string;
}

export interface BulkCashbackSubmitPayload {
    cashback: number;
    tagId?: string;
    itemIds?: string[];
}

interface BulkCashbackModalProps {
    isOpen: boolean;
    selectedItems: SelectedItemSummary[];
    loading?: boolean;
    error?: string | null;
    onClose: () => void;
    onSubmit: (payload: BulkCashbackSubmitPayload) => Promise<void>;
}

function validateCashback(value: string): string | null {
    if (!value.trim()) {
        return 'Укажите кешбек от 0 до 100';
    }

    const parsedValue = Number(value);
    if (!Number.isInteger(parsedValue) || parsedValue < 0 || parsedValue > 100) {
        return 'Кешбек должен быть целым числом от 0 до 100';
    }

    return null;
}

export default function BulkCashbackModal(props: BulkCashbackModalProps) {
    const {isOpen, selectedItems, loading = false, error = null, onClose, onSubmit} = props;
    const [tagId, setTagId] = useState('');
    const [cashback, setCashback] = useState('0');
    const [validationError, setValidationError] = useState<string | null>(null);
    const hasSelection = selectedItems.length > 0;
    const loadTagOptions = useTagLookupOptionsLoader(TAGS_PAGE_SIZE);

    useEffect(() => {
        if (!isOpen) {
            return;
        }

        setTagId('');
        setCashback('0');
        setValidationError(null);
    }, [hasSelection, isOpen]);

    const handleSubmit = async () => {
        const cashbackError = validateCashback(cashback);
        if (cashbackError) {
            setValidationError(cashbackError);
            return;
        }

        if (!hasSelection && !tagId) {
            setValidationError('Выберите тег для массового обновления');
            return;
        }

        setValidationError(null);
        await onSubmit({
            cashback: Number(cashback),
            tagId: hasSelection ? undefined : tagId,
            itemIds: hasSelection ? selectedItems.map((item) => item.id) : undefined,
        });
    };

    return (
        <FormModal
            isOpen={isOpen}
            onClose={onClose}
            onSubmit={() => void handleSubmit()}
            title="Массовое обновление кешбека"
            error={validationError ?? error}
            loading={loading}
        >
            {hasSelection ? (
                <Stack gap={4}>
                    <Text color="fg.muted" fontSize="sm">
                        Кешбек будет обновлен только у выбранных элементов.
                    </Text>
                    <Box
                        maxH="240px"
                        overflowY="auto"
                        border="1px solid"
                        borderColor="glass.border"
                        borderRadius="md"
                        bg="bg.layer2"
                        px={3}
                        py={2}
                    >
                        <Stack gap={2}>
                            {selectedItems.map((item) => (
                                <Flex key={item.id} align="center" justify="space-between" gap={3}>
                                    <Text color="neon.blue" fontWeight="medium">
                                        {item.label}
                                    </Text>
                                    <Text color="fg.muted" fontSize="xs">
                                        {item.id}
                                    </Text>
                                </Flex>
                            ))}
                        </Stack>
                    </Box>
                </Stack>
            ) : (
                <Stack gap={4}>
                    <Text color="fg.muted" fontSize="sm">
                        Выберите один тег, и кешбек обновится у всех элементов этого тега.
                    </Text>
                    <AsyncSelectField
                        label="Тег"
                        value={tagId}
                        placeholder="Выберите тег"
                        emptyText="Теги не найдены"
                        loadOptions={loadTagOptions}
                        onChange={setTagId}
                    />
                </Stack>
            )}

            <NumberField
                label="Кешбек (%)"
                value={cashback}
                defaultValue="0"
                step={1}
                min={0}
                max={100}
                onChange={setCashback}
            />
        </FormModal>
    );
}
