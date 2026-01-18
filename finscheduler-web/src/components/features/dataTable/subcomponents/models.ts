import {createListCollection} from "@chakra-ui/react";

export const pageSizes = createListCollection({
    items: [
        {label: "10", value: 10},
        {label: "25", value: 25},
        {label: "50", value: 50},
        {label: "100", value: 100},
    ],
})