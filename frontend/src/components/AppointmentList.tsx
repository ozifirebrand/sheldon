import type { Appointment } from '../api'
import { formatSlot, slotDurationMinutes } from '../utils/slotUtils'
import { CLINIC_NAME } from '../constants'

export interface AppointmentListProps {
  appointments: Appointment[]
  onDelete: (id: string) => void
}

export default function AppointmentList({ appointments, onDelete }: AppointmentListProps) {
  if (appointments.length === 0) {
    return <p>No appointments at {CLINIC_NAME}.</p>
  }
  return (
    <ul style={{ listStyle: 'none', padding: 0 }}>
      {appointments.map((a) => (
        <li
          key={a.id}
          style={{
            padding: '0.5rem 0.75rem',
            background: '#fff',
            marginBottom: 4,
            borderRadius: 4,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <span>
            {formatSlot(a.start_time_unix)} – {slotDurationMinutes(a.duration_slots)} min · {a.patient_name}
            {a.patient_contact && ` · ${a.patient_contact}`}
          </span>
          <button type="button" onClick={() => onDelete(a.id)} style={{ marginLeft: 8 }}>
            Cancel
          </button>
        </li>
      ))}
    </ul>
  )
}
