import {Box, NumberInput} from "@chakra-ui/react";
import {filterWidthProps} from "./shared.ts";

type NumberInputFilterProps = {
    value: string;
    placeholder: string;
    onChange: (value: string) => void;
    onApply: () => void;
};

export default function NumberInputFilter(props: NumberInputFilterProps) {
    return (
        <Box {...filterWidthProps}>
            <NumberInput.Root
                value={props.value}
                onValueChange={(details) => props.onChange(details.value)}
                bg="bg.layer1"
                borderColor="glass.border"
                color="neon.blue"
                width="100%"
                gap={0}
            >
                <NumberInput.Control>
                    <NumberInput.IncrementTrigger
                        bg="bg.layer1"
                        color="neon.blue"
                        p={0}
                        _hover={{
                            filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                            boxShadow: "0 0 12px rgba(0, 212, 255, 0.45)",
                        }}
                    >
                        +
                    </NumberInput.IncrementTrigger>
                    <NumberInput.DecrementTrigger
                        bg="bg.layer1"
                        color="neon.blue"
                        p={0}
                        _hover={{
                            filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                            boxShadow: "0 0 12px rgba(0, 212, 255, 0.45)",
                        }}
                    >
                        -
                    </NumberInput.DecrementTrigger>
                </NumberInput.Control>
                <NumberInput.Input
                    placeholder={props.placeholder}
                    _placeholder={{color: 'text.placeholder'}}
                    onKeyDown={(event) => {
                        if (event.key === 'Enter') {
                            props.onApply();
                        }
                    }}
                />
            </NumberInput.Root>
        </Box>
    );
}
