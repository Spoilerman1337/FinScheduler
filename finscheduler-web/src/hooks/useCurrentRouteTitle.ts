import {type UIMatch, useMatches} from 'react-router-dom';

interface AppRouteHandle {
    title?: string;
}

export function useCurrentRouteTitle() {
    const matches = useMatches() as UIMatch<unknown, AppRouteHandle>[];

    for (let index = matches.length - 1; index >= 0; index -= 1) {
        const title = matches[index].handle?.title;

        if (title) {
            return title;
        }
    }

    return undefined;
}
