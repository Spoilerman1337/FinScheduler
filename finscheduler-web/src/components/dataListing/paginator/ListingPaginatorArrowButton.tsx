import {forwardRef} from 'react';
import {IconButton} from '@chakra-ui/react';
import {ChevronLeftIcon, ChevronRightIcon} from 'lucide-react';
import type {ComponentProps} from 'react';

interface ListingPaginatorArrowButtonProps extends ComponentProps<typeof IconButton> {
    direction: 'previous' | 'next';
}

const ListingPaginatorArrowButton = forwardRef<HTMLButtonElement, ListingPaginatorArrowButtonProps>(
    (props, ref) => {
        const {direction, ...restProps} = props;

        return (
            <IconButton
                ref={ref}
                color="neon.blue"
                borderColor="neon.blue"
                backdropFilter="blur(12px)"
                bg="glass.bgHover"
                border="1px solid"
                _hover={{
                    filter: 'drop-shadow(0 0 16px rgba(212, 0,255,0.9))',
                    boxShadow: '0 0 20px rgba(212, 0,255,1)',
                    color: 'neon.purple',
                    bg: 'glass.bgHover',
                    backdropFilter: 'blur(12px)',
                    borderColor: 'neon.purple',
                }}
                _disabled={{
                    color: 'neon.blue',
                    borderColor: 'glass.borderStrong',
                    bg: 'glass.bgHover',
                    opacity: 0.7,
                    filter: 'none',
                    boxShadow: 'none',
                    cursor: 'default',
                }}
                transition="all 0.3s ease-in-out"
                focusRing="none"
                {...restProps}
            >
                {direction === 'previous' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
            </IconButton>
        );
    },
);

ListingPaginatorArrowButton.displayName = 'ListingPaginatorArrowButton';

export default ListingPaginatorArrowButton;
