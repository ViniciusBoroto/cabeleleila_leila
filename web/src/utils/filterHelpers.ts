import type { FilterPeriod } from "../components/AppointmentFilter";

interface Service {
  id: number;
  name: string;
  price: number;
  duration_minutes: number;
}

interface User {
  id: number;
  email?: string;
  name?: string;
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

export const filterAppointmentsByPeriod = (
  appointments: Appointment[],
  period: FilterPeriod,
  customStartDate?: string,
  customEndDate?: string
): Appointment[] => {
  const now = new Date();

  switch (period) {
    case "all":
      return appointments;

    case "today": {
      const startOfDay = new Date(now);
      startOfDay.setHours(0, 0, 0, 0);
      const endOfDay = new Date(now);
      endOfDay.setHours(23, 59, 59, 999);

      return appointments.filter((appointment) => {
        const appointmentDate = new Date(appointment.date);
        return appointmentDate >= startOfDay && appointmentDate <= endOfDay;
      });
    }

    case "week": {
      const startOfWeek = new Date(now);
      const day = startOfWeek.getDay();
      const diff = startOfWeek.getDate() - day; // Domingo como início da semana
      startOfWeek.setDate(diff);
      startOfWeek.setHours(0, 0, 0, 0);

      const endOfWeek = new Date(startOfWeek);
      endOfWeek.setDate(startOfWeek.getDate() + 6);
      endOfWeek.setHours(23, 59, 59, 999);

      return appointments.filter((appointment) => {
        const appointmentDate = new Date(appointment.date);
        return appointmentDate >= startOfWeek && appointmentDate <= endOfWeek;
      });
    }

    case "month": {
      const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
      startOfMonth.setHours(0, 0, 0, 0);

      const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
      endOfMonth.setHours(23, 59, 59, 999);

      return appointments.filter((appointment) => {
        const appointmentDate = new Date(appointment.date);
        return (
          appointmentDate >= startOfMonth && appointmentDate <= endOfMonth
        );
      });
    }

    case "custom": {
      if (!customStartDate || !customEndDate) {
        return appointments;
      }

      const startDate = new Date(customStartDate);
      startDate.setHours(0, 0, 0, 0);

      const endDate = new Date(customEndDate);
      endDate.setHours(23, 59, 59, 999);

      return appointments.filter((appointment) => {
        const appointmentDate = new Date(appointment.date);
        return appointmentDate >= startDate && appointmentDate <= endDate;
      });
    }

    default:
      return appointments;
  }
};

export const getDefaultCustomDates = () => {
  const today = new Date();
  const startDate = new Date(today);
  startDate.setDate(today.getDate() - 7); // 7 dias atrás

  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  };

  return {
    start: formatDate(startDate),
    end: formatDate(today),
  };
};