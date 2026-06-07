import {Button, Flex} from '@chakra-ui/react';
import type {ReactNode} from 'react';
import {filterWidthProps} from './shared.ts';

type FilterWrapperProps = {
    children: ReactNode;
    onReset: () => void;
};

export default function FilterWrapper(props: FilterWrapperProps) {
    return (
        <Flex
            mb={4}
            p={4}
            bg="bg.layer2"
            borderRadius="lg"
            border="1px solid"
            borderColor="glass.border"
            width="100%"
            align="center"
            gap={4}
            flexWrap="wrap"
            justifyContent="flex-start"
        >
            {props.children}
            <Button
                {...filterWidthProps}
                onClick={props.onReset}
                bg="glass.bgHover"
                color="neon.purple"
                borderColor="neon.purple"
                backdropFilter="blur(12px)"
                border="1px solid"
                transition="all 0.3s ease-in-out"
                _hover={{
                    filter: 'drop-shadow(0 0 16px rgba(212, 0,255,0.9))',
                    boxShadow: '0 0 20px rgba(212, 0,255,1)',
                    color: 'neon.pink',
                    bg: 'glass.bgHover',
                    backdropFilter: 'blur(12px)',
                    borderColor: 'neon.purple',
                }}
                focusRing="none"
            >
                Сброс
            </Button>
        </Flex>
    );
}
