import { pb } from "@/types/pb"
import { useEffect, useState } from "react";
import useSWR from "swr";

function userFetcher(url: string) {
    return fetch(url).then(res => res.json()).then(data => data.data as pb.User[]);
}

export function useUsers() {
    const { data, error, isLoading } = useSWR("/api/users", userFetcher);

    return {
        users: data,
        error,
        isLoading,
    }
}

export function useFetchUsers() {
    const [isLoading, setIsLoading] = useState(true);
    const [users, setUsers] = useState<pb.User[] | undefined>();
    const [error, setError]  = useState<any>(null);

    useEffect(() => {
        setIsLoading(true);
        const doFetch = async () => {
            try {
                const response = await fetch("/api/users");
                const data = await response.json();
                setUsers(data.data);
            } catch (err) {
                setError(err);
            }
            setIsLoading(false);
        }

        doFetch();
    }, [])

    return {
        users,
        setUsers,
        isLoading,
        error
    }
}
