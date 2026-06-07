import {Flex} from '@chakra-ui/react';
import type {ComponentProps, ReactNode} from 'react';

interface DataTableCellProps extends ComponentProps<typeof Flex> {
    children: ReactNode;
    isHeader?: boolean;
}

export default function DataTableCell(props: DataTableCellProps) {
    const {children, isHeader = false, ...restProps} = props;

    return (
        <Flex
            py={isHeader ? 4 : 3}
            px={5}
            flex={1}
            minWidth={isHeader ? '100px' : '0'}
            color={isHeader ? 'textMuted' : undefined}
            textTransform={isHeader ? 'uppercase' : undefined}
            fontSize={isHeader ? 'xs' : undefined}
            fontWeight={isHeader ? 'bold' : undefined}
            align={isHeader ? undefined : 'center'}
            overflow={isHeader ? undefined : 'hidden'}
            textOverflow={isHeader ? undefined : 'ellipsis'}
            whiteSpace={isHeader ? undefined : 'nowrap'}
            flexShrink={isHeader ? undefined : 1}
            {...restProps}
        >
            {children}
        </Flex>
    );
}
