import React, { useState, useEffect } from "react";

// Ícones SVG
const Calendar = ({ className = "w-5 h-5" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
    />
  </svg>
);

const Plus = ({ className = "w-5 h-5" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M12 4v16m8-8H4"
    />
  </svg>
);

const Clock = ({ className = "w-4 h-4" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const CheckCircle = ({ className = "w-3 h-3" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const XCircle = ({ className = "w-3 h-3" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const AlertCircle = ({ className = "w-3 h-3" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const X = ({ className = "w-6 h-6" }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M6 18L18 6M6 6l12 12"
    />
  </svg>
);

// Tipos baseados na API
interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface Customer {
  id: number;
  user_id: number;
  is_active: boolean;
}

interface Appointment {
  id: number;
  customer_id: number;
  date: string;
  status: "PENDING" | "CONFIRMED" | "DONE" | "CANCELED";
  services: Service[];
  customer?: Customer;
  created_at?: string;
  updated_at?: string;
}

const API_BASE = "http://localhost:8080/api";

const SalonDashboard = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [showNewAppointment, setShowNewAppointment] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Estados do formulário de novo agendamento
  const [newAppointmentDate, setNewAppointmentDate] = useState("");
  const [selectedServices, setSelectedServices] = useState<number[]>([]);

  // Lista de serviços disponíveis (você pode adaptar conforme sua necessidade)
  const availableServices: Service[] = [
    { id: 1, name: "Corte de Cabelo", price: 50.0, duration_minutes: 30 },
    { id: 2, name: "Coloração", price: 120.0, duration_minutes: 90 },
    { id: 3, name: "Escova", price: 40.0, duration_minutes: 45 },
    { id: 4, name: "Hidratação", price: 80.0, duration_minutes: 60 },
    { id: 5, name: "Manicure", price: 35.0, duration_minutes: 40 },
  ];

  useEffect(() => {
    fetchAppointments();
  }, []);

  const fetchAppointments = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE}/appointments`, {
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

  const handleCreateAppointment = async () => {
    const token = localStorage.getItem("token");
    if (!token || selectedServices.length === 0) return;

    setLoading(true);
    setError("");

    try {
      const services = availableServices.filter((s) =>
        selectedServices.includes(s.id)
      );

      const response = await fetch(`${API_BASE}/appointments`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          date: newAppointmentDate,
          services: services,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setAppointments([...appointments, data]);
        setShowNewAppointment(false);
        setNewAppointmentDate("");
        setSelectedServices([]);
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Erro ao criar agendamento");
      }
    } catch (err) {
      setError("Erro ao conectar com o servidor");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const toggleService = (serviceId: number) => {
    setSelectedServices((prev) =>
      prev.includes(serviceId)
        ? prev.filter((id) => id !== serviceId)
        : [...prev, serviceId]
    );
  };

  const getStatusBadge = (status: string) => {
    const badges = {
      PENDING: {
        icon: AlertCircle,
        color: "bg-yellow-100 text-yellow-800",
        text: "Pendente",
      },
      CONFIRMED: {
        icon: CheckCircle,
        color: "bg-blue-100 text-blue-800",
        text: "Confirmado",
      },
      DONE: {
        icon: CheckCircle,
        color: "bg-green-100 text-green-800",
        text: "Concluído",
      },
      CANCELED: {
        icon: XCircle,
        color: "bg-red-100 text-red-800",
        text: "Cancelado",
      },
    };

    const badge = badges[status as keyof typeof badges] || badges.PENDING;
    const Icon = badge.icon;

    return (
      <span
        className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium ${badge.color}`}
      >
        <Icon className="w-3 h-3" />
        {badge.text}
      </span>
    );
  };

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

  const calculateTotal = (services: Service[]) => {
    return services.reduce((sum, service) => sum + service.price, 0);
  };

  const calculateTotalDuration = (services: Service[]) => {
    return services.reduce((sum, service) => sum + service.duration_minutes, 0);
  };

  const handleCloseModal = () => {
    setShowNewAppointment(false);
    setError("");
    setSelectedServices([]);
    setNewAppointmentDate("");
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">
              Meus Agendamentos
            </h1>
            <p className="text-gray-600 mt-1">
              Gerencie seus horários no salão
            </p>
          </div>
          <button
            onClick={() => setShowNewAppointment(true)}
            className="flex items-center gap-2 bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition"
          >
            <Plus className="w-5 h-5" />
            Novo Agendamento
          </button>
        </div>

        {/* Mensagem de erro */}
        {error && !showNewAppointment && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Loading */}
        {loading && !showNewAppointment && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
            <p className="mt-4 text-gray-600">Carregando agendamentos...</p>
          </div>
        )}

        {/* Lista de Agendamentos */}
        {!loading && appointments.length === 0 && (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <Calendar className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-700 mb-2">
              Nenhum agendamento encontrado
            </h3>
            <p className="text-gray-500 mb-6">
              Comece criando seu primeiro agendamento
            </p>
            <button
              onClick={() => setShowNewAppointment(true)}
              className="bg-purple-600 text-white px-6 py-2 rounded-lg hover:bg-purple-700 transition"
            >
              Criar Agendamento
            </button>
          </div>
        )}

        {!loading && appointments.length > 0 && (
          <div className="grid gap-4">
            {appointments.map((appointment) => (
              <div
                key={appointment.id}
                className="bg-white rounded-lg shadow hover:shadow-md transition p-6"
              >
                <div className="flex justify-between items-start mb-4">
                  <div className="flex items-start gap-3">
                    <Calendar className="w-5 h-5 text-purple-600 mt-1" />
                    <div>
                      <h3 className="font-semibold text-lg text-gray-900">
                        {formatDate(appointment.date)}
                      </h3>
                      <div className="flex items-center gap-2 mt-1 text-sm text-gray-600">
                        <Clock className="w-4 h-4" />
                        <span>
                          {calculateTotalDuration(appointment.services)} minutos
                        </span>
                      </div>
                    </div>
                  </div>
                  {getStatusBadge(appointment.status)}
                </div>

                <div className="border-t pt-4">
                  <h4 className="font-medium text-gray-700 mb-3">Serviços:</h4>
                  <div className="space-y-2">
                    {appointment.services.map((service) => (
                      <div
                        key={service.id}
                        className="flex justify-between items-center bg-gray-50 rounded p-3"
                      >
                        <div>
                          <p className="font-medium text-gray-900">
                            {service.name}
                          </p>
                          <p className="text-sm text-gray-600">
                            {service.duration_minutes} min
                          </p>
                        </div>
                        <p className="font-semibold text-purple-600">
                          R$ {service.price.toFixed(2)}
                        </p>
                      </div>
                    ))}
                  </div>
                  <div className="mt-4 pt-4 border-t flex justify-between items-center">
                    <span className="font-semibold text-gray-700">Total:</span>
                    <span className="text-xl font-bold text-purple-600">
                      R$ {calculateTotal(appointment.services).toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Modal Novo Agendamento */}
        {showNewAppointment && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
              <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
                <h2 className="text-2xl font-bold text-gray-900">
                  Novo Agendamento
                </h2>
                <button
                  onClick={handleCloseModal}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              <div className="p-6">
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
                    {error}
                  </div>
                )}

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Data e Hora
                  </label>
                  <input
                    type="datetime-local"
                    value={newAppointmentDate}
                    onChange={(e) => setNewAppointmentDate(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent"
                  />
                </div>

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-3">
                    Selecione os Serviços
                  </label>
                  <div className="space-y-2">
                    {availableServices.map((service) => (
                      <label
                        key={service.id}
                        className={`flex items-center justify-between p-4 border-2 rounded-lg cursor-pointer transition ${
                          selectedServices.includes(service.id)
                            ? "border-purple-600 bg-purple-50"
                            : "border-gray-200 hover:border-gray-300"
                        }`}
                      >
                        <div className="flex items-center gap-3">
                          <input
                            type="checkbox"
                            checked={selectedServices.includes(service.id)}
                            onChange={() => toggleService(service.id)}
                            className="w-5 h-5 text-purple-600 rounded focus:ring-purple-600"
                          />
                          <div>
                            <p className="font-medium text-gray-900">
                              {service.name}
                            </p>
                            <p className="text-sm text-gray-600">
                              {service.duration_minutes} minutos
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

                {selectedServices.length > 0 && (
                  <div className="bg-purple-50 rounded-lg p-4 mb-6">
                    <div className="flex justify-between items-center mb-2">
                      <span className="font-medium text-gray-700">
                        Total de serviços:
                      </span>
                      <span className="font-semibold text-gray-900">
                        {selectedServices.length}
                      </span>
                    </div>
                    <div className="flex justify-between items-center mb-2">
                      <span className="font-medium text-gray-700">
                        Duração total:
                      </span>
                      <span className="font-semibold text-gray-900">
                        {calculateTotalDuration(
                          availableServices.filter((s) =>
                            selectedServices.includes(s.id)
                          )
                        )}{" "}
                        min
                      </span>
                    </div>
                    <div className="flex justify-between items-center pt-2 border-t border-purple-200">
                      <span className="font-bold text-gray-900">
                        Valor total:
                      </span>
                      <span className="text-xl font-bold text-purple-600">
                        R${" "}
                        {calculateTotal(
                          availableServices.filter((s) =>
                            selectedServices.includes(s.id)
                          )
                        ).toFixed(2)}
                      </span>
                    </div>
                  </div>
                )}

                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={handleCloseModal}
                    className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
                  >
                    Cancelar
                  </button>
                  <button
                    type="button"
                    onClick={handleCreateAppointment}
                    disabled={loading || selectedServices.length === 0}
                    className="flex-1 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
                  >
                    {loading ? "Criando..." : "Confirmar Agendamento"}
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

export default SalonDashboard;
