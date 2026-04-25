import {Button} from "@chakra-ui/react";
import {PencilIcon} from "lucide-react";

interface DataTableEditButtonProps {
    onClick: () => void;
}

export default function DataTableEditButton(props: DataTableEditButtonProps) {
    const {onClick} = props;

    return (
        <Button
            onClick={onClick}
            bg="neon.green"
            color="bg.base"
            _hover={{bg: "neon.green", opacity: 0.8}}
        >
            <PencilIcon size={18} style={{marginRight: "8px"}}/>
            Редактировать
        </Button>
    );
}
