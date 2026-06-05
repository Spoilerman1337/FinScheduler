import {IconButton} from "@chakra-ui/react";
import {EllipsisIcon} from "lucide-react";

export default function ListingPaginatorEllipsisButton() {
    return (
        <IconButton
            color="neon.blue"
            borderColor="neon.blue"
            backdropFilter="blur(12px)"
            bg="glass.bgHover"
            border="1px solid"
            disabled
            transition="all 0.3s ease-in-out"
            _selected={{
                filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                boxShadow: "0 0 20px rgba(0,212,255,1)",
                color: "neon.blue",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
            }}
            _hover={{
                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                color: "neon.purple",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
                borderColor: "neon.purple",
            }}
            _disabled={{
                color: "neon.blue",
                borderColor: "glass.borderStrong",
                bg: "glass.bgHover",
                opacity: 0.7,
                filter: "none",
                boxShadow: "none",
                cursor: "default",
            }}
            focusRing="none"
        >
            <EllipsisIcon/>
        </IconButton>
    );
}
