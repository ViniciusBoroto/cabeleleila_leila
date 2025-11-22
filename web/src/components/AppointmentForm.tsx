import { useState } from "react";
import { XMarkIcon } from "@heroicons/react/24/outline";

interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface AppointmentFormProps {
  availableServices: Service[];
  loading: boolean;
  error: string;
  onSubmit: (date: string, selectedServices: number[]) => Promise<void>;
  onClose: () => void;
}

const AppointmentForm = ({
  availableServices,
  loading,
  error,
  onSubmit,
  onClose,
}: AppointmentFormProps) => {
  const [newAppointmentDate, setNewAppointmentDate] = useState("");
  const [selectedServices, setSelectedServices] = useState<number[]>([]);
  const [validationError, setValidationError] = useState("");

  const toggleService = (serviceId: number) => {
    setSelectedServices((prev) =>
      prev.includes(serviceId)
        ? prev.filter((id) => id !== serviceId)
        : [...prev, serviceId]
    );
  };

  const calculateTotal = (services: Service[]) => {
    return services.reduce((sum, service) => sum + service.price, 0);
  };

  const calculateTotalDuration = (services: Service[]) => {
    return services.reduce((sum, service) => sum + service.duration_minutes, 0);
  };

  const getMinDateTime = () => {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    return now.toISOString().slice(0, 16);
  };

  const handleSubmit = async () => {
    setValidationError("");

    // Validate at least one service is selected
    if (selectedServices.length === 0) {
      setValidationError("Selecione pelo menos um serviço");
      return;
    }

    // Validate date is not in the past
    const selectedDate = new Date(newAppointmentDate);
    const now = new Date();
    if (selectedDate < now) {
      setValidationError("A data não pode ser anterior a agora");
      return;
    }

    // Validate date is provided
    if (!newAppointmentDate) {
      setValidationError("Selecione uma data e hora");
      return;
    }

    await onSubmit(newAppointmentDate, selectedServices);
    handleClose();
  };

  const handleClose = () => {
    setNewAppointmentDate("");
    setSelectedServices([]);
    setValidationError("");
    onClose();
  };

  const selectedServicesList = availableServices.filter((s) =>
    selectedServices.includes(s.id)
  );

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          <h2 className="text-2xl font-bold text-gray-900">Novo Agendamento</h2>
          <button
            onClick={handleClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        <div className="p-6">
          {(error || validationError) && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
              {error || validationError}
            </div>
          )}

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Data e Hora *
            </label>
            <input
              type="datetime-local"
              value={newAppointmentDate}
              onChange={(e) => setNewAppointmentDate(e.target.value)}
              min={getMinDateTime()}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent"
            />
            <p className="text-xs text-gray-500 mt-1">
              A data deve ser a partir de agora
            </p>
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-3">
              Selecione os Serviços *
            </label>
            {availableServices.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <p>Carregando serviços disponíveis...</p>
              </div>
            ) : (
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
            )}
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
                  {calculateTotalDuration(selectedServicesList)} min
                </span>
              </div>
              <div className="flex justify-between items-center pt-2 border-t border-purple-200">
                <span className="font-bold text-gray-900">Valor total:</span>
                <span className="text-xl font-bold text-purple-600">
                  R$ {calculateTotal(selectedServicesList).toFixed(2)}
                </span>
              </div>
            </div>
          )}

          <div className="flex gap-3">
            <button
              type="button"
              onClick={handleClose}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
            >
              Cancelar
            </button>
            <button
              type="button"
              onClick={handleSubmit}
              disabled={loading}
              className="flex-1 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {loading ? "Criando..." : "Confirmar Agendamento"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AppointmentForm;
