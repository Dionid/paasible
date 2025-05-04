import { UseFormReturn } from "react-hook-form";

export const mergeErrors = <T extends Record<PropertyKey, unknown>>(
  form: UseFormReturn<T>
): string => {
  let merged = "";

  for (const [key, value] of Object.entries(form.formState.errors)) {
    if (value && value.message) {
      merged += `${key}: ${value.message}\n`;
    }
  }

  return merged;
};
