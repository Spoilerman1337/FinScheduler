import {Box, Input} from '@chakra-ui/react';
import {SearchIcon} from 'lucide-react';
import {filterWidthProps} from './shared.ts';

type TextInputFilterProps = {
    value: string;
    placeholder: string;
    onChange: (value: string) => void;
    onApply: () => void;
};

export default function TextInputFilter(props: TextInputFilterProps) {
    return (
        <Box {...filterWidthProps} position="relative">
            <Box
                position="absolute"
                left="3"
                top="50%"
                transform="translateY(-50%)"
                zIndex="1"
                pointerEvents="none"
            >
                <SearchIcon size={18} color="rgba(255,255,255,0.6)" />
            </Box>
            <Input
                placeholder={props.placeholder}
                value={props.value}
                onChange={(e) => props.onChange(e.target.value)}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        props.onApply();
                    }
                }}
                pl="10"
                bg="bg.layer1"
                borderColor="glass.border"
                color="neon.blue"
                _placeholder={{color: 'text.placeholder'}}
            />
        </Box>
    );
}
