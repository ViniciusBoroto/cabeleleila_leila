import { AlertTriangle, X } from "lucide-react";

interface CancelConfirmationModalProps {
  appointmentDate: string;
  loading: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

const CancelConfirmationModal = ({
  appointmentDate,
  loading,
  onConfirm,
  onClose,
}: CancelConfirmationModalProps) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString("pt-BR", {
      weekday: "long",
      day: "2-digit",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
        <div className="px-6 py-4 border-b flex justify-between items-center">
          <div className="flex items-center gap-3">
            <AlertTriangle className="w-6 h-6 text-red-600" />
            <h2 className="text-xl font-bold text-gray-900">
              Confirmar Cancelamento
            </h2>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
            disabled={loading}
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        <div className="p-6">
          <p className="text-gray-700 mb-4">
            Tem certeza que deseja cancelar o agendamento para:
          </p>
          <div className="bg-gray-50 rounded-lg p-4 mb-6">
            <p className="font-semibold text-gray-900">
              {formatDate(appointmentDate)}
            </p>
          </div>
          <p className="text-sm text-gray-600 mb-6">
            Esta ação não pode ser desfeita. Você precisará criar um novo
            agendamento caso mude de ideia.
          </p>

          <div className="flex gap-3">
            <button
              type="button"
              onClick={onClose}
              disabled={loading}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition disabled:opacity-50"
            >
              Voltar
            </button>
            <button
              type="button"
              onClick={onConfirm}
              disabled={loading}
              className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition disabled:bg-gray-400"
            >
              {loading ? "Cancelando..." : "Confirmar Cancelamento"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CancelConfirmationModal;
