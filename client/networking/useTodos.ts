import { pb } from "@/types/pb"
import useSWR from "swr";
import { useEffect, useState } from "react";

function todoFetcher(url: string) {
    return fetch(url).then(res => res.json()).then(data => data.data as pb.Todo[]);
}

export function useTodos() {
    const { data, error, isLoading } = useSWR("/api/todos", todoFetcher);

    return {
        todos: data,
        error,
        isLoading,
    }
}

export function useFetchTodos() {
    const [isLoading, setIsLoading] = useState(true);
    const [todos, setTodos] = useState<pb.Todo[] | undefined>();
    const [error, setError]  = useState<any>(null);

    useEffect(() => {
        setIsLoading(true);
        const doFetch = async () => {
            try {
                const response = await fetch("/api/todos");
                const data = await response.json();
                setTodos(data.data);
            } catch (err) {
                setError(err);
            }
            setIsLoading(false);
        }

        doFetch();
    }, [])

    return {
        todos,
        setTodos,
        isLoading,
        error
    }
}