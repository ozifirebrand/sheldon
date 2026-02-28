import { API_BASE } from '../config'
import { API_PATH_APPOINTMENTS, API_PATH_DOCTORS } from '../constants'

const api = (path: string) => `${API_BASE}${path}`

export interface Appointment {
  id: string
  doctor_id: string
  doctor_name: string
  start_time_unix: number
  duration_slots: number
  patient_name: string
  patient_contact: string
}

export interface Doctor {
  id: string
  name: string
}

function ensureArray<T>(x: T[] | null | undefined): T[] {
  return Array.isArray(x) ? x : []
}

export async function listAppointments(from?: number, to?: number): Promise<Appointment[]> {
  const params = new URLSearchParams()
  if (from != null) params.set('from', String(from))
  if (to != null) params.set('to', String(to))
  const res = await fetch(api(`${API_PATH_APPOINTMENTS}?${params}`))
  if (!res.ok) throw new Error(await res.text())
  return ensureArray(await res.json())
}

export async function createAppointment(data: {
  doctor_id?: string
  doctor_name?: string
  start_time_unix: number
  duration_slots: number
  patient_name: string
  patient_contact: string
}): Promise<Appointment> {
  const res = await fetch(api(API_PATH_APPOINTMENTS), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (res.status === 409) {
    const msg = await res.text()
    throw new Error(msg || 'This time slot is already booked.')
  }
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function deleteAppointment(id: string): Promise<void> {
  const res = await fetch(api(`${API_PATH_APPOINTMENTS}/${id}`), { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

export async function getDoctors(): Promise<Doctor[]> {
  const res = await fetch(api(API_PATH_DOCTORS))
  if (!res.ok) throw new Error(await res.text())
  return ensureArray(await res.json())
}

export interface AppointmentEvent {
  action: string
  appointment_id: string
  appointment?: Appointment
}

export function subscribeAppointmentUpdates(onEvent: (event: AppointmentEvent) => void): () => void {
  try {
    const es = new EventSource(api(`${API_PATH_APPOINTMENTS}/stream`))
    es.onmessage = (e) => {
      try {
        onEvent(JSON.parse(e.data))
      } catch {
      }
    }
    es.onerror = () => es.close()
    return () => es.close()
  } catch {
    return () => {}
  }
}
