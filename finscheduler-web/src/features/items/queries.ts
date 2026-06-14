import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query';
import ItemsService from '../../api/items.ts';
import type {ItemDetailedDto, ItemFilter, ItemModification} from '../../api/items.types.ts';

const itemsService = new ItemsService();

export const itemsQueryKeys = {
    all: ['items'] as const,
    lists: () => [...itemsQueryKeys.all, 'list'] as const,
    list: (filter: ItemFilter) => [...itemsQueryKeys.lists(), filter] as const,
    details: () => [...itemsQueryKeys.all, 'detail'] as const,
    detail: (itemId: string) => [...itemsQueryKeys.details(), itemId] as const,
};

export function useItemsListQuery(filter: ItemFilter) {
    return useQuery({
        queryKey: itemsQueryKeys.list(filter),
        queryFn: () => itemsService.getListingInfo(filter),
    });
}

export function useItemDetailsQuery(itemId?: string) {
    return useQuery<ItemDetailedDto | null>({
        queryKey: itemsQueryKeys.detail(itemId ?? ''),
        queryFn: () => itemsService.getDetailedInfo(itemId!),
        enabled: Boolean(itemId),
    });
}

export function useCreateItemMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (item: ItemModification) => itemsService.createItem(item),
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: itemsQueryKeys.lists()});
        },
    });
}

export function useUpdateItemMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (params: {itemId: string; item: ItemModification}) => {
            await itemsService.updateItem(params.itemId, params.item);
            return params.itemId;
        },
        onSuccess: async (itemId) => {
            await Promise.all([
                queryClient.invalidateQueries({queryKey: itemsQueryKeys.lists()}),
                queryClient.invalidateQueries({queryKey: itemsQueryKeys.detail(itemId)}),
            ]);
        },
    });
}

export function useDeleteItemsMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (itemIds: string[]) => {
            await Promise.all(itemIds.map((itemId) => itemsService.deleteItem(itemId)));
            return itemIds;
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: itemsQueryKeys.all});
        },
    });
}

export function useBulkCashbackMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (payload: {cashback: number; itemIds?: string[]; tagId?: string}) => {
            if (payload.itemIds && payload.itemIds.length > 0) {
                await itemsService.updateCashbackByItems(payload.itemIds, payload.cashback);
                return payload;
            }

            if (payload.tagId) {
                await itemsService.updateCashbackByTag(payload.tagId, payload.cashback);
                return payload;
            }

            throw new Error('Не удалось определить цель обновления кешбэка');
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: itemsQueryKeys.all});
        },
    });
}
