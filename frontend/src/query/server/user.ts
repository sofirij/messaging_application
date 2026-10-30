import { getUser } from "@/lib/api/server/user"
import { queryOptions } from "@tanstack/react-query"

export const userQueryOptions = queryOptions({
    queryKey: ["user"],
    queryFn: async () => await getUser(),
    staleTime: Infinity,
})