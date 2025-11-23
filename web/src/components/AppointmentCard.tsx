import {
  Calendar,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Pencil,
  Trash2,
} from "lucide-react";
import {
  canEditAppointment,
  canCancelAppointment,
} from "../utils/appointmentHelpers";

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
  created_at?: string;
  updated_at?: string;
}

interface AppointmentCardProps {
  appointment: Appointment;
  onEdit: (appointment: Appointment) => void;
  onCancel: (appointmentId: number) => void;
}

const AppointmentCard = ({
  appointment,
  onEdit,
  onCancel,
}: AppointmentCardProps) => {
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

  const editCheck = canEditAppointment(appointment.date);
  const cancelCheck = canCancelAppointment(appointment.status);

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-md transition p-6">
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
        <div className="flex items-center gap-2">
          {getStatusBadge(appointment.status)}

          {/* Edit Button - sempre visível mas pode estar desabilitado */}
          <div className="relative group">
            <button
              onClick={() => editCheck.canEdit && onEdit(appointment)}
              disabled={!editCheck.canEdit}
              className={`p-2 rounded-lg transition ${
                editCheck.canEdit
                  ? "text-purple-600 hover:bg-purple-50 cursor-pointer"
                  : "text-gray-400 cursor-not-allowed"
              }`}
              title={
                editCheck.canEdit ? "Editar agendamento" : editCheck.reason
              }
            >
              <Pencil className="w-4 h-4" />
            </button>
            {/* Tooltip quando desabilitado */}
            {!editCheck.canEdit && (
              <div className="absolute hidden group-hover:block right-0 top-full mt-2 w-64 bg-gray-900 text-white text-xs rounded-lg p-3 z-10 shadow-lg">
                <div className="absolute -top-1 right-4 w-2 h-2 bg-gray-900 transform rotate-45"></div>
                {editCheck.reason}
              </div>
            )}
          </div>

          {/* Cancel Button - sempre visível mas pode estar desabilitado */}
          <div className="relative group">
            <button
              onClick={() => cancelCheck.canCancel && onCancel(appointment.id)}
              disabled={!cancelCheck.canCancel}
              className={`p-2 rounded-lg transition ${
                cancelCheck.canCancel
                  ? "text-red-600 hover:bg-red-50 cursor-pointer"
                  : "text-gray-400 cursor-not-allowed"
              }`}
              title={
                cancelCheck.canCancel
                  ? "Cancelar agendamento"
                  : cancelCheck.reason
              }
            >
              <Trash2 className="w-4 h-4" />
            </button>
            {/* Tooltip quando desabilitado */}
            {!cancelCheck.canCancel && (
              <div className="absolute hidden group-hover:block right-0 top-full mt-2 w-64 bg-gray-900 text-white text-xs rounded-lg p-3 z-10 shadow-lg">
                <div className="absolute -top-1 right-4 w-2 h-2 bg-gray-900 transform rotate-45"></div>
                {cancelCheck.reason}
              </div>
            )}
          </div>
        </div>
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
                <p className="font-medium text-gray-900">{service.name}</p>
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
  );
};

export default AppointmentCard;
