import {Button} from "@chakra-ui/react";
import {TrashIcon} from "lucide-react";

interface DataTableDeleteButtonProps {
    selectedCount: number;
    onClick: () => void;
}

export default function DataTableDeleteButton(props: DataTableDeleteButtonProps) {
    const {selectedCount, onClick} = props;

    return (
        <Button
            onClick={onClick}
            color="neon.blue"
            borderColor="neon.blue"
            backdropFilter="blur(12px)"
            bg="glass.bgHover"
            border="1px solid"
            transition="all 0.3s ease-in-out"
            disabled={selectedCount === 0}
            opacity={selectedCount === 0 ? 0.5 : 1}
            _hover={selectedCount > 0 ? {
                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                color: "neon.purple",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
                borderColor: "neon.purple",
            } : {}}
            focusRing="none"
        >
            <TrashIcon size={18} style={{marginRight: "8px"}}/>
            Удалить {selectedCount > 0 ? `(${selectedCount})` : ""}
        </Button>
    );
}
