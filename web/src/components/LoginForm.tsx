import { useState } from "react";

interface LoginFormProps {
  onSubmit: (email: string, senha: string, lembrarMe: boolean) => Promise<void>;
  loading: boolean;
  error: string;
}

const LoginForm = ({ onSubmit, loading, error }: LoginFormProps) => {
  const [email, setEmail] = useState(() => {
    // Load saved email from localStorage if available
    return localStorage.getItem("email") || "";
  });
  const [senha, setSenha] = useState("");
  const [lembrarMe, setLembrarMe] = useState(!!localStorage.getItem("email"));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(email, senha, lembrarMe);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
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
          disabled={loading}
          required
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
          disabled={loading}
          required
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
            disabled={loading}
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

      {/* Mensagem de erro */}
      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      {/* Botão de Entrar */}
      <button
        type="submit"
        disabled={loading || !email || !senha}
        className="w-full bg-indigo-600 text-white py-3 rounded-lg font-semibold hover:bg-indigo-700 focus:ring-4 focus:ring-indigo-200 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        {loading ? "Entrando..." : "Entrar"}
      </button>
    </form>
  );
};

export default LoginForm;
