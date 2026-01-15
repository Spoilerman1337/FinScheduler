import {Box, ButtonGroup, createListCollection, Flex, IconButton, Pagination, Select, Text} from "@chakra-ui/react";
import {CheckIcon, ChevronDownIcon, ChevronLeftIcon, ChevronRightIcon, EllipsisIcon} from "lucide-react";

export interface PaginatorProps  {
    total: number;
    page: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
}

export default function Paginator(props: PaginatorProps) {
    const items = createListCollection({
        items: [
            { label: "10", value: 10 },
            { label: "25", value: 25 },
            { label: "50", value: 50 },
            { label: "100", value: 100 },
        ],
    })

    const totalPages = Math.ceil(props.total / props.pageSize);

    return (
        <Flex justify="space-between" align="center" mt={2} py={5}>
            <Select.Root
                collection={items}
                width="60px"
                size="lg"
                value={[props.pageSize.toString()]}
                onValueChange={(e) => {
                    const newPageSize = parseInt(e.value[0] || '10');
                    props.onPageSizeChange(newPageSize);
                }}
            >
                <Select.HiddenSelect />
                <Select.Label
                    color="neon.blue"
                    mb="2"
                    fontWeight="medium"
                    _empty={{ display: "none" }}
                />

                {/* CONTROL/TRIGGER - Стилизован как неоновая кнопка/пагинатор */}
                <Select.Control asChild>
                    <Flex
                        align="center"
                        justify="space-between"

                    >
                        <Select.Trigger asChild>
                            <Box flex="1"
                                 borderColor="neon.blue"
                                 boxShadow="0 0 10px rgba(0, 212, 255, 0.6)"
                                 transition="all 0.3s ease-in-out"
                                 backdropFilter="blur(12px)"
                                 bg="glass.bg"
                                 color="neon.blue"
                                 p="1.5"
                                 maxW="100px"
                                 h="36px"
                                 focusRing="none"
                                 _hover={{
                                     filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                     boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                     color: "neon.purple",
                                     bg: "glass.bgHover",
                                     borderColor: "neon.purple", // Рамка тоже меняет цвет на ховере
                                     cursor: "pointer",
                                 }}>
                                <Select.ValueText
                                    placeholder="10" // Обычно тут отображается число
                                    color="textPrimary"
                                    _placeholder={{ color: 'textMuted' }}
                                    // 👇 Уменьшаем размер текста
                                    fontSize="sm"
                                />
                            </Box>
                        </Select.Trigger>

                        <Select.IndicatorGroup asChild>
                            <Flex align="center">
                                <Select.Indicator asChild>
                                    <ChevronDownIcon />
                                </Select.Indicator>
                            </Flex>
                        </Select.IndicatorGroup>
                    </Flex>
                </Select.Control>

                {/* ПОЗИЦИОНЕР И КОНТЕНТ: Уменьшаем отступы и размеры шрифта */}
                <Select.Positioner>
                    <Select.Content
                        // ... (Стили контейнера остаются прежними, но могут использовать меньший padding/radius)
                        p="1" // Меньший внутренний отступ
                        borderRadius="sm"
                        zIndex="dropdown"
                        mt="1" // Меньший отступ сверху
                        border="1px solid"
                        borderColor="glass.borderStrong"
                        backdropFilter="blur(16px)"
                        bg="glass.bg"
                        boxShadow="lg"
                        width="--trigger-width"
                        maxH="200px" // Меньшая высота
                        overflowY="auto"
                    >
                        {items.items.map((item) => (
                            <Select.Item
                                item={item}
                                key={item.value}
                                display="flex"
                                justifyContent="space-between"
                                alignItems="center"
                                // 👇 Меньшие отступы для элементов списка
                                py="1.5"
                                px="2"
                                borderRadius="sm"
                                color="neon.blue"
                                fontSize="sm" // Меньший размер текста
                                transition="all 0.2s"
                                cursor="pointer"
                                _hover={{
                                    filter: "drop-shadow(0 0 8px rgba(212, 0,255,0.9))",
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                }}
                                _selected={{
                                    filter: "drop-shadow(0 0 8px rgba(0,212,255,0.9))",
                                    bg: "glass.bgHover",
                                    fontWeight: "semibold",
                                }}
                                _focus={{
                                    outline: "none",
                                    bg: "glass.bgHover",
                                    boxShadow: "0 0 0 2px rgba(0, 212, 255, 0.5)",
                                }}
                            >
                                <Text>{item.label}</Text>
                                <Select.ItemIndicator asChild>
                                    <CheckIcon />
                                </Select.ItemIndicator>
                            </Select.Item>
                        ))}
                    </Select.Content>
                </Select.Positioner>
            </Select.Root>

            <Text fontSize="lg" color="white" fontFamily="body">Элементов всего: {props.total}</Text>

            <Pagination.Root 
                count={totalPages} 
                pageSize={1} 
                page={props.page} 
                onPageChange={(e) => props.onPageChange(e.page)}
                key={totalPages}
            >
                <ButtonGroup variant="ghost" size="lg" my={-5}>
                    <Pagination.PrevTrigger asChild>
                        <IconButton color={"neon.blue"}
                                    borderColor={"neon.blue"}
                                    backdropFilter={"blur(12px)"}
                                    bg={"glass.bgHover"}
                                    _hover={{
                                        filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                        boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                        color: "neon.purple",
                                        bg: "glass.bgHover",
                                        backdropFilter: "blur(12px)",
                                        borderColor: "neon.purple",
                                    }}
                                    transition="all 0.3s ease-in-out"
                                    focusRing={"none"}>
                            <ChevronLeftIcon />
                        </IconButton>
                    </Pagination.PrevTrigger>

                    <Pagination.Items
                        render={(page) => (
                            <IconButton color={"neon.blue"}
                                        borderColor={"neon.blue"}
                                        backdropFilter={"blur(12px)"}
                                        bg={"glass.bgHover"}
                                        transition="all 0.3s ease-in-out"
                                        _selected={{
                                            filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                                            boxShadow: "0 0 20px rgba(0,212,255,1)",
                                            color: "neon.blue",
                                            bg: "glass.bgHover",
                                            backdropFilter: "blur(12px)",
                                        }}
                                        _hover={{
                                            filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                            boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                            color: "neon.purple",
                                            bg: "glass.bgHover",
                                            backdropFilter: "blur(12px)",
                                            borderColor: "neon.purple"
                                        }}
                                        focusRing={"none"}>
                                <Text color="text.secondary">{page.value}</Text>
                            </IconButton>
                        )}
                        ellipsis={<IconButton color={"neon.blue"}
                                              borderColor={"neon.blue"}
                                              backdropFilter={"blur(12px)"}
                                              bg={"glass.bgHover"}
                                              transition="all 0.3s ease-in-out"
                                              _selected={{
                                                  filter: "drop-shadow(0 0 16px rgba(0,212,255,0.9))",
                                                  boxShadow: "0 0 20px rgba(0,212,255,1)",
                                                  color: "neon.blue",
                                                  bg: "glass.bgHover",
                                                  backdropFilter: "blur(12px)",
                                              }}
                                              _hover={{
                                                  filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                                  boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                                  color: "neon.purple",
                                                  bg: "glass.bgHover",
                                                  backdropFilter: "blur(12px)",
                                                  borderColor: "neon.purple"
                                              }}
                                              focusRing={"none"}>
                            <EllipsisIcon />
                        </IconButton>}
                    />

                    <Pagination.NextTrigger asChild>
                        <IconButton color={"neon.blue"}
                                    borderColor={"neon.blue"}
                                    backdropFilter={"blur(12px)"}
                                    bg={"glass.bgHover"}
                                    _hover={{
                                        filter: "drop-shadow(0 0 16px rgba(212, 0,255,0.9))",
                                        boxShadow: "0 0 20px rgba(212, 0,255,1)",
                                        color: "neon.purple",
                                        bg: "glass.bgHover",
                                        backdropFilter: "blur(12px)",
                                        borderColor: "neon.purple",
                                    }}
                                    transition="all 0.3s ease-in-out"
                                    focusRing={"none"}>
                            <ChevronRightIcon />
                        </IconButton>
                    </Pagination.NextTrigger>
                </ButtonGroup>
            </Pagination.Root>
        </Flex>
    )
}