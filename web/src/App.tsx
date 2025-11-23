import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage.tsx";
import Dashboard from "./pages/Agendar.tsx";
import PrivateRoute from "./components/PrivateRoute";
import NotFoundPage from "./pages/NotFound.tsx";
import ForbiddenPage from "./pages/Forbidden.tsx";
import AdminDashboardPage from "./pages/AdminDashboard.tsx";
import RegisterPage from "./pages/RegisterPage.tsx";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Rota raiz - redireciona para login ou dashboard */}
        <Route
          path="/"
          element={
            localStorage.getItem("token") ? (
              <Navigate to="/agendar" replace />
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />

        <Route path="/login" element={<LoginPage />} />

        <Route path="/registrar" element={<RegisterPage />} />

        {/* Rota de Dashboard - protegida, só acessa se estiver logado */}
        <Route
          path="/agendar"
          element={
            <PrivateRoute allowedRoles={["admin", "customer"]}>
              <Dashboard />
            </PrivateRoute>
          }
        />

        {/* Rota de Administração - protegida, só acessa se estiver logado e ser admin */}
        <Route
          path="/admin/dashboard"
          element={
            <PrivateRoute allowedRoles={["admin"]}>
              <AdminDashboardPage />
            </PrivateRoute>
          }
        />

        {/* Rota 404 - qualquer caminho não encontrado */}
        <Route path="*" element={<NotFoundPage />} />

        {/* Rota Forbidden - qualquer caminho não autorizado */}
        <Route path="/forbidden" element={<ForbiddenPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
