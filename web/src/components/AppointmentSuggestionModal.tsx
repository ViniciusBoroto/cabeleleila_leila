import {
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
} from "@heroicons/react/24/outline";

interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface Appointment {
  id: number;
  date: string;
  services: Service[];
}

interface AppointmentSuggestionModalProps {
  existingAppointment: Appointment;
  newServices: Service[];
  loading: boolean;
  onMerge: () => Promise<void>;
  onReject: () => Promise<void>;
  onClose: () => void;
}

const AppointmentSuggestionModal = ({
  existingAppointment,
  newServices,
  loading,
  onMerge,
  onReject,
  onClose,
}: AppointmentSuggestionModalProps) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString("pt-BR", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
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

  const handleMerge = async () => {
    await onMerge();
  };

  const handleReject = async () => {
    await onReject();
  };

  const allServices = [...existingAppointment.services, ...newServices];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          <div className="flex items-center gap-3">
            <CheckCircleIcon className="w-6 h-6 text-blue-600" />
            <h2 className="text-2xl font-bold text-gray-900">
              Agendamento Existente Detectado
            </h2>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        <div className="p-6">
          <div className="mb-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-blue-800">
              Você já possui um agendamento na mesma semana. Deseja mesclar os
              serviços no agendamento existente ou criar um novo?
            </p>
          </div>

          <div className="mb-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Agendamento Existente:
            </h3>
            <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
              <p className="text-sm text-gray-600 mb-2">
                <span className="font-medium">Data:</span>{" "}
                {formatDate(existingAppointment.date)}
              </p>
              <p className="text-sm text-gray-600 mb-3">
                <span className="font-medium">Serviços:</span>
              </p>
              <div className="space-y-2">
                {existingAppointment.services.map((service) => (
                  <div
                    key={service.id}
                    className="flex justify-between items-center bg-white p-2 rounded border border-gray-200"
                  >
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {service.name}
                      </p>
                      <p className="text-xs text-gray-500">
                        {service.duration_minutes} min
                      </p>
                    </div>
                    <span className="text-sm font-semibold text-gray-900">
                      R$ {service.price.toFixed(2)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="mb-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Novos Serviços a Adicionar:
            </h3>
            <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
              <div className="space-y-2">
                {newServices.map((service) => (
                  <div
                    key={service.id}
                    className="flex justify-between items-center bg-white p-2 rounded border border-purple-200"
                  >
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {service.name}
                      </p>
                      <p className="text-xs text-gray-500">
                        {service.duration_minutes} min
                      </p>
                    </div>
                    <span className="text-sm font-semibold text-purple-600">
                      R$ {service.price.toFixed(2)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="mb-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Total após Mesclagem:
            </h3>
            <div className="bg-gradient-to-r from-purple-50 to-blue-50 rounded-lg p-4 border border-purple-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-medium text-gray-700">
                  Total de serviços:
                </span>
                <span className="font-semibold text-gray-900">
                  {allServices.length}
                </span>
              </div>
              <div className="flex justify-between items-center mb-2">
                <span className="font-medium text-gray-700">
                  Duração total:
                </span>
                <span className="font-semibold text-gray-900">
                  {calculateTotalDuration(allServices)} min
                </span>
              </div>
              <div className="flex justify-between items-center pt-2 border-t border-purple-200">
                <span className="font-bold text-gray-900">Valor total:</span>
                <span className="text-xl font-bold text-purple-600">
                  R$ {calculateTotal(allServices).toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex gap-3">
            <button
              type="button"
              onClick={handleReject}
              disabled={loading}
              className="flex-1 px-4 py-3 border-2 border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition disabled:bg-gray-100 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <XCircleIcon className="w-5 h-5" />
              Criar Novo
            </button>
            <button
              type="button"
              onClick={handleMerge}
              disabled={loading}
              className="flex-1 px-4 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <CheckCircleIcon className="w-5 h-5" />
              {loading ? "Processando..." : "Mesclar Serviços"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AppointmentSuggestionModal;
