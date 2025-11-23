import { useState, useEffect } from "react";
import { Plus } from "lucide-react";
import AppointmentForm from "../components/AppointmentForm";
import AppointmentList from "../components/AppointmentList";
import AppointmentSuggestionModal from "../components/AppointmentSuggestionModal";
import CancelConfirmationModal from "../components/CancelConfirmationModal";

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
  const [editingAppointment, setEditingAppointment] = useState<
    Appointment | undefined
  >();
  const [cancelingAppointment, setCancelingAppointment] = useState<
    Appointment | undefined
  >();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [suggestion, setSuggestion] = useState<{
    appointment: Appointment;
    newServices: Service[];
  } | null>(null);

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
        let data = await response.json();
        data = data.filter(
          (ap: Appointment): boolean => ap.status !== "CANCELED"
        );
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

  const handleSubmitAppointment = async (
    date: string,
    selectedServices: number[],
    appointmentId?: number
  ) => {
    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const services = availableServices.filter((s) =>
        selectedServices.includes(s.id)
      );

      const url = appointmentId
        ? `${API_BASE}/appointments/${appointmentId}`
        : `${API_BASE}/appointments`;

      const method = appointmentId ? "PUT" : "POST";

      const response = await fetch(url, {
        method,
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

        if (appointmentId) {
          // Update existing appointment
          setAppointments(
            appointments.map((ap) => (ap.id === appointmentId ? data : ap))
          );
          setShowNewAppointment(false);
          setEditingAppointment(undefined);
        } else {
          // Creating new appointment - check for suggestion
          if (data.suggestion) {
            setSuggestion({
              appointment: data.suggestion,
              newServices: services,
            });
            setShowNewAppointment(false);
          } else {
            // No suggestion, appointment was created
            setAppointments([...appointments, data.appointment]);
            setShowNewAppointment(false);
          }
        }
      } else {
        const errorData = await response.json();
        setError(
          errorData.error ||
            `Erro ao ${appointmentId ? "atualizar" : "criar"} agendamento`
        );
      }
    } catch (err) {
      setError("Erro ao conectar com o servidor");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelAppointment = async () => {
    if (!cancelingAppointment) return;

    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `${API_BASE}/appointments/${cancelingAppointment.id}/cancel`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        // Update the appointment status in the list
        setAppointments(appointments.filter((ap) => ap.id !== data.id));
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Erro ao cancelar agendamento");
      }
    } catch (err) {
      setError("Erro ao conectar com o servidor");
      console.error(err);
    } finally {
      setLoading(false);

      setCancelingAppointment(undefined);
    }
  };

  const handleEdit = (appointment: Appointment) => {
    setEditingAppointment(appointment);
    setShowNewAppointment(true);
  };

  const handleCancelClick = (appointmentId: number) => {
    const appointment = appointments.find((ap) => ap.id === appointmentId);
    if (appointment) {
      setCancelingAppointment(appointment);
    }
  };

  const handleCloseModal = () => {
    setShowNewAppointment(false);
    setEditingAppointment(undefined);
    setError("");
  };

  const handleCloseCancelModal = () => {
    setCancelingAppointment(undefined);
  };

  const handleMergeSuggestion = async () => {
    if (!suggestion) return;

    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `${API_BASE}/appointments/${suggestion.appointment.id}/merge`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            services: suggestion.newServices,
          }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        // Update the appointment in the list
        setAppointments(
          appointments.map((ap) => (ap.id === data.id ? data : ap))
        );
        setSuggestion(null);
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Erro ao mesclar agendamentos");
      }
    } catch (err) {
      setError("Erro ao conectar com o servidor");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleRejectSuggestion = async () => {
    if (!suggestion) return;

    const token = localStorage.getItem("token");
    if (!token) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE}/appointments`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          date: formatDateToISO(suggestion.appointment.date),
          services: suggestion.newServices,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setAppointments([...appointments, data.appointment || data]);
        setSuggestion(null);
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

  const handleCloseSuggestion = () => {
    setSuggestion(null);
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
        {error && !showNewAppointment && !cancelingAppointment && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Lista de Agendamentos */}
        <AppointmentList
          appointments={appointments}
          loading={loading}
          onCreateNew={() => setShowNewAppointment(true)}
          onEdit={handleEdit}
          onCancel={handleCancelClick}
        />

        {/* Modal Novo/Editar Agendamento */}
        {showNewAppointment && (
          <AppointmentForm
            availableServices={availableServices}
            loading={loading}
            error={error}
            onSubmit={handleSubmitAppointment}
            onClose={handleCloseModal}
            editingAppointment={editingAppointment}
          />
        )}

        {/* Modal Confirmação de Cancelamento */}
        {cancelingAppointment && (
          <CancelConfirmationModal
            appointmentDate={cancelingAppointment.date}
            loading={loading}
            onConfirm={handleCancelAppointment}
            onClose={handleCloseCancelModal}
          />
        )}

        {/* Modal Sugestão de Mesclagem */}
        {suggestion && (
          <AppointmentSuggestionModal
            existingAppointment={suggestion.appointment}
            newServices={suggestion.newServices}
            loading={loading}
            onMerge={handleMergeSuggestion}
            onReject={handleRejectSuggestion}
            onClose={handleCloseSuggestion}
          />
        )}
      </div>
    </div>
  );
};

export default SalonDashboard;
