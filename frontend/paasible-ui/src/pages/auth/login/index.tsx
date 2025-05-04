
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useDiq } from "@/lib/diq"
import { usePaasibleApi } from "@/lib/paasible"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle } from "lucide-react"
import { z } from "zod"
import { useForm, UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Form } from "@/components/ui/form"
import { mergeErrors } from "@/lib/rhf"

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8)
})

type LoginFormProps = React.ComponentProps<"div"> & {
  form: UseFormReturn<z.infer<typeof formSchema>>,
  onFormSubmit: (values: z.infer<typeof formSchema>) => void
}

export function LoginForm({
  className,
  onFormSubmit,
  form,
  ...props
}: LoginFormProps) {
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="overflow-hidden py-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onFormSubmit)} className="p-6 md:p-8">
              <div className="flex flex-col gap-6">
                <div className="flex flex-col items-center text-center">
                  <h1 className="text-2xl font-bold">Sign In</h1>
                  <p className="text-balance text-muted-foreground">
                    Your admin must give you login and password
                  </p>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="m@example.com"
                    required
                    {...form.register("email")}
                  />
                </div>
                <div className="grid gap-2">
                  <div className="flex items-center">
                    <Label htmlFor="password">Password</Label>
                  </div>
                  <Input
                    id="password"
                    type="password"
                    required
                    {...form.register("password")}
                    />
                </div>
                {
                    mergeErrors(form) != "" && (
                      <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4" />
                        <AlertTitle>Error</AlertTitle>
                        <AlertDescription>
                          { mergeErrors(form) }
                        </AlertDescription>
                      </Alert>
                    )
                  }
                <Button type="submit" className="w-full">
                  Login
                </Button>
              </div>
            </form>
          </Form>
          <div className="relative hidden bg-muted md:block">
            <img
              src="https://ui.shadcn.com/placeholder.svg"
              alt="Image"
              className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
            />
          </div>
        </CardContent>
      </Card>
      <div className="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 hover:[&_a]:text-primary">
        By clicking continue, you agree to our <a href="#">Terms of Service</a>{" "}
        and <a href="#">Privacy Policy</a>.
      </div>
    </div>
  )
}

export const LoginPage = () => {
  const api = usePaasibleApi()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  })

  const signInMutation = useDiq(
    api.Mutations.signIn
  )
 
  const onFormSubmit = async (values: z.infer<typeof formSchema>) => {
    await signInMutation.request(values)
    // TODO: move to main page
    // ...
  }

  return (
      <div className="flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10">
          <div className="w-full max-w-sm md:max-w-3xl">
              <LoginForm
                form={ form }
                onFormSubmit={ onFormSubmit }
              />
          </div>
      </div>
  )
}