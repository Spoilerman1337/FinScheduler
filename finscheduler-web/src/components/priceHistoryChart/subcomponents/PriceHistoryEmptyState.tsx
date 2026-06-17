import {Flex, Stack, Text} from '@chakra-ui/react';
import {TrendingUp} from 'lucide-react';

interface PriceHistoryEmptyStateProps {
    title: string;
    description: string;
}

export default function PriceHistoryEmptyState({
    title,
    description,
}: PriceHistoryEmptyStateProps) {
    return (
        <Flex
            minH="240px"
            borderRadius="xl"
            border="1px dashed"
            borderColor="app.cardBorder"
            bg="linear-gradient(180deg, rgba(16, 29, 49, 0.72), rgba(8, 16, 32, 0.78))"
            align="center"
            justify="center"
            px={5}
            py={6}
        >
            <Stack maxW="360px" align="center" textAlign="center" gap={3}>
                <Flex
                    width="14"
                    height="14"
                    borderRadius="full"
                    align="center"
                    justify="center"
                    bg="rgba(143, 120, 255, 0.14)"
                    color="neon.purple"
                    boxShadow="app.glowViolet"
                >
                    <TrendingUp size={22} />
                </Flex>
                <Text color="fg" fontWeight="600" textStyle="lg">
                    {title}
                </Text>
                <Text color="fg.subtle">{description}</Text>
            </Stack>
        </Flex>
    );
}
