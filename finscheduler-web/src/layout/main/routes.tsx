import Dashboard from "../../features/dashboard/Dashboard.tsx";
import Items from "../../features/items/Items.tsx";
import Tags from "../../features/tags/Tags.tsx";

export const routes = [
    {
        path: "/",
        title: "Дашборды",
        element: <Dashboard/>,
    },
    {
        path: "/items",
        title: "Каталог",
        element: <Items/>,
    },
    {
        path: "/tags",
        title: "Теги",
        element: <Tags/>,
    },
];
