export const dashboardPath = '/';

export const itemsListPath = '/items';
export const newItemPath = '/items/new';
export const itemEditRoutePath = '/items/:itemId/edit';

export const tagsListPath = '/tags';
export const newTagPath = '/tags/new';
export const tagEditRoutePath = '/tags/:tagId/edit';

export function buildEditItemPath(itemId: string) {
    return `/items/${itemId}/edit`;
}

export function buildEditTagPath(tagId: string) {
    return `/tags/${tagId}/edit`;
}
