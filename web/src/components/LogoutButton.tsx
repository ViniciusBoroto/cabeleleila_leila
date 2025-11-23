import { LogOut } from "lucide-react";

export default function LogoutButton() {
  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    localStorage.removeItem("name");
    window.location.href = "/login";
  };

  return (
    <button
      onClick={handleLogout}
      className="flex max-h-10 items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors cursor-pointer"
    >
      <LogOut size={20} />
      <span>Sair</span>
    </button>
  );
}
