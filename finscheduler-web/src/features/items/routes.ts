export const itemsListPath = '/items';
export const newItemPath = '/items/new';

export function buildEditItemPath(itemId: string) {
    return `/items/${itemId}/edit`;
}
