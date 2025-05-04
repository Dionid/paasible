import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { usePaasibleApi } from "@/lib/paasible";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { mergeErrors } from "@/lib/rhf";
import { AlertCircle } from "lucide-react";

const formSchema = z.object({
  machineName: z.string().min(6),
});

export const SetLocalMachineIdWidget = ({
  className,
  machineName,
  onSuccess,
  ...props
}: {
  machineName: string;
  onSuccess: () => void;
} & React.ComponentPropsWithoutRef<"div">) => {
  const api = usePaasibleApi();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      machineName,
    },
  });

  const onFormSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await api.Mutations.createCurrentMachine(values.machineName);
      onSuccess();
    } catch (err) {
      form.setError("machineName", err as Error);
    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Set local machine name</CardTitle>
          <CardDescription>
            This will help you identify your current ansible playbooks machine
            from others.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={form.handleSubmit(onFormSubmit)}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Input
                  id="machineName"
                  type="text"
                  placeholder="some_unique_machine_name"
                  required
                  {...form.register("machineName")}
                />
              </div>
              {mergeErrors(form) != "" && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>{mergeErrors(form)}</AlertDescription>
                </Alert>
              )}
              <Button type="submit" className="w-full">
                Set machine name
              </Button>
            </div>
            <div className="mt-4 text-sm text-gray-500">
              If you don't know how to set this up and you see that this field
              is already filled, just click "Set machine name".
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
