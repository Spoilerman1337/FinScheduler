import {Box, HStack, Icon, Text} from '@chakra-ui/react';
import type {LucideIcon} from 'lucide-react';
import {NavLink, useLocation} from 'react-router-dom';

export interface SidebarItemProps {
    icon: LucideIcon;
    label: string;
    path?: string;
    disabled?: boolean;
    compact?: boolean;
}

function SidebarItemContent(props: SidebarItemProps & {isActive: boolean}) {
    const {icon, label, compact = false, isActive, disabled = false} = props;

    return (
        <HStack
            w="full"
            minH="12"
            px={{md: '3', lg: compact ? '3' : '4'}}
            py="3"
            gap="3"
            borderRadius="l2"
            borderWidth="1px"
            borderColor={isActive ? 'app.cardBorderActive' : 'transparent'}
            bg={isActive ? 'app.accentSoft' : 'transparent'}
            color={isActive ? 'app.accent' : 'fg.muted'}
            opacity={disabled ? 0.56 : 1}
            transition="all 0.2s ease"
            _hover={{
                bg: disabled ? 'transparent' : 'transparent',
                bgGradient: disabled
                    ? undefined
                    : 'linear(to-r, app.sidebarItemHoverStart, app.sidebarItemHoverEnd)',
                color: disabled ? 'fg.muted' : 'fg',
                borderColor: disabled ? 'transparent' : 'app.sidebarItemHoverBorder',
                boxShadow: disabled ? 'none' : 'app.sidebarItemHover',
            }}
        >
            <Box
                boxSize="6"
                display="flex"
                alignItems="center"
                justifyContent="center"
                flexShrink={0}
            >
                <Icon as={icon} boxSize="4.5" />
            </Box>

            <Text
                display={{base: 'none', md: 'none', lg: compact ? 'none' : 'block'}}
                fontSize="sm"
                fontWeight={isActive ? '700' : '600'}
                color="currentColor"
                whiteSpace="nowrap"
            >
                {label}
            </Text>
        </HStack>
    );
}

export default function SidebarItem(props: SidebarItemProps) {
    const location = useLocation();
    const isActive = Boolean(props.path) && location.pathname === props.path;

    if (!props.path || props.disabled) {
        return (
            <Box as="button" w="full" textAlign="left" cursor="not-allowed" aria-disabled="true">
                <SidebarItemContent {...props} isActive={false} />
            </Box>
        );
    }

    return (
        <NavLink to={props.path} aria-label={props.label}>
            <SidebarItemContent {...props} isActive={isActive} />
        </NavLink>
    );
}
