import { useState } from "react";
import LoginForm from "../components/LoginForm";

export default function LoginPage() {
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState("");

  const handleLogin = async (
    email: string,
    senha: string,
    lembrarMe: boolean
  ) => {
    // Limpa mensagens de erro anteriores
    setErro("");
    setCarregando(true);

    try {
      // Faz a requisição para seu backend
      const response = await fetch("http://localhost:8080/api/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          password: senha,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Login bem-sucedido!
        console.log("Login realizado com sucesso!");
        console.log("Dados do usuário:", data.user);

        // Salva o token JWT e os dados do usuário no localStorage
        localStorage.setItem("token", data.token);
        localStorage.setItem("user", JSON.stringify(data.user));

        // Se o usuário marcou "Lembrar-me", salva o email também
        if (lembrarMe) {
          localStorage.setItem("email", email);
        }

        // Redireciona para o dashboard
        window.location.href = "/agendar";
      } else {
        // Erro no login (credenciais inválidas, etc)
        setErro(data.message || "Email ou senha incorretos");
      }
    } catch (error) {
      // Erro de conexão ou outro erro
      console.error("Erro ao fazer login:", error);
      setErro("Erro ao conectar com o servidor. Verifique sua conexão.");
    } finally {
      setCarregando(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md p-8">
        {/* Cabeçalho */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Bem-vindo de volta
          </h1>
          <p className="text-gray-600">Entre na sua conta</p>
        </div>

        {/* Formulário de Login */}
        <LoginForm onSubmit={handleLogin} loading={carregando} error={erro} />

        {/* Link para criar conta */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            Não tem uma conta?{" "}
            <a
              href="/registrar"
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Criar conta
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
