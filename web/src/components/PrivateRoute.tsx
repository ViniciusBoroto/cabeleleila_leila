import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";

interface PrivateRouteProps {
  children: React.ReactNode;
  allowedRoles?: string[];
}

const API_BASE = "http://localhost:8080/api";

export default function PrivateRoute({
  children,
  allowedRoles,
}: PrivateRouteProps) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [userRole, setUserRole] = useState<string | null>(null);

  const validateToken = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      setIsLoading(false);
      return;
    }

    const response = await fetch(`${API_BASE}/auth/validate`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    if (response.ok) {
      const data = await response.json();
      setUserRole(data.role);
      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(data.user));
      localStorage.setItem("role", data.role);
      setIsAuthenticated(true);
    } else {
      console.error("Erro ao validar token");
      localStorage.removeItem("token");
      setIsAuthenticated(false);
    }
    setIsLoading(false);
  };

  useEffect(() => {
    validateToken();
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated || !userRole) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && !allowedRoles.includes(userRole)) {
    return <Navigate to="/forbidden" replace />;
  }

  return <>{children}</>;
}
