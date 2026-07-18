import {Box, DatePicker, Text, useDatePickerContext} from '@chakra-ui/react';
import type {CalendarDayDetails} from './shared.ts';
import CalendarDayCell from './CalendarDayCell.tsx';

interface CalendarGridPanelProps {
    detailsByDay: Record<string, CalendarDayDetails>;
}

interface CalendarDayGridProps {
    detailsByDay: Record<string, CalendarDayDetails>;
}

function CalendarDayGrid(props: CalendarDayGridProps) {
    const {detailsByDay} = props;
    const datePicker = useDatePickerContext();

    return (
        <DatePicker.Table
            unstyled
            aria-label={'Календарная сетка'}
            w="full"
            tableLayout="fixed"
            borderCollapse="separate"
            style={{borderSpacing: '0.5rem'}}
        >
            <DatePicker.TableHead unstyled>
                <DatePicker.TableRow unstyled>
                    {datePicker.weekDays.map((weekDay) => (
                        <DatePicker.TableHeader
                            key={weekDay.short}
                            unstyled
                            px="2"
                            py="2.5"
                            borderRadius="lg"
                            bg="rgba(6, 16, 34, 0.3)"
                            textAlign="center"
                        >
                            <Text
                                fontSize="xs"
                                fontWeight="700"
                                color="fg.subtle"
                                textTransform="uppercase"
                                letterSpacing="0.08em"
                            >
                                {weekDay.short}
                            </Text>
                        </DatePicker.TableHeader>
                    ))}
                </DatePicker.TableRow>
            </DatePicker.TableHead>

            <DatePicker.TableBody unstyled>
                {datePicker.weeks.map((week, weekIndex) => (
                    <DatePicker.TableRow key={weekIndex} unstyled>
                        {week.map((dayValue) => (
                            <CalendarDayCell
                                key={dayValue.toString()}
                                dayValue={dayValue}
                                detailsByDay={detailsByDay}
                            />
                        ))}
                    </DatePicker.TableRow>
                ))}
            </DatePicker.TableBody>
        </DatePicker.Table>
    );
}

export default function CalendarGridPanel(props: CalendarGridPanelProps) {
    const {detailsByDay} = props;

    return (
        <Box
            borderWidth="1px"
            borderColor="app.cardBorder"
            bg="app.cardBg"
            borderRadius="2xl"
            boxShadow="card"
            px={{base: '4', md: '5'}}
            py={{base: '4', md: '5'}}
            display="flex"
            flexDirection="column"
        >
            <Box minH={{base: '32rem', md: '38rem', xl: '42rem'}}>
                <CalendarDayGrid detailsByDay={detailsByDay} />
            </Box>
        </Box>
    );
}
