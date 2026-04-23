import {Flex} from "@chakra-ui/react";
import DataTableAddButton from "./DataTableAddButton.tsx";
import DataTableDeleteButton from "./DataTableDeleteButton.tsx";
import DataTableEditButton from "./DataTableEditButton.tsx";

interface DataTableActionButtonsProps {
    selectedCount: number;
    onAdd?: () => void;
    onEditSelected?: () => void;
    onDeleteSelected?: () => void;
}

export default function DataTableActionButtons(props: DataTableActionButtonsProps) {
    const {selectedCount, onAdd, onEditSelected, onDeleteSelected} = props;

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
            {onAdd && <DataTableAddButton onClick={onAdd}/>}
            {onEditSelected && selectedCount === 1 && <DataTableEditButton onClick={onEditSelected}/>}
            {onDeleteSelected && selectedCount === 1 &&  (
                <DataTableDeleteButton
                    selectedCount={selectedCount}
                    onClick={onDeleteSelected}
                />
            )}
        </Flex>
    );
}
