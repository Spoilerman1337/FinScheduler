import {Button} from "@chakra-ui/react";
import {PlusIcon} from "lucide-react";

interface DataTableAddButtonProps {
    onClick: () => void;
}

export default function DataTableAddButton(props: DataTableAddButtonProps) {
    const {onClick} = props;

    return (
        <Button
            onClick={onClick}
            bg="neon.blue"
            color="bg.base"
            _hover={{bg: "neon.blue", opacity: 0.8}}
        >
            <PlusIcon size={18} style={{marginRight: "8px"}}/>
            Добавить
        </Button>
    );
}
