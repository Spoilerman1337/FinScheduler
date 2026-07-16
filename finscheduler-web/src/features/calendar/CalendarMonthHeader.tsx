import {Button, DatePicker, Flex, HStack, Text, VStack} from '@chakra-ui/react';
import {ChevronLeft, ChevronRight} from 'lucide-react';
import {formatMonthLabel} from './shared.ts';

interface CalendarMonthHeaderProps {
    visibleMonthStart: Date;
    onTodayClick: () => void;
}

export default function CalendarMonthHeader(props: CalendarMonthHeaderProps) {
    const {visibleMonthStart, onTodayClick} = props;

    return (
        <Flex justify="space-between" align={{base: 'start', lg: 'center'}} gap="4" wrap="wrap">
            <VStack align="start" gap="1">
                <Text fontSize={{base: '2xl', md: '3xl'}} fontWeight="700">
                    {formatMonthLabel(visibleMonthStart)}
                </Text>
            </VStack>

            <HStack gap="2" wrap="wrap">
                <DatePicker.PrevTrigger
                    unstyled
                    display="inline-flex"
                    alignItems="center"
                    justifyContent="center"
                    h="9"
                    minW="9"
                    px="3.5"
                    borderWidth="1px"
                    borderColor="app.cardBorder"
                    borderRadius="l2"
                    bg="app.cardBg"
                    color="fg"
                    boxShadow="card"
                    _hover={{bg: 'app.cardBgHover', borderColor: 'app.cardBorderActive'}}
                >
                    <ChevronLeft size={16} />
                </DatePicker.PrevTrigger>
                <Button variant="outline" onClick={onTodayClick}>
                    {'Сегодня'}
                </Button>
                <DatePicker.NextTrigger
                    unstyled
                    display="inline-flex"
                    alignItems="center"
                    justifyContent="center"
                    h="9"
                    minW="9"
                    px="3.5"
                    borderWidth="1px"
                    borderColor="app.cardBorder"
                    borderRadius="l2"
                    bg="app.cardBg"
                    color="fg"
                    boxShadow="card"
                    _hover={{bg: 'app.cardBgHover', borderColor: 'app.cardBorderActive'}}
                >
                    <ChevronRight size={16} />
                </DatePicker.NextTrigger>
            </HStack>
        </Flex>
    );
}
