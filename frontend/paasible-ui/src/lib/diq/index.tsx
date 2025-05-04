import { useMemo, useState } from "react";

export type DiqResult<P extends unknown[], R> = {
  isPending: boolean;
  error: Error | null;
  data: R | undefined;
  request: (...args: P) => Promise<R | Error>;
};

export function useDiq<P extends unknown[], R>(
  query: (...args: P) => Promise<R>,
  opts: {
    ignoreError?: boolean;
  } = {}
): DiqResult<P, R> {
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [data, setData] = useState<R>();

  const request = useMemo(() => {
    return async (...args: P) => {
      setIsPending(true);
      setError(null);
      try {
        const result = await query(...args);
        setData(result);
        return result;
      } catch (err) {
        if (opts.ignoreError) {
          return err as Error;
        }

        setError(err as Error);
        return err as Error;
      } finally {
        setIsPending(false);
      }
    };
  }, [query, opts]);

  const result = useMemo(() => {
    return {
      isPending,
      error,
      data,
      request,
    };
  }, [data, error, request, isPending]);

  return result;
}
