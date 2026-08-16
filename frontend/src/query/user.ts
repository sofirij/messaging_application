"use client"
import { getUser, searchUsername } from "@/lib/api/user"
import { queryOptions } from "@tanstack/react-query"

export const userQueryOptions = queryOptions({
    queryKey: ["user"],
    queryFn: async () => await getUser(),
    staleTime: Infinity,
})

export const userInvalidateQueryOptions = {
    queryKey: userQueryOptions.queryKey,
    exact: true,
    refetchType: "none"
} as const

export const userSearchQueryOptions = (query: string) => queryOptions({
    queryKey: ["userSearch", query],
    queryFn: () => searchUsername(query),
    staleTime: 15000,
})