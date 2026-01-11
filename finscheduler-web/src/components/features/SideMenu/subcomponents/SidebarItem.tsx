'use client';

import type { ReactNode } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Box, Text } from '@chakra-ui/react';
import 'focus-visible/dist/focus-visible'
import {Link, useLocation} from "react-router-dom";

interface SidebarItemProps {
    icon: ReactNode;
    label: string;
    isCollapsed: boolean;
    onClick?: () => void;
    id: string;
    path: string;
}

export default function SidebarItem(props: SidebarItemProps) {
    const location = useLocation();
    const isActive = location.pathname === props.path;

    const handleClick = () => {
        props.onClick?.()
    }

    return (
        <Link to={props.path}>
            <Box
                as="button"
                w="full"
                h="36px"
                display="flex"
                alignItems="center"
                gap={3}
                px={2}
                borderRadius="xl"
                cursor="pointer"
                position="relative"
                overflow="hidden"

                color={isActive ? "text.primary" : "text.secondary"}
                fontWeight={isActive ? "semibold" : "medium"}
                bg={isActive ? "rgba(255,255,255,0.08)" : "transparent"}
                backdropFilter={isActive ? "blur(12px)" : undefined}
                style={{
                    textShadow: isActive
                        ? "0 0 16px rgba(0,212,255,0.8)"
                        : undefined,
                }}
                _hover={{
                    filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                    textShadow: "0 0 20px rgba(0,212,255,1)",
                    color: "white",
                    bg: "rgba(255,255,255,0.12)",
                    backdropFilter: "blur(12px)",
                }}
                focusRing={"none"}

                transition="all 0.25s cubic-bezier(0.4, 0, 0.2, 1)"
                onClick={handleClick}
            >
                {/* Иконка всегда на месте */}
                <Box flexShrink={0} w={6} h={6} display="flex" alignItems="center" justifyContent="center" color={isActive ? "neon.blue" : "text.primary"}>
                    {props.icon}
                </Box>

                {/* Текст плавно появляется/исчезает */}
                <AnimatePresence>
                    {!props.isCollapsed && (
                        <motion.div
                            initial={{ opacity: 0, width: 0 }}
                            animate={{ opacity: 1, width: 'auto' }}
                            exit={{ opacity: 0, width: 0 }}
                            transition={{ duration: 0.2, ease: 'easeInOut' }}
                            style={{ overflow: 'hidden', whiteSpace: 'nowrap' }}
                        >
                            <Text fontSize="sm" fontWeight={isActive ? "semibold" : "medium"} fontFamily={"Montserrat"} color={isActive ? "neon.blue" : "text.primary"}>
                                {props.label}
                            </Text>
                        </motion.div>
                    )}
                </AnimatePresence>
                {isActive && (
                    <Box
                        position="absolute"
                        left={0}
                        top="4px"
                        bottom="4px"
                        w="3px"
                        bg="neon.blue"
                        borderRadius="full"
                        boxShadow="0 0 12px rgba(0,212,255,0.8)"
                    />
                )}
            </Box>
        </Link>
    );
};