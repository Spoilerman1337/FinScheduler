import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query';
import {useCallback} from 'react';
import type {TagDetailedDto, TagFilter, TagLookupFilter, TagModification} from '../../api/tags.types.ts';
import TagsService from '../../api/tags.ts';
import {mapLookupsToSelectOptions} from '../shared.ts';
import {itemsQueryKeys} from '../items/queries.ts';

const tagsService = new TagsService();

export const tagsQueryKeys = {
    all: ['tags'] as const,
    lists: () => [...tagsQueryKeys.all, 'list'] as const,
    list: (filter: TagFilter) => [...tagsQueryKeys.lists(), filter] as const,
    details: () => [...tagsQueryKeys.all, 'detail'] as const,
    detail: (tagId: string) => [...tagsQueryKeys.details(), tagId] as const,
    lookups: () => [...tagsQueryKeys.all, 'lookup'] as const,
    lookup: (filter: TagLookupFilter) => [...tagsQueryKeys.lookups(), filter] as const,
};

export function useTagsListQuery(filter: TagFilter) {
    return useQuery({
        queryKey: tagsQueryKeys.list(filter),
        queryFn: () => tagsService.getListingInfo(filter),
    });
}

export function useTagDetailsQuery(tagId?: string) {
    return useQuery<TagDetailedDto | null>({
        queryKey: tagsQueryKeys.detail(tagId ?? ''),
        queryFn: () => tagsService.getDetailedInfo(tagId!),
        enabled: Boolean(tagId),
    });
}

export function useCreateTagMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (tag: TagModification) => tagsService.createTag(tag),
        onSuccess: async () => {
            await Promise.all([
                queryClient.invalidateQueries({queryKey: tagsQueryKeys.lists()}),
                queryClient.invalidateQueries({queryKey: tagsQueryKeys.lookups()}),
            ]);
        },
    });
}

export function useUpdateTagMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (params: {tagId: string; tag: TagModification}) => {
            await tagsService.updateTag(params.tagId, params.tag);
            return params.tagId;
        },
        onSuccess: async (tagId) => {
            await Promise.all([
                queryClient.invalidateQueries({queryKey: tagsQueryKeys.lists()}),
                queryClient.invalidateQueries({queryKey: tagsQueryKeys.detail(tagId)}),
                queryClient.invalidateQueries({queryKey: tagsQueryKeys.lookups()}),
                queryClient.invalidateQueries({queryKey: itemsQueryKeys.all}),
            ]);
        },
    });
}

export function useTagLookupOptionsLoader(pageSize = 20) {
    const queryClient = useQueryClient();

    return useCallback(
        async ({page, search}: {page: number; search: string}) => {
            const filter: TagLookupFilter = {
                page,
                pageSize,
                name: search || undefined,
            };

            const result = await queryClient.fetchQuery({
                queryKey: tagsQueryKeys.lookup(filter),
                queryFn: () => tagsService.getLookup(filter),
                staleTime: 5 * 60_000,
            });

            return {
                options: mapLookupsToSelectOptions(result.data),
                hasMore: (page + 1) * pageSize < result.count,
            };
        },
        [pageSize, queryClient],
    );
}
