import {Theme} from "@chakra-ui/react";
import {render} from "@testing-library/react";
import type {ReactElement} from "react";
import {MemoryRouter} from "react-router-dom";
import {Provider} from "../components/ui/provider.tsx";

interface RenderWithProvidersOptions {
    route?: string;
    initialEntries?: string[];
}

export function renderWithProviders(
    ui: ReactElement,
    options: RenderWithProvidersOptions = {},
) {
    const initialEntries = options.initialEntries ?? [options.route ?? "/"];

    return render(
        <Provider>
            <Theme>
                <MemoryRouter initialEntries={initialEntries}>
                    {ui}
                </MemoryRouter>
            </Theme>
        </Provider>,
    );
}
