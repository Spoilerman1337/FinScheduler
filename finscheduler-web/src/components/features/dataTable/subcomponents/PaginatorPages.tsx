import {ButtonGroup, IconButton, Pagination, Text} from "@chakra-ui/react";
import {ChevronLeftIcon, ChevronRightIcon, EllipsisIcon} from "lucide-react";

interface PaginatorPagesProps {
    totalPages?: number;
    page?: number;
    onPageChange: (page: number) => void;
}

export default function PaginatorPages(props: PaginatorPagesProps){
    return( <Pagination.Root
        count={props.totalPages}
        pageSize={1}
        page={props.page}
        onPageChange={(e) => props.onPageChange(e.page)}
        key={props.totalPages}
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
                    <ChevronLeftIcon/>
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
                    <EllipsisIcon/>
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
                    <ChevronRightIcon/>
                </IconButton>
            </Pagination.NextTrigger>
        </ButtonGroup>
    </Pagination.Root>)
}