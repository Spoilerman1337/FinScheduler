import {Box, Input} from "@chakra-ui/react";
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
            <Input
                placeholder={props.placeholder}
                type="number"
                value={props.value}
                onChange={(e) => props.onChange(e.target.value)}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        props.onApply();
                    }
                }}
                bg="bg.layer1"
                borderColor="glass.border"
                color="neon.blue"
                _placeholder={{color: 'text.placeholder'}}
            />
        </Box>
    );
}
