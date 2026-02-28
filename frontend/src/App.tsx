import { useEffect, useState } from 'react'
import {
  listAppointments,
  createAppointment,
  deleteAppointment,
  getDoctors,
  subscribeAppointmentUpdates,
  type Appointment,
  type Doctor,
} from './api'
import BookingForm from './components/BookingForm'
import AppointmentList from './components/AppointmentList'
import { slotOptions, getTakenSlotUnixSet } from './utils/slotUtils'
import { CLINIC_NAME } from './constants'

function normalizeList(list: Appointment[] | null | undefined): Appointment[] {
  return Array.isArray(list) ? list : []
}

function mergeAndSort(appointments: Appointment[], created: Appointment): Appointment[] {
  if (appointments.some((a) => a.id === created.id)) return appointments
  return [...appointments, created].sort((a, b) => a.start_time_unix - b.start_time_unix)
}

export default function App() {
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formError, setFormError] = useState<string | null>(null)
  const [doctorId, setDoctorId] = useState('')
  const [date, setDate] = useState(() => {
    const d = new Date()
    d.setHours(0, 0, 0, 0)
    return d.toISOString().slice(0, 10)
  })
  const [slot, setSlot] = useState('')
  const [duration, setDuration] = useState<number>(1)
  const [patientName, setPatientName] = useState('')
  const [patientContact, setPatientContact] = useState('')

  const refetch = async () => {
    try {
      const [list, docList] = await Promise.all([listAppointments(), getDoctors()])
      const doctorsList: Doctor[] = Array.isArray(docList) ? docList : []
      setAppointments(normalizeList(list))
      setDoctors(doctorsList)
      if (doctorsList.length) {
        setDoctorId((prev) => prev || doctorsList[0].id)
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    refetch()
  }, [])

  useEffect(() => {
    const unsub = subscribeAppointmentUpdates(() => refetch())
    return unsub
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setFormError(null)
    if (!slot) {
      setFormError('Please select a time slot.')
      return
    }
    if (!patientName.trim()) {
      setFormError('Please enter your name.')
      return
    }
    try {
      const selectedDoctor = doctors.find((d) => d.id === (doctorId || ''))
      const created = await createAppointment({
        doctor_id: doctorId || undefined,
        doctor_name: selectedDoctor?.name,
        start_time_unix: parseInt(slot, 10),
        duration_slots: duration,
        patient_name: patientName.trim(),
        patient_contact: patientContact.trim(),
      })
      setPatientName('')
      setPatientContact('')
      setSlot('')
      setAppointments((prev) => mergeAndSort(normalizeList(prev), created))
      await refetch()
    } catch (e) {
      setFormError(e instanceof Error ? e.message : 'Booking failed. Please try again.')
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteAppointment(id)
      await refetch()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Delete failed')
    }
  }

  const selectedDate = new Date(date + 'T12:00:00')
  const slots = slotOptions(selectedDate)
  const takenSlotUnixSet = getTakenSlotUnixSet(appointments, selectedDate)

  if (loading) {
    return <div style={{ padding: '1rem', fontFamily: 'system-ui' }}>Loading...</div>
  }
  if (error) {
    return (
      <div style={{ padding: '1rem', fontFamily: 'system-ui', color: 'crimson' }}>
        {error}
      </div>
    )
  }

  return (
    <div style={{ fontFamily: 'system-ui', color: '#111' }}>
      <h1>{CLINIC_NAME} – Appointments</h1>
      <BookingForm
        doctors={doctors}
        doctorId={doctorId}
        setDoctorId={setDoctorId}
        date={date}
        setDate={setDate}
        slot={slot}
        setSlot={setSlot}
        duration={duration}
        setDuration={setDuration}
        patientName={patientName}
        setPatientName={setPatientName}
        patientContact={patientContact}
        setPatientContact={setPatientContact}
        formError={formError}
        onSubmit={handleSubmit}
        slots={slots}
        takenSlotUnixSet={takenSlotUnixSet}
      />
      <section>
        <h2>Your appointments</h2>
        <AppointmentList appointments={appointments} onDelete={handleDelete} />
      </section>
    </div>
  )
}
