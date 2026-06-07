import {Checkbox, Flex} from '@chakra-ui/react';
import type {ComponentProps, MouseEventHandler} from 'react';

interface ListingSelectionCheckboxProps {
    checked: boolean;
    onCheckedChange: (checked: boolean) => void;
    isHeader?: boolean;
    onClick?: MouseEventHandler<HTMLElement>;
    containerProps?: ComponentProps<typeof Flex>;
}

export default function ListingSelectionCheckbox(props: ListingSelectionCheckboxProps) {
    const {checked, onCheckedChange, isHeader = false, onClick, containerProps} = props;

    return (
        <Flex
            as={isHeader ? 'th' : 'td'}
            py={isHeader ? 4 : 3}
            px={5}
            flexBasis="50px"
            minWidth="50px"
            maxWidth="50px"
            align="center"
            justify="center"
            {...containerProps}
            onClick={onClick ?? containerProps?.onClick}
        >
            <Checkbox.Root
                checked={checked}
                onCheckedChange={(details) => {
                    onCheckedChange(details.checked === true);
                }}
            >
                <Checkbox.HiddenInput />
                <Checkbox.Control
                    filter="drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"
                    transition="all 0.3s ease-in-out"
                    color="neon.blue"
                    borderColor="neon.blue"
                    _hover={{
                        boxShadow: '0 0 12px rgba(0, 212, 255, 0.9)',
                        filter: 'drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))',
                    }}
                >
                    <Checkbox.Indicator />
                </Checkbox.Control>
            </Checkbox.Root>
        </Flex>
    );
}
