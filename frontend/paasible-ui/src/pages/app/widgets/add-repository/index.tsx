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

export const AddRepositoryWidget = ({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"div">) => {
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Add repository</CardTitle>
          <CardDescription>
            You can do this by adding git repository URL or adding repository to
            the dedicated folder.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Input
                  id="email"
                  type="email"
                  placeholder="github.com/.../..."
                  required
                />
              </div>
              <Button type="submit" className="w-full">
                Add
              </Button>
            </div>
            <div className="mt-4 text-sm text-gray-500">
              This is needed because first project must be created from
              repository. Even if you don't want to use git, just create
              repository from your project and move it to dedicated folder.
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
