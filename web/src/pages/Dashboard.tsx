import { useState, useEffect } from "react";
import { PlusIcon } from "@heroicons/react/24/outline";
import AppointmentForm from "../components/AppointmentForm";
import AppointmentList from "../components/AppointmentList";

// Tipos baseados na API
interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface User {
  id: number;
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

const API_BASE = "http://localhost:8080/api";

const formatDateToISO = (dateString: string) => {
  const date = new Date(dateString);
  return date.toISOString();
};

const SalonDashboard = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [availableServices, setAvailableServices] = useState<Service[]>([]);
  const [showNewAppointment, setShowNewAppointment] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchAppointments();
    fetchServices();
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

  const fetchServices = async () => {
    try {
      const response = await fetch(`${API_BASE}/services`, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        setAvailableServices(data || []);
      } else {
        console.error("Erro ao carregar serviços");
      }
    } catch (err) {
      console.error("Erro de conexão ao carregar serviços:", err);
    }
  };

  const handleCreateAppointment = async (
    date: string,
    selectedServices: number[]
  ) => {
    const token = localStorage.getItem("token");
    if (!token) return;

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
          date: formatDateToISO(date),
          services: services,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setAppointments([...appointments, data]);
        setShowNewAppointment(false);
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

  const handleCloseModal = () => {
    setShowNewAppointment(false);
    setError("");
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
            <PlusIcon className="w-5 h-5" />
            Novo Agendamento
          </button>
        </div>

        {/* Mensagem de erro */}
        {error && !showNewAppointment && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Lista de Agendamentos */}
        <AppointmentList
          appointments={appointments}
          loading={loading}
          onCreateNew={() => setShowNewAppointment(true)}
        />

        {/* Modal Novo Agendamento */}
        {showNewAppointment && (
          <AppointmentForm
            availableServices={availableServices}
            loading={loading}
            error={error}
            onSubmit={handleCreateAppointment}
            onClose={handleCloseModal}
          />
        )}
      </div>
    </div>
  );
};

export default SalonDashboard;
