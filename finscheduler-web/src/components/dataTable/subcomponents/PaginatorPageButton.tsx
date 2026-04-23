import {forwardRef} from "react";
import {IconButton, Text} from "@chakra-ui/react";
import type {ComponentProps} from "react";

interface PaginatorPageButtonProps extends ComponentProps<typeof IconButton> {
    value: number;
}

const PaginatorPageButton = forwardRef<HTMLButtonElement, PaginatorPageButtonProps>((props, ref) => {
    const {value, ...restProps} = props;

    return (
        <IconButton
            ref={ref}
            color={"neon.blue"}
            borderColor={"neon.blue"}
            backdropFilter={"blur(12px)"}
            bg={"glass.bgHover"}
            border={"1px solid"}
            transition="all 0.3s ease-in-out"
            _selected={{
                filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                boxShadow: "0 0 20px rgba(0,212,255,1)",
                color: "neon.blue",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
                borderColor: "neon.blue",
            }}
            _hover={{
                filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                boxShadow: "0 0 20px rgba(212, 0,255,1)",
                color: "neon.purple",
                bg: "glass.bgHover",
                backdropFilter: "blur(12px)",
                borderColor: "neon.purple",
            }}
            focusRing={"none"}
            {...restProps}
        >
            <Text color="currentColor">{value}</Text>
        </IconButton>
    );
});

PaginatorPageButton.displayName = "PaginatorPageButton";

export default PaginatorPageButton;
