export const canEditAppointment = (appointmentDate: string): { canEdit: boolean; reason?: string } => {
  const appointment = new Date(appointmentDate);
  const now = new Date();
  const diffInMs = appointment.getTime() - now.getTime();
  const diffInDays = diffInMs / (1000 * 60 * 60 * 24);

  if (diffInDays < 2) {
    return {
      canEdit: false,
      reason: "Não é possível editar agendamentos com menos de 2 dias de antecedência"
    };
  }

  return { canEdit: true };
};

export const canCancelAppointment = (appointmentStatus: string): { canCancel: boolean; reason?: string } => {
  if (appointmentStatus === "DONE") {
    return {
      canCancel: false,
      reason: "Não é possível cancelar agendamentos já concluídos"
    };
  }

  if (appointmentStatus === "CANCELED") {
    return {
      canCancel: false,
      reason: "Este agendamento já está cancelado"
    };
  }

  return { canCancel: true };
};