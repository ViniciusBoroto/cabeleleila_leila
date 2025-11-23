import { Calendar, Filter, X } from "lucide-react";

export type FilterPeriod = "all" | "today" | "week" | "month" | "custom";

interface AppointmentFilterProps {
  selectedPeriod: FilterPeriod;
  customStartDate: string;
  customEndDate: string;
  totalCount: number;
  filteredCount: number;
  onPeriodChange: (period: FilterPeriod) => void;
  onCustomDateChange: (start: string, end: string) => void;
}

const AppointmentFilter = ({
  selectedPeriod,
  customStartDate,
  customEndDate,
  totalCount,
  filteredCount,
  onPeriodChange,
  onCustomDateChange,
}: AppointmentFilterProps) => {
  const periods = [
    { value: "all" as FilterPeriod, label: "Todos" },
    { value: "today" as FilterPeriod, label: "Hoje" },
    { value: "week" as FilterPeriod, label: "Esta Semana" },
    { value: "month" as FilterPeriod, label: "Este Mês" },
    { value: "custom" as FilterPeriod, label: "Período Personalizado" },
  ];

  return (
    <div className="bg-white rounded-lg shadow p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-purple-600" />
          <h3 className="text-lg font-semibold text-gray-900">
            Filtrar Agendamentos
          </h3>
        </div>
        {selectedPeriod !== "all" && (
          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-600">
              Mostrando {filteredCount} de {totalCount}
            </span>
            <button
              onClick={() => onPeriodChange("all")}
              className="text-gray-500 hover:text-gray-700 transition"
              title="Limpar filtro"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        )}
      </div>

      {/* Botões de Período */}
      <div className="flex flex-wrap gap-2 mb-4">
        {periods.map((period) => (
          <button
            key={period.value}
            onClick={() => onPeriodChange(period.value)}
            className={`px-4 py-2 rounded-lg font-medium transition ${
              selectedPeriod === period.value
                ? "bg-purple-600 text-white"
                : "bg-gray-100 text-gray-700 hover:bg-gray-200"
            }`}
          >
            {period.label}
          </button>
        ))}
      </div>

      {/* Datas Personalizadas */}
      {selectedPeriod === "custom" && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Data Inicial
            </label>
            <input
              type="date"
              value={customStartDate}
              onChange={(e) =>
                onCustomDateChange(e.target.value, customEndDate)
              }
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Data Final
            </label>
            <input
              type="date"
              value={customEndDate}
              onChange={(e) =>
                onCustomDateChange(customStartDate, e.target.value)
              }
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent"
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default AppointmentFilter;
