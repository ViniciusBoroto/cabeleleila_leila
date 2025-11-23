import { Calendar } from "lucide-react";
import AppointmentCard from "./AppointmentCard";

interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface Appointment {
  id: number;
  user_id: number;
  date: string;
  status: "PENDING" | "CONFIRMED" | "DONE" | "CANCELED";
  services: Service[];
  user?: { id: number };
  created_at?: string;
  updated_at?: string;
}

interface AppointmentListProps {
  appointments: Appointment[];
  loading: boolean;
  onCreateNew: () => void;
  onEdit: (appointment: Appointment) => void;
  onCancel: (appointmentId: number) => void;
}

const AppointmentList = ({
  appointments,
  loading,
  onCreateNew,
  onEdit,
  onCancel,
}: AppointmentListProps) => {
  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
        <p className="mt-4 text-gray-600">Carregando agendamentos...</p>
      </div>
    );
  }

  if (appointments.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-12 text-center">
        <Calendar className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <h3 className="text-xl font-semibold text-gray-700 mb-2">
          Nenhum agendamento encontrado
        </h3>
        <p className="text-gray-500 mb-6">
          Comece criando seu primeiro agendamento
        </p>
        <button
          onClick={onCreateNew}
          className="bg-purple-600 text-white px-6 py-2 rounded-lg hover:bg-purple-700 transition"
        >
          Criar Agendamento
        </button>
      </div>
    );
  }

  return (
    <div className="grid gap-4">
      {appointments.map((appointment) => (
        <AppointmentCard
          key={appointment.id}
          appointment={appointment}
          onEdit={onEdit}
          onCancel={onCancel}
        />
      ))}
    </div>
  );
};

export default AppointmentList;
