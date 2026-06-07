import {Box, Card, Flex, IconButton, SimpleGrid, Text} from '@chakra-ui/react';
import {ArrowUpRight} from 'lucide-react';
import ListingSelectionCheckbox from '../../dataListing/selectionCheckbox/ListingSelectionCheckbox.tsx';
import type {DataListingColumn} from '../../dataListing/types.ts';
import DataShowcaseField from './DataShowcaseField.tsx';

interface DataShowcaseCardProps<T> {
    row: T;
    rowId: string;
    columns: DataListingColumn<T>[];
    selectable: boolean;
    isSelected: boolean;
    onSelectRow: (rowId: string, checked: boolean) => void;
    onEdit?: (row: T) => void;
    onRowEdit: (row: T) => void;
}

const accentMap = {
    cyan: {
        color: 'app.accent',
        surface: 'app.accentSoft',
        shadow: 'app.glowCyan',
    },
    violet: {
        color: 'app.accentViolet',
        surface: 'rgba(143, 120, 255, 0.16)',
        shadow: 'app.glowViolet',
    },
} as const;

export default function DataShowcaseCard<T extends object>(props: DataShowcaseCardProps<T>) {
    const {row, rowId, columns, selectable, isSelected, onSelectRow, onEdit, onRowEdit} = props;
    const [primaryColumn, ...secondaryColumns] = columns;
    const accent = isSelected ? accentMap.violet : accentMap.cyan;

    return (
        <Card.Root
            variant="outline"
            minH="100%"
            cursor={onEdit ? 'pointer' : 'default'}
            overflow="hidden"
            onDoubleClick={() => onEdit && onRowEdit(row)}
            _hover={{
                borderColor: accent.color,
                boxShadow: accent.shadow,
                transform: 'translateY(-2px)',
                bg: 'app.cardBgHover',
            }}
            borderColor={isSelected ? 'app.cardBorderActive' : undefined}
            bg={isSelected ? 'app.cardBgHover' : undefined}
            boxShadow={isSelected ? 'app.glowViolet' : undefined}
        >
            <Box
                position="absolute"
                top="-10"
                left="-10"
                boxSize="40"
                borderRadius="full"
                pointerEvents="none"
                bg={accent.surface}
                filter="blur(42px)"
                opacity={0.75}
            />

            <Card.Body gap="5" position="relative" pe={onEdit || selectable ? '16' : undefined}>
                <Flex flex="1" minW={0}>
                    {primaryColumn ? (
                        <DataShowcaseField row={row} column={primaryColumn} isPrimary />
                    ) : (
                        <Text color="fg.muted">—</Text>
                    )}
                </Flex>

                {secondaryColumns.length > 0 ? (
                    <SimpleGrid columns={{base: 1, sm: 2}} gap="4" mt="auto" pt="2">
                        {secondaryColumns.map((column) => (
                            <DataShowcaseField key={column.key} row={row} column={column} />
                        ))}
                    </SimpleGrid>
                ) : null}
            </Card.Body>

            {selectable ? (
                <Flex
                    position="absolute"
                    top="5"
                    right="5"
                    boxSize="9"
                    align="center"
                    justify="center"
                    zIndex={1}
                    onClick={(event) => event.stopPropagation()}
                >
                    <ListingSelectionCheckbox
                        checked={isSelected}
                        onCheckedChange={(checked) => onSelectRow(rowId, checked)}
                        containerProps={{
                            as: 'div',
                            py: 0,
                            px: 0,
                            flexBasis: '100%',
                            minWidth: '100%',
                            maxWidth: '100%',
                            h: '100%',
                        }}
                    />
                </Flex>
            ) : null}

            {onEdit ? (
                <IconButton
                    aria-label="Открыть карточку"
                    variant="ghost"
                    size="sm"
                    position="absolute"
                    right="5"
                    bottom="5"
                    boxSize="9"
                    zIndex={1}
                    onClick={(event) => {
                        event.stopPropagation();
                        onRowEdit(row);
                    }}
                >
                    <ArrowUpRight />
                </IconButton>
            ) : null}
        </Card.Root>
    );
}
