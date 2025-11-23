import { useState, useEffect, useMemo } from "react";
import { Calendar, TrendingUp, DollarSign, Users, Clock } from "lucide-react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  BarChart,
  Bar,
} from "recharts";
import LogoutButton from "../components/LogoutButton";
import AppointmentFilter from "../components/AppointmentFilter";
import type { FilterPeriod } from "../components/AppointmentFilter";
import {
  filterAppointmentsByPeriod,
  getDefaultCustomDates,
} from "../utils/filterHelpers";

// Tipos
interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface User {
  id: number;
  name?: string;
  email?: string;
}

interface Appointment {
  id: number;
  user_id: number;
  date: string;
  status: "PENDING" | "CONFIRMED" | "DONE" | "CANCELED";
  services: Service[];
  user?: User;
  created_at?: string;
  updated_at?: string;
}

interface WeeklyStats {
  week: string;
  appointments: number;
  revenue: number;
  services: number;
}

const API_BASE = "http://localhost:8080/api";

const AdminDashboardPage = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [availableServices, setAvailableServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [selectedAppointment, setSelectedAppointment] =
    useState<Appointment | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [filter, setFilter] = useState<
    "all" | "PENDING" | "CONFIRMED" | "DONE" | "CANCELED"
  >("all");

  // Estados do filtro de período
  const [filterPeriod, setFilterPeriod] = useState<FilterPeriod>("all");
  const [customStartDate, setCustomStartDate] = useState("");
  const [customEndDate, setCustomEndDate] = useState("");

  // Estados do formulário
  const [editDate, setEditDate] = useState("");
  const [editServices, setEditServices] = useState<number[]>([]);
  const [editStatus, setEditStatus] =
    useState<Appointment["status"]>("PENDING");

  useEffect(() => {
    fetchAppointments();
    fetchServices();

    // Inicializar datas personalizadas
    const defaultDates = getDefaultCustomDates();
    setCustomStartDate(defaultDates.start);
    setCustomEndDate(defaultDates.end);
  }, []);

  // Aplicar filtro de período nos agendamentos
  const periodFilteredAppointments = useMemo(() => {
    return filterAppointmentsByPeriod(
      appointments,
      filterPeriod,
      customStartDate,
      customEndDate
    );
  }, [appointments, filterPeriod, customStartDate, customEndDate]);

  const fetchAppointments = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE}/admin/appointments`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        setAppointments(data || []);
      } else {
        setError("Erro ao carregar agendamentos");
      }
    } catch (err) {
      setError("Erro de conexão com o servidor");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchServices = async () => {
    try {
      const response = await fetch(`${API_BASE}/services`);
      if (response.ok) {
        const data = await response.json();
        setAvailableServices(data || []);
      }
    } catch (err) {
      console.error("Erro ao carregar serviços:", err);
    }
  };

  const handleEdit = (appointment: Appointment) => {
    setSelectedAppointment(appointment);
    const date = new Date(appointment.date);
    const localDate = new Date(
      date.getTime() - date.getTimezoneOffset() * 60000
    );
    setEditDate(localDate.toISOString().slice(0, 16));
    setEditServices(appointment.services.map((s) => s.id));
    setEditStatus(appointment.status);
    setShowEditModal(true);
  };

  const handleUpdateAppointment = async () => {
    if (!selectedAppointment) return;

    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const services = availableServices.filter((s) =>
        editServices.includes(s.id)
      );

      const response = await fetch(
        `${API_BASE}/admin/appointments/${selectedAppointment.id}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            date: new Date(editDate).toISOString(),
            services,
            status: editStatus,
          }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        setAppointments(
          appointments.map((ap) => (ap.id === data.id ? data : ap))
        );
        setShowEditModal(false);
        setSelectedAppointment(null);
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Erro ao atualizar agendamento");
      }
    } catch (err) {
      setError("Erro ao conectar com o servidor");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handlePeriodChange = (period: FilterPeriod) => {
    setFilterPeriod(period);
  };

  const handleCustomDateChange = (start: string, end: string) => {
    setCustomStartDate(start);
    setCustomEndDate(end);
  };

  // Calcular estatísticas semanais
  const calculateWeeklyStats = (): WeeklyStats[] => {
    const stats: { [key: string]: WeeklyStats } = {};

    appointments
      .filter((ap) => ap.status !== "CANCELED")
      .forEach((ap) => {
        const date = new Date(ap.date);
        const weekStart = new Date(date);
        weekStart.setDate(date.getDate() - date.getDay());
        const weekKey = weekStart.toISOString().split("T")[0];

        if (!stats[weekKey]) {
          stats[weekKey] = {
            week: new Date(weekKey).toLocaleDateString("pt-BR", {
              day: "2-digit",
              month: "short",
            }),
            appointments: 0,
            revenue: 0,
            services: 0,
          };
        }

        stats[weekKey].appointments += 1;
        stats[weekKey].revenue += ap.services.reduce(
          (sum, s) => sum + s.price,
          0
        );
        stats[weekKey].services += ap.services.length;
      });

    return Object.values(stats)
      .sort((a, b) => new Date(a.week).getTime() - new Date(b.week).getTime())
      .slice(-8);
  };

  const weeklyStats = calculateWeeklyStats();

  // Estatísticas gerais
  const totalRevenue = appointments
    .filter((ap) => ap.status !== "CANCELED")
    .reduce(
      (sum, ap) => sum + ap.services.reduce((s, srv) => s + srv.price, 0),
      0
    );

  const totalAppointments = appointments.filter(
    (ap) => ap.status !== "CANCELED"
  ).length;
  const pendingAppointments = appointments.filter(
    (ap) => ap.status === "PENDING"
  ).length;
  const avgDuration =
    appointments.length > 0
      ? appointments.reduce(
          (sum, ap) =>
            sum + ap.services.reduce((s, srv) => s + srv.duration_minutes, 0),
          0
        ) / appointments.length
      : 0;

  // Aplicar filtro de status na lista já filtrada por período
  const filteredAppointments =
    filter === "all"
      ? periodFilteredAppointments
      : periodFilteredAppointments.filter((ap) => ap.status === filter);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getStatusBadge = (status: string) => {
    const badges = {
      PENDING: { color: "bg-yellow-100 text-yellow-800", text: "Pendente" },
      CONFIRMED: { color: "bg-blue-100 text-blue-800", text: "Confirmado" },
      DONE: { color: "bg-green-100 text-green-800", text: "Concluído" },
      CANCELED: { color: "bg-red-100 text-red-800", text: "Cancelado" },
    };
    const badge = badges[status as keyof typeof badges] || badges.PENDING;
    return (
      <span
        className={`px-3 py-1 rounded-full text-xs font-medium ${badge.color}`}
      >
        {badge.text}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">
              Dashboard Admin
            </h1>
            <p className="text-gray-600">
              Visão completa dos agendamentos e performance
            </p>
          </div>
          <LogoutButton />
        </div>

        {/* Cards de Estatísticas */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-xl shadow-sm p-6 border-l-4 border-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Receita Total</p>
                <p className="text-2xl font-bold text-gray-900">
                  R$ {totalRevenue.toFixed(2)}
                </p>
              </div>
              <div className="bg-purple-100 p-3 rounded-lg">
                <DollarSign className="w-6 h-6 text-purple-600" />
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 border-l-4 border-blue-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">
                  Total de Agendamentos
                </p>
                <p className="text-2xl font-bold text-gray-900">
                  {totalAppointments}
                </p>
              </div>
              <div className="bg-blue-100 p-3 rounded-lg">
                <Calendar className="w-6 h-6 text-blue-600" />
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 border-l-4 border-yellow-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Pendentes</p>
                <p className="text-2xl font-bold text-gray-900">
                  {pendingAppointments}
                </p>
              </div>
              <div className="bg-yellow-100 p-3 rounded-lg">
                <Users className="w-6 h-6 text-yellow-600" />
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 border-l-4 border-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Duração Média</p>
                <p className="text-2xl font-bold text-gray-900">
                  {avgDuration.toFixed(0)} min
                </p>
              </div>
              <div className="bg-green-100 p-3 rounded-lg">
                <Clock className="w-6 h-6 text-green-600" />
              </div>
            </div>
          </div>
        </div>

        {/* Gráficos */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <div className="bg-white rounded-xl shadow-sm p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-purple-600" />
              Receita Semanal
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={weeklyStats}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="week" />
                <YAxis />
                <Tooltip
                  formatter={(value: number) => `R$ ${value.toFixed(2)}`}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="revenue"
                  stroke="#9333ea"
                  strokeWidth={2}
                  name="Receita"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <Calendar className="w-5 h-5 text-blue-600" />
              Agendamentos por Semana
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={weeklyStats}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="week" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar
                  dataKey="appointments"
                  fill="#3b82f6"
                  name="Agendamentos"
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Filtro de Período */}
        <AppointmentFilter
          selectedPeriod={filterPeriod}
          customStartDate={customStartDate}
          customEndDate={customEndDate}
          totalCount={appointments.length}
          filteredCount={periodFilteredAppointments.length}
          onPeriodChange={handlePeriodChange}
          onCustomDateChange={handleCustomDateChange}
        />

        {/* Filtros */}
        <div className="bg-white rounded-xl shadow-sm p-4 mb-6">
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => setFilter("all")}
              className={`px-4 py-2 rounded-lg transition ${
                filter === "all"
                  ? "bg-purple-600 text-white"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
              }`}
            >
              Todos ({periodFilteredAppointments.length})
            </button>
            <button
              onClick={() => setFilter("PENDING")}
              className={`px-4 py-2 rounded-lg transition ${
                filter === "PENDING"
                  ? "bg-yellow-600 text-white"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
              }`}
            >
              Pendentes (
              {
                periodFilteredAppointments.filter((a) => a.status === "PENDING")
                  .length
              }
              )
            </button>
            <button
              onClick={() => setFilter("CONFIRMED")}
              className={`px-4 py-2 rounded-lg transition ${
                filter === "CONFIRMED"
                  ? "bg-blue-600 text-white"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
              }`}
            >
              Confirmados (
              {
                periodFilteredAppointments.filter(
                  (a) => a.status === "CONFIRMED"
                ).length
              }
              )
            </button>
            <button
              onClick={() => setFilter("DONE")}
              className={`px-4 py-2 rounded-lg transition ${
                filter === "DONE"
                  ? "bg-green-600 text-white"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
              }`}
            >
              Concluídos (
              {
                periodFilteredAppointments.filter((a) => a.status === "DONE")
                  .length
              }
              )
            </button>
            <button
              onClick={() => setFilter("CANCELED")}
              className={`px-4 py-2 rounded-lg transition ${
                filter === "CANCELED"
                  ? "bg-red-600 text-white"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
              }`}
            >
              Cancelados (
              {
                periodFilteredAppointments.filter(
                  (a) => a.status === "CANCELED"
                ).length
              }
              )
            </button>
          </div>
        </div>

        {/* Lista de Agendamentos */}
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Cliente
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Data/Hora
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Serviços
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Valor
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Ações
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {filteredAppointments.length === 0 ? (
                  <tr>
                    <td
                      colSpan={6}
                      className="px-6 py-8 text-center text-gray-500"
                    >
                      Nenhum agendamento encontrado
                    </td>
                  </tr>
                ) : (
                  filteredAppointments.map((appointment) => (
                    <tr key={appointment.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">
                          {appointment.user?.name ||
                            `Usuário #${appointment.user_id}`}
                        </div>
                        <div className="text-sm text-gray-500">
                          {appointment.user?.email || ""}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {formatDate(appointment.date)}
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900">
                          {appointment.services.map((s) => s.name).join(", ")}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {getStatusBadge(appointment.status)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                        R${" "}
                        {appointment.services
                          .reduce((sum, s) => sum + s.price, 0)
                          .toFixed(2)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <button
                          onClick={() => handleEdit(appointment)}
                          className="text-purple-600 hover:text-purple-900 font-medium"
                        >
                          Editar
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Modal de Edição */}
        {showEditModal && selectedAppointment && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
              <div className="sticky top-0 bg-white border-b px-6 py-4">
                <h2 className="text-2xl font-bold text-gray-900">
                  Editar Agendamento
                </h2>
              </div>

              <div className="p-6">
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
                    {error}
                  </div>
                )}

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Cliente
                  </label>
                  <input
                    type="text"
                    value={
                      selectedAppointment.user?.name ||
                      `Usuário #${selectedAppointment.user_id}`
                    }
                    disabled
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100"
                  />
                </div>

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Data e Hora
                  </label>
                  <input
                    type="datetime-local"
                    value={editDate}
                    onChange={(e) => setEditDate(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600"
                  />
                </div>

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Status
                  </label>
                  <select
                    value={editStatus}
                    onChange={(e) =>
                      setEditStatus(e.target.value as Appointment["status"])
                    }
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600"
                  >
                    <option value="PENDING">Pendente</option>
                    <option value="CONFIRMED">Confirmado</option>
                    <option value="DONE">Concluído</option>
                    <option value="CANCELED">Cancelado</option>
                  </select>
                </div>

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-3">
                    Serviços
                  </label>
                  <div className="space-y-2 max-h-60 overflow-y-auto">
                    {availableServices.map((service) => (
                      <label
                        key={service.id}
                        className={`flex items-center justify-between p-3 border-2 rounded-lg cursor-pointer transition ${
                          editServices.includes(service.id)
                            ? "border-purple-600 bg-purple-50"
                            : "border-gray-200 hover:border-gray-300"
                        }`}
                      >
                        <div className="flex items-center gap-3">
                          <input
                            type="checkbox"
                            checked={editServices.includes(service.id)}
                            onChange={() => {
                              setEditServices((prev) =>
                                prev.includes(service.id)
                                  ? prev.filter((id) => id !== service.id)
                                  : [...prev, service.id]
                              );
                            }}
                            className="w-4 h-4 text-purple-600 rounded"
                          />
                          <div>
                            <p className="font-medium text-gray-900">
                              {service.name}
                            </p>
                            <p className="text-sm text-gray-600">
                              {service.duration_minutes} min
                            </p>
                          </div>
                        </div>
                        <span className="font-semibold text-purple-600">
                          R$ {service.price.toFixed(2)}
                        </span>
                      </label>
                    ))}
                  </div>
                </div>

                <div className="flex gap-3">
                  <button
                    onClick={() => {
                      setShowEditModal(false);
                      setSelectedAppointment(null);
                      setError("");
                    }}
                    className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
                  >
                    Cancelar
                  </button>
                  <button
                    onClick={handleUpdateAppointment}
                    disabled={loading || editServices.length === 0}
                    className="flex-1 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:bg-gray-400"
                  >
                    {loading ? "Salvando..." : "Salvar Alterações"}
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default AdminDashboardPage;
