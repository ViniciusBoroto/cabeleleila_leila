import { useState } from "react";
import RegisterForm from "../components/RegisterForm";

export default function RegisterPage() {
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState("");
  const [sucesso, setSucesso] = useState(false);

  const handleRegister = async (
    nome: string,
    email: string,
    senha: string,
    confirmarSenha: string
  ) => {
    // Limpa mensagens anteriores
    setErro("");
    setSucesso(false);

    // Validações básicas
    if (senha !== confirmarSenha) {
      setErro("As senhas não coincidem");
      return;
    }

    if (senha.length < 6) {
      setErro("A senha deve ter no mínimo 6 caracteres");
      return;
    }

    setCarregando(true);

    try {
      // Faz a requisição para seu backend
      const response = await fetch("http://localhost:8080/api/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: nome,
          email: email,
          password: senha,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Registro bem-sucedido!
        console.log("Registro realizado com sucesso!");
        setSucesso(true);

        // Aguarda 2 segundos e redireciona para o login
        setTimeout(() => {
          window.location.href = "/login";
        }, 2000);
      } else {
        // Erro no registro (email já existe, etc)
        setErro(data.message || "Erro ao criar conta. Tente novamente.");
      }
    } catch (error) {
      // Erro de conexão ou outro erro
      console.error("Erro ao fazer registro:", error);
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
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Criar Conta</h1>
          <p className="text-gray-600">Preencha os dados para começar</p>
        </div>

        {/* Mensagem de Sucesso */}
        {sucesso && (
          <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
            <p className="text-green-800 text-sm text-center">
              Conta criada com sucesso! Redirecionando para o login...
            </p>
          </div>
        )}

        {/* Formulário de Registro */}
        <RegisterForm
          onSubmit={handleRegister}
          loading={carregando}
          error={erro}
        />

        {/* Link para fazer login */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            Já tem uma conta?{" "}
            <a
              href="/login"
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Fazer login
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
