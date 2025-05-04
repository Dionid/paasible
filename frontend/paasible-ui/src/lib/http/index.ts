export type ExpandedResponse<T extends Record<PropertyKey, unknown>> = {
  status: number;
  response: Response;
  data: T;
};

export const expandedFetch = async <T extends Record<PropertyKey, unknown>>(
  request: Promise<Response>
): Promise<ExpandedResponse<T>> => {
  const response = await request;

  return {
    status: response.status,
    response,
    data: await response.json(),
  };
};

export type EitherResponse<T extends Record<PropertyKey, unknown>> =
  | {
      status: "success";
      code: number;
      data: T;
    }
  | {
      status: "error";
      code: number;
      message: string;
    };

export const eitherFetch = async <T extends Record<PropertyKey, unknown>>(
  request: Promise<Response>
): Promise<EitherResponse<T>> => {
  try {
    const response = await request;

    if (response.ok) {
      return {
        status: "success",
        code: response.status,
        data: (await response.json()) as T,
      };
    } else {
      return {
        status: "error",
        code: response.status,
        message: response.statusText,
      };
    }
  } catch (error) {
    return {
      status: "error",
      code: 500,
      message: error instanceof Error ? error.message : "Unknown error",
    };
  }
};
