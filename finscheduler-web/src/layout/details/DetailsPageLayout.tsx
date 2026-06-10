import {Badge, Box, Button, Card, Flex, Stack, Text} from '@chakra-ui/react';
import {CheckCircle2, Save, X} from 'lucide-react';
import type {ReactNode} from 'react';
import Breadcrumbs, {type BreadcrumbItem} from '../../components/ui/Breadcrumbs.tsx';

interface DetailsPageLayoutProps {
    breadcrumbItems: BreadcrumbItem[];
    title: string;
    subtitle: string;
    isActive: boolean;
    isDirty: boolean;
    saving: boolean;
    error?: string | null;
    onSave: () => void;
    onSaveAndClose: () => void;
    onBack: () => void;
    children: ReactNode;
}

export default function DetailsPageLayout({
    breadcrumbItems,
    title,
    subtitle,
    isActive,
    isDirty,
    saving,
    error = null,
    onSave,
    onSaveAndClose,
    onBack,
    children,
}: DetailsPageLayoutProps) {
    return (
        <Stack width="100%" gap={6} pb={6}>
            <Breadcrumbs items={breadcrumbItems} />

            <Card.Root overflow="visible">
                <Box
                    position="absolute"
                    inset="0"
                    pointerEvents="none"
                    borderRadius="inherit"
                    background="linear-gradient(90deg, rgba(32, 208, 255, 0.08), rgba(143, 120, 255, 0.04))"
                />
                <Card.Body position="relative" gap={6} pb={{base: 7, xl: 6}}>
                    <Flex
                        direction={{base: 'column', xl: 'row'}}
                        align="flex-start"
                        justify="space-between"
                        gap={6}
                    >
                        <Stack gap={2} flex="1" minW={0}>
                            <Text textStyle="3xl" color="fg" fontWeight="700" letterSpacing="tight" lineHeight="tight">
                                {title}
                            </Text>
                            <Text color="fg.muted" textStyle="lg" lineHeight="snug">
                                {subtitle}
                            </Text>
                            <Badge
                                alignSelf="flex-start"
                                px={3}
                                py={1}
                                borderRadius="full"
                                bg={isActive ? 'neon.green' : 'neon.pink'}
                                color="bg.base"
                            >
                                {isActive ? 'Активен' : 'Неактивен'}
                            </Badge>
                        </Stack>

                        <Flex wrap="wrap" gap={3}>
                            {isDirty ? (
                                <>
                                    <Button onClick={onSave} loading={saving}>
                                        <Save />
                                        Сохранить
                                    </Button>
                                    <Button
                                        variant="surface"
                                        borderColor="app.cardBorderActive"
                                        boxShadow="app.glowViolet"
                                        _hover={{
                                            bg: 'rgba(143, 120, 255, 0.18)',
                                            borderColor: 'app.cardBorderActive',
                                        }}
                                        onClick={onSaveAndClose}
                                        loading={saving}
                                    >
                                        <CheckCircle2 />
                                        Сохранить и закрыть
                                    </Button>
                                    <Button variant="outline" onClick={onBack} disabled={saving}>
                                        <X />
                                        Отмена
                                    </Button>
                                </>
                            ) : (
                                <Button variant="outline" onClick={onBack} disabled={saving}>
                                    <X />
                                    Назад
                                </Button>
                            )}
                        </Flex>
                    </Flex>
                </Card.Body>
            </Card.Root>

            {error ? (
                <Card.Root borderColor="border.error">
                    <Card.Body>
                        <Text color="fg.error">{error}</Text>
                    </Card.Body>
                </Card.Root>
            ) : null}

            {children}
        </Stack>
    );
}
