import { PocketBaseProvider, usePocketBase } from "@/lib/paasible";
import { PaasibleApiProvider } from "@/lib/paasible";
import { LoginPage } from "./pages/auth/login";
import {
  BrowserRouter,
  Outlet,
  Route,
  Routes,
  useNavigate,
} from "react-router";
import { DashboardPage } from "./pages/app/dashboard";
import { AppLayout } from "./pages/app";

const AuthLayout = () => {
  const navigate = useNavigate();
  const pb = usePocketBase();

  if (pb.authStore.isValid) {
    navigate("/app");

    return;
  }

  return <Outlet />;
};

const AppInner = () => {
  const pb = usePocketBase();

  return (
    // TODO: use env variable for host
    <PaasibleApiProvider config={{ host: "http://127.0.0.1:8090", pb }}>
      {/* <AppAuth /> */}

      <Routes>
        {/* <Route index element={<Home />} /> */}

        <Route path="auth" element={<AuthLayout />}>
          <Route index path="login" element={<LoginPage />} />
        </Route>

        <Route path="app" element={<AppLayout />}>
          <Route index element={<DashboardPage />} />
        </Route>
      </Routes>
    </PaasibleApiProvider>
  );
};

const App = () => {
  return (
    <BrowserRouter>
      <PocketBaseProvider>
        <AppInner />
      </PocketBaseProvider>
    </BrowserRouter>
  );
};

export default App;
