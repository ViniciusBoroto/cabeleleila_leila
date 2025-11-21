import { Navigate } from "react-router-dom";

interface PrivateRouteProps {
  children: React.ReactNode;
}

export default function PrivateRoute({ children }: PrivateRouteProps) {
  // Verifica se o usuário está logado (tem token)
  const token = localStorage.getItem("token");

  // Se não tiver token, redireciona para o login
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  // Se tiver token, mostra a página
  return <>{children}</>;
}
