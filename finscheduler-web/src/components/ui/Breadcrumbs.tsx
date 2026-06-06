import {Fragment} from 'react';
import {HStack, Link, Text} from '@chakra-ui/react';
import {ChevronRight} from 'lucide-react';
import {Link as RouterLink} from 'react-router-dom';

export interface BreadcrumbItem {
    label: string;
    to?: string;
}

interface BreadcrumbsProps {
    items: BreadcrumbItem[];
}

export default function Breadcrumbs({items}: BreadcrumbsProps) {
    return (
        <HStack gap={2} flexWrap="wrap" color="fg.muted" textStyle="sm">
            {items.map((item, index) => {
                const isLast = index === items.length - 1;

                return (
                    <Fragment key={`${item.label}-${index}`}>
                        {item.to && !isLast ? (
                            <Link asChild color="fg.muted" _hover={{color: 'app.accent'}}>
                                <RouterLink to={item.to}>{item.label}</RouterLink>
                            </Link>
                        ) : (
                            <Text color={isLast ? 'fg' : 'fg.muted'}>{item.label}</Text>
                        )}
                        {!isLast && <ChevronRight size={14} />}
                    </Fragment>
                );
            })}
        </HStack>
    );
}
