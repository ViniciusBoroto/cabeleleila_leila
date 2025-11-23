import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface RegisterFormProps {
  onSubmit: (
    nome: string,
    email: string,
    senha: string,
    confirmarSenha: string
  ) => void;
  loading: boolean;
  error: string;
}

export default function RegisterForm({
  onSubmit,
  loading,
  error,
}: RegisterFormProps) {
  const [nome, setNome] = useState("");
  const [email, setEmail] = useState("");
  const [senha, setSenha] = useState("");
  const [confirmarSenha, setConfirmarSenha] = useState("");
  const [mostrarSenha, setMostrarSenha] = useState(false);
  const [mostrarConfirmarSenha, setMostrarConfirmarSenha] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(nome, email, senha, confirmarSenha);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {/* Mensagem de Erro */}
      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800 text-sm">{error}</p>
        </div>
      )}

      {/* Campo Nome */}
      <div>
        <label
          htmlFor="nome"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          Nome Completo
        </label>
        <input
          id="nome"
          type="text"
          value={nome}
          onChange={(e) => setNome(e.target.value)}
          required
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
          placeholder="Seu nome completo"
          disabled={loading}
        />
      </div>

      {/* Campo Email */}
      <div>
        <label
          htmlFor="email"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
          placeholder="seu@email.com"
          disabled={loading}
        />
      </div>

      {/* Campo Senha */}
      <div>
        <label
          htmlFor="senha"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          Senha
        </label>
        <div className="relative">
          <input
            id="senha"
            type={mostrarSenha ? "text" : "password"}
            value={senha}
            onChange={(e) => setSenha(e.target.value)}
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition pr-12"
            placeholder="Mínimo 6 caracteres"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() => setMostrarSenha(!mostrarSenha)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
            disabled={loading}
          >
            {mostrarSenha ? (
              <EyeOff className="w-5 h-5" />
            ) : (
              <Eye className="w-5 h-5" />
            )}
          </button>
        </div>
      </div>

      {/* Campo Confirmar Senha */}
      <div>
        <label
          htmlFor="confirmarSenha"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          Confirmar Senha
        </label>
        <div className="relative">
          <input
            id="confirmarSenha"
            type={mostrarConfirmarSenha ? "text" : "password"}
            value={confirmarSenha}
            onChange={(e) => setConfirmarSenha(e.target.value)}
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition pr-12"
            placeholder="Digite a senha novamente"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() => setMostrarConfirmarSenha(!mostrarConfirmarSenha)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
            disabled={loading}
          >
            {mostrarConfirmarSenha ? (
              <EyeOff className="w-5 h-5" />
            ) : (
              <Eye className="w-5 h-5" />
            )}
          </button>
        </div>
      </div>

      {/* Botão de Submit */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-4 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? "Criando conta..." : "Criar Conta"}
      </button>
    </form>
  );
}
