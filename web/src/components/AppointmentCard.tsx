import {
  CalendarIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon,
  ExclamationCircleIcon,
} from "@heroicons/react/24/outline";

interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface AppointmentCardProps {
  date: string;
  status: "PENDING" | "CONFIRMED" | "DONE" | "CANCELED";
  services: Service[];
}

const AppointmentCard = ({ date, status, services }: AppointmentCardProps) => {
  const getStatusBadge = (status: string) => {
    const badges = {
      PENDING: {
        icon: ExclamationCircleIcon,
        color: "bg-yellow-100 text-yellow-800",
        text: "Pendente",
      },
      CONFIRMED: {
        icon: CheckCircleIcon,
        color: "bg-blue-100 text-blue-800",
        text: "Confirmado",
      },
      DONE: {
        icon: CheckCircleIcon,
        color: "bg-green-100 text-green-800",
        text: "Concluído",
      },
      CANCELED: {
        icon: XCircleIcon,
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

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-md transition p-6">
      <div className="flex justify-between items-start mb-4">
        <div className="flex items-start gap-3">
          <CalendarIcon className="w-5 h-5 text-purple-600 mt-1" />
          <div>
            <h3 className="font-semibold text-lg text-gray-900">
              {formatDate(date)}
            </h3>
            <div className="flex items-center gap-2 mt-1 text-sm text-gray-600">
              <ClockIcon className="w-4 h-4" />
              <span>{calculateTotalDuration(services)} minutos</span>
            </div>
          </div>
        </div>
        {getStatusBadge(status)}
      </div>

      <div className="border-t pt-4">
        <h4 className="font-medium text-gray-700 mb-3">Serviços:</h4>
        <div className="space-y-2">
          {services.map((service) => (
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
            R$ {calculateTotal(services).toFixed(2)}
          </span>
        </div>
      </div>
    </div>
  );
};

export default AppointmentCard;
