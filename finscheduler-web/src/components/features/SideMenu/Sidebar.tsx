import {Box, Flex, IconButton, Text, VStack} from "@chakra-ui/react";
import {AnimatePresence, motion} from "framer-motion";
import {useState} from "react";
import {
    ChevronsRight,
    LayoutDashboardIcon,
    ShoppingCartIcon,
} from "lucide-react";
import SvgLogoIcon from "../svgIcon/SvgLogoIcon.tsx";
import SidebarItem from "./subcomponents/SidebarItem.tsx";

export default function Sidebar() {
    const [isCollapsed, setIsCollapsed] = useState(false);

    const sidebarWidth = isCollapsed ? '68px' : '240px';

    return (
        <motion.div
            animate={{width: sidebarWidth}}
            transition={{duration: 0.3, ease: 'easeInOut'}}
            style={{
                height: '100vh',
                background: 'linear-gradient(to bottom, #1a1a2e, #16213e)',
                position: 'relative',
                overflow: 'hidden',
            }}
        >
            <Flex direction="column" h="full" color="white" bg={"bg.layer1"}>
                {/* Логотип / Заголовок */}
                <Box p={3} pb={6}>
                    <Flex align="center" gap={3}>
                        <SvgLogoIcon />
                        <AnimatePresence>
                            {!isCollapsed && (
                                <motion.div
                                    initial={{opacity: 0}}
                                    animate={{opacity: 1}}
                                    exit={{opacity: 0}}
                                    transition={{duration: 0.2}}
                                >
                                <Text fontSize="lg" fontWeight="bold">
                                        FinScheduler
                                    </Text>
                                </motion.div>
                            )}
                        </AnimatePresence>
                    </Flex>
                </Box>

                <VStack align="stretch" flex={1} px={2} gap={1}>
                    <SidebarItem
                        isCollapsed={isCollapsed}
                        icon={<LayoutDashboardIcon/>}
                        label="Дашборды"
                        id={'dashboard'}
                        path={"/"}
                    />
                    <SidebarItem
                        isCollapsed={isCollapsed}
                        icon={<ShoppingCartIcon/>}
                        label="Виды расходов"
                        id={'items'}
                        path={"/items"}
                    />
                </VStack>

                <Box position="relative" h="96px" mt="auto">
                    <Box position="absolute" right={2} bottom={3}>
                        <IconButton
                            aria-label="Свернуть/развернуть"
                            size="sm"
                            borderRadius="full"
                            variant="ghost"
                            onClick={() => setIsCollapsed(!isCollapsed)}
                            focusRing={"none"}
                            color={isCollapsed ? "neon.blue" : "text.primary"}
                            borderColor={isCollapsed ? "neon.blue" : "text.primary"}
                            backdropFilter={isCollapsed ? "blur(12px)" : undefined}
                            bg={isCollapsed ? "rgba(255,255,255,0.12)" : "bg.layer2"}
                            filter={isCollapsed ? "drop-shadow(0 0 16px rgba(0,212,255,0.9))" : undefined}
                            textShadow={isCollapsed ? "0 0 20px rgba(0,212,255,1)" : undefined}
                            _hover={{
                                filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                                textShadow: "0 0 20px rgba(0,212,255,1)",
                                color: "white",
                                bg: "rgba(255,255,255,0.12)",
                                backdropFilter: "blur(12px)",
                            }}
                        >
                            <motion.div
                                animate={{rotate: isCollapsed ? 180 : 0}}
                                transition={{duration: 0.3}}
                            >
                                <ChevronsRight size={18}/>
                            </motion.div>
                        </IconButton>
                    </Box>
                </Box>
            </Flex>
        </motion.div>
    );
}