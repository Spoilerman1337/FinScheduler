import {Box, Button, Circle, Flex, HStack, Stack, Text} from '@chakra-ui/react';
import {ChevronsLeft, ChevronsRight} from 'lucide-react';
import {useState} from 'react';
import SvgLogoIcon from '../../components/svgIcon/SvgLogoIcon.tsx';
import {routedNavigationItems} from './navigation.tsx';
import SidebarItem from './subcomponents/SidebarItem.tsx';

export default function Sidebar() {
    const [isCollapsed, setIsCollapsed] = useState(false);

    return (
        <Box
            display={{base: 'none', md: 'block'}}
            w={{md: '5.5rem', lg: isCollapsed ? '5.5rem' : '15.5rem'}}
            transition="width 0.25s ease"
            flexShrink={0}
            borderRightWidth="1px"
            borderColor="app.cardBorder"
            bg="app.sidebarBg"
            backdropFilter="blur(22px)"
            boxShadow="card"
            position="sticky"
            top="0"
            h="100vh"
            overflow="hidden"
        >
            <Flex direction="column" h="full" px={{md: '3', lg: isCollapsed ? '3' : '4'}} py="5">
                <HStack
                    align="center"
                    gap="3"
                    minH="12"
                    px={{md: '1.5', lg: isCollapsed ? '1.5' : '2'}}
                    color="app.accent"
                >
                    <Circle
                        size="10"
                        borderWidth="1px"
                        borderColor="app.cardBorderActive"
                        bg="app.accentSoft"
                        boxShadow="app.glowCyan"
                        flexShrink={0}
                    >
                        <Box boxSize="6">
                            <SvgLogoIcon />
                        </Box>
                    </Circle>

                    <Stack
                        gap="0"
                        display={{base: 'none', md: 'none', lg: isCollapsed ? 'none' : 'flex'}}
                    >
                        <Text
                            fontSize="2xl"
                            fontWeight="700"
                            color="app.accent"
                            letterSpacing="tight"
                        >
                            FinScheduler
                        </Text>
                    </Stack>
                </HStack>

                <Stack gap="1.5" mt="8" flex="1">
                    {routedNavigationItems.map((item) => (
                        <SidebarItem
                            key={item.id}
                            icon={item.icon}
                            label={item.label}
                            path={item.path}
                            disabled={item.disabled}
                            compact={isCollapsed}
                        />
                    ))}
                </Stack>

                <HStack mt="4" justify="center" gap="2">
                    <Button
                        variant="ghost"
                        size="sm"
                        justifyContent="flex-start"
                        color="fg.muted"
                        onClick={() => setIsCollapsed(!isCollapsed)}
                    >
                        <HStack gap="2">
                            {isCollapsed ? (
                                <ChevronsRight size={16} />
                            ) : (
                                <>
                                    <Text
                                        fontSize="sm"
                                        display={{
                                            base: 'none',
                                            md: 'none',
                                            lg: isCollapsed ? 'none' : 'inline-flex',
                                        }}
                                    >
                                        Свернуть меню
                                    </Text>
                                    <ChevronsLeft size={16} />
                                </>
                            )}
                        </HStack>
                    </Button>
                </HStack>
            </Flex>
        </Box>
    );
}
