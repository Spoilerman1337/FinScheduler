import {IconButton} from "@chakra-ui/react";
import type {ComponentProps} from "react";

interface PaginatorButtonProps extends ComponentProps<typeof IconButton> {
    interactive?: boolean;
}

export default function PaginatorButton(props: PaginatorButtonProps) {
    const {interactive = true, ...restProps} = props;

    return (
        <IconButton
            color="neon.blue"
            borderColor="neon.blue"
            backdropFilter="blur(12px)"
            bg="glass.bgHover"
            border="1px solid"
            transition="all 0.3s ease-in-out"
            focusRing="none"
            cursor={interactive ? "pointer" : "default"}
            pointerEvents={interactive ? "auto" : "none"}
            _hover={interactive ? {
                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                color: "neon.purple",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
                borderColor: "neon.purple",
            } : {
                color: "neon.blue",
                bg: "glass.bgHover",
                borderColor: "neon.blue",
            }}
            _disabled={{
                color: "textMuted",
                borderColor: "glass.border",
                bg: "bg.layer2",
                opacity: 1,
                filter: "none",
                boxShadow: "none",
                cursor: "default",
            }}
            {...restProps}
        />
    );
}
