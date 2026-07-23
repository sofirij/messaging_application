import { getUser, searchUsername } from "@/lib/api/user"
import { queryOptions } from "@tanstack/react-query"

export const userQueryOptions = queryOptions({
    queryKey: ["user"],
    queryFn: async () => await getUser(),
    staleTime: Infinity,
    retry: false
})

export const userRemoveQueryOptions = {
    queryKey: userQueryOptions.queryKey,
    exact: true
}

export const userSearchQueryOptions = (query: string) => queryOptions({
    queryKey: ["userSearch", query],
    queryFn: () => searchUsername(query),
    staleTime: 15000
})