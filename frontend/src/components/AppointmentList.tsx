import type { Appointment } from '../api'
import { formatSlot, slotDurationMinutes } from '../utils/slotUtils'
import { CLINIC_NAME } from '../constants'
import './AppointmentList.css'

export interface AppointmentListProps {
  appointments: Appointment[]
  onDelete: (id: string) => void
}

export default function AppointmentList({ appointments, onDelete }: AppointmentListProps) {
  if (appointments.length === 0) {
    return <p className="emptyMessage">No appointments at {CLINIC_NAME}.</p>
  }
  return (
    <ul className="appointmentList">
      {appointments.map((a) => (
        <li key={a.id} className="appointmentItem">
          <span>
            {formatSlot(a.start_time_unix)} – {slotDurationMinutes(a.duration_slots)} min · {a.patient_name}
            {a.patient_contact && ` · ${a.patient_contact}`}
          </span>
          <button type="button" onClick={() => onDelete(a.id)} className="cancelBtn">
            Cancel
          </button>
        </li>
      ))}
    </ul>
  )
}
