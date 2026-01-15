import { apiClient } from './client';
import type { ItemDto, PaginatedList, ItemFilter } from './types';

export const itemsApi = {
    getItems: async (filter?: ItemFilter): Promise<PaginatedList<ItemDto>> => {
        return apiClient.getItems(filter);
    },

    getItemById: async (id: string): Promise<ItemDto> => {
        return apiClient.getItemById(id);
    },

    createItem: async (item: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>): Promise<string> => {
        return apiClient.createItem(item);
    },

    updateItem: async (id: string, item: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>): Promise<void> => {
        return apiClient.updateItem(id, item);
    },

    deleteItem: async (id: string): Promise<void> => {
        return apiClient.deleteItem(id);
    },
};

