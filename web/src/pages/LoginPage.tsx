import { useState } from "react";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [senha, setSenha] = useState("");
  const [lembrarMe, setLembrarMe] = useState(false);
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState("");

  const handleLogin = async () => {
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
        window.location.href = "/dashboard";
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

        {/* Mensagem de erro */}
        {erro && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-sm text-red-600">{erro}</p>
          </div>
        )}

        <div className="space-y-6">
          {/* Campo de Email */}
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Endereço de Email
            </label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition"
              placeholder="seu@email.com"
              disabled={carregando}
            />
          </div>

          {/* Campo de Senha */}
          <div>
            <label
              htmlFor="senha"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Senha
            </label>
            <input
              type="password"
              id="senha"
              value={senha}
              onChange={(e) => setSenha(e.target.value)}
              className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition"
              placeholder="••••••••"
              disabled={carregando}
            />
          </div>

          {/* Lembrar-me e Esqueci a senha */}
          <div className="flex items-center justify-between">
            <label className="flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={lembrarMe}
                onChange={(e) => setLembrarMe(e.target.checked)}
                className="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
                disabled={carregando}
              />
              <span className="ml-2 text-sm text-gray-600">Lembrar-me</span>
            </label>
            <a
              href="#"
              className="text-sm text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Esqueceu a senha?
            </a>
          </div>

          {/* Botão de Entrar */}
          <button
            onClick={handleLogin}
            disabled={carregando || !email || !senha}
            className="w-full bg-indigo-600 text-white py-3 rounded-lg font-semibold hover:bg-indigo-700 focus:ring-4 focus:ring-indigo-200 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {carregando ? "Entrando..." : "Entrar"}
          </button>
        </div>

        {/* Link para criar conta */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            Não tem uma conta?{" "}
            <a
              href="#"
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
