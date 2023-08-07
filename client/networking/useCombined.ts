import { pb } from "@/types/pb"
import { useEffect, useState } from "react";
import useSWR from "swr";

function combinedFetcher(url: string) {
    return fetch(url).then(res => res.json()).then(data => data.data as pb.Combined[]);
}

export function useCombined() {
    const { data, error, isLoading } = useSWR("/api/combined", combinedFetcher);

    return {
        combined: data,
        error,
        isLoading,
    }
}

export function useFetchCombined() {
    const [isLoading, setIsLoading] = useState(true);
    const [combined, setCombined] = useState<pb.Combined[] | undefined>();
    const [error, setError]  = useState<any>(null);

    useEffect(() => {
        setIsLoading(true);
        const doFetch = async () => {
            try {
                const response = await fetch("/api/combined");
                const data = await response.json();
                setCombined(data.data);
            } catch (err) {
                setError(err);
            }
            setIsLoading(false);
        }

        doFetch();
    }, [])

    return {
        combined: combined,
        setCombined: setCombined,
        isLoading,
        error
    }
}
