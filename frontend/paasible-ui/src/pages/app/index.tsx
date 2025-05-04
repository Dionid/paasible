import { Separator } from "@/components/ui/separator";
import {
  SidebarProvider,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarGroup,
  SidebarGroupContent,
  SidebarContent,
  SidebarGroupLabel,
  SidebarRail,
  SidebarInset,
  SidebarTrigger,
  Sidebar,
} from "@/components/ui/sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
} from "@/components/ui/breadcrumb";
import {
  AudioWaveform,
  BellIcon,
  ClipboardListIcon,
  Command,
  FileCodeIcon,
  FoldersIcon,
  GalleryVerticalEnd,
  HelpCircleIcon,
  KeyRoundIcon,
  LayoutDashboardIcon,
  SettingsIcon,
  TimerIcon,
} from "lucide-react";
import { NavUser } from "@/components/ui/nav-user";
import { NavSecondary } from "@/components/ui/nav-secondary";
import { TeamSwitcher } from "@/components/ui/team-switcher";
import { Outlet, useNavigate } from "react-router";
import { usePaasibleApi } from "@/lib/paasible";
import { Button } from "@/components/ui/button";
import { useDiq } from "@/lib/diq";
import { useEffect, useMemo } from "react";
import { AddRepositoryWidget } from "./widgets/add-repository";
import { SetLocalMachineIdWidget } from "./widgets/set-local-machine-id";

const navigationData = {
  navMain: [
    {
      title: "Main",
      url: "#",
      icon: SettingsIcon,
      items: [
        {
          title: "Dashboard",
          url: "#",
          icon: LayoutDashboardIcon,
          isActive: true,
        },
        {
          title: "Runs",
          url: "#",
          icon: ClipboardListIcon,
        },
        {
          title: "Cron",
          url: "#",
          icon: TimerIcon,
        },
      ],
    },
    {
      title: "Local",
      url: "#",
      icon: SettingsIcon,
      items: [
        {
          title: "Files",
          url: "#",
          icon: FoldersIcon,
        },
        {
          title: "Playbooks",
          url: "#",
          icon: FileCodeIcon,
        },
        {
          title: "Inventories",
          url: "#",
          icon: KeyRoundIcon,
        },
      ],
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      url: "#",
      icon: SettingsIcon,
    },
    {
      title: "Get Help",
      url: "#",
      icon: HelpCircleIcon,
    },
  ],
  teams: [
    {
      name: "Acme Inc",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Acme Corp.",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Evil Corp.",
      logo: Command,
      plan: "Free",
    },
  ],
};

export const AppLayout = () => {
  const navigate = useNavigate();
  const api = usePaasibleApi();
  const { pb } = api;

  const getRepositoriesQ = useDiq(api.Queries.getRepositories);
  const checkLocalMachineIdQ = useDiq(api.Queries.checkLocalMachineId);

  useEffect(() => {
    getRepositoriesQ.request();
    checkLocalMachineIdQ.request();
  }, []);

  useEffect(() => {
    if (!pb.authStore.isValid) {
      navigate("/auth");

      return;
    }
  }, [pb.authStore.isValid, navigate]);

  const onMachineSetSuccess = useMemo(() => {
    return () => {
      checkLocalMachineIdQ.request();
    };
  }, [checkLocalMachineIdQ]);

  if (getRepositoriesQ.isPending || checkLocalMachineIdQ.isPending) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (getRepositoriesQ.error || checkLocalMachineIdQ.error) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-red-500">
          {getRepositoriesQ.error?.message ||
            checkLocalMachineIdQ.error?.message}
        </div>
      </div>
    );
  }

  if (
    checkLocalMachineIdQ.data &&
    checkLocalMachineIdQ.data.data.status != "found"
  ) {
    return (
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
          <SetLocalMachineIdWidget
            machineName={checkLocalMachineIdQ.data.data.value}
            onSuccess={onMachineSetSuccess}
          />
        </div>
      </div>
    );
  }

  if (getRepositoriesQ.data && getRepositoriesQ.data.length === 0) {
    return (
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
          <AddRepositoryWidget />
        </div>
      </div>
    );
  }

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>
          <TeamSwitcher teams={navigationData.teams} />
        </SidebarHeader>
        <SidebarContent>
          {navigationData.navMain.map((item) => (
            <SidebarGroup key={item.title}>
              <SidebarGroupLabel>{item.title}</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {item.items.map((subItem) => (
                    <SidebarMenuItem key={subItem.title}>
                      <SidebarMenuButton asChild isActive={subItem.isActive}>
                        <a href={subItem.url}>
                          <subItem.icon />
                          {subItem.title}
                        </a>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          ))}
          <NavSecondary
            items={navigationData.navSecondary}
            className="mt-auto"
          />
        </SidebarContent>
        <SidebarRail />
      </Sidebar>
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink href="#">Dashboard</BreadcrumbLink>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div className="ml-auto flex items-center gap-2">
            <Button size="lg" variant="ghost">
              <BellIcon />
            </Button>
            <NavUser
              user={{
                name: pb.authStore.record!["name"],
                email: pb.authStore.record!["email"],
                avatar: "https://ui.shadcn.com/avatars/shadcn.jpg",
              }}
            />
          </div>
        </header>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  );
};
