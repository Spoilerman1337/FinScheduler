import {Flex} from '@chakra-ui/react';
import ListingAddButton from './ListingAddButton.tsx';
import ListingDeleteButton from './ListingDeleteButton.tsx';
import ListingEditButton from './ListingEditButton.tsx';

interface ListingActionButtonsProps {
    selectedCount: number;
    onAdd?: () => void;
    onEditSelected?: () => void;
    onDeleteSelected?: () => void;
}

export default function ListingActionButtons(props: ListingActionButtonsProps) {
    const {selectedCount, onAdd, onEditSelected, onDeleteSelected} = props;
    const hasSelection = selectedCount > 0;

    return (
        <Flex
            mt={2}
            p={4}
            bg="bg.layer2"
            borderRadius="lg"
            border="1px solid"
            borderColor="glass.border"
            justify="flex-end"
            gap={2}
        >
            {onAdd && !hasSelection && <ListingAddButton onClick={onAdd} />}
            {onEditSelected && selectedCount === 1 && (
                <ListingEditButton onClick={onEditSelected} />
            )}
            {onDeleteSelected && selectedCount > 0 && (
                <ListingDeleteButton selectedCount={selectedCount} onClick={onDeleteSelected} />
            )}
        </Flex>
    );
}
