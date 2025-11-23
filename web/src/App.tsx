import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage.tsx";
import Dashboard from "./pages/Dashboard";
import PrivateRoute from "./components/PrivateRoute";
import NotFoundPage from "./pages/NotFound.tsx";
import ForbiddenPage from "./pages/Forbidden.tsx";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Rota raiz - redireciona para login ou dashboard */}
        <Route
          path="/"
          element={
            localStorage.getItem("token") ? (
              <Navigate to="/dashboard" replace />
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />

        {/* Rota de Login - acessível sem autenticação */}
        <Route path="/login" element={<LoginPage />} />

        {/* Rota de Dashboard - protegida, só acessa se estiver logado */}
        <Route
          path="/dashboard"
          element={
            <PrivateRoute allowedRoles={["admin", "customer"]}>
              <Dashboard />
            </PrivateRoute>
          }
        />

        {/* Rota de Administração - protegida, só acessa se estiver logado e ser admin */}
        <Route
          path="/admin"
          element={
            <PrivateRoute allowedRoles={["admin"]}>
              <h1>Admin Page</h1>
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
