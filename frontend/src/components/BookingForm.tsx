import type { Doctor } from '../api'
import { DURATION_SLOTS, CLINIC_NAME } from '../constants'
import type { SlotOption } from '../utils/slotUtils'

export interface BookingFormProps {
  doctors: Doctor[]
  doctorId: string
  setDoctorId: (v: string) => void
  date: string
  setDate: (v: string) => void
  slot: string
  setSlot: (v: string) => void
  duration: number
  setDuration: (v: number) => void
  patientName: string
  setPatientName: (v: string) => void
  patientContact: string
  setPatientContact: (v: string) => void
  formError: string | null
  onSubmit: (e: React.FormEvent) => void
  slots: SlotOption[]
  takenSlotUnixSet: Set<number>
}

export default function BookingForm({
  doctors,
  doctorId,
  setDoctorId,
  date,
  setDate,
  slot,
  setSlot,
  duration,
  setDuration,
  patientName,
  setPatientName,
  patientContact,
  setPatientContact,
  formError,
  onSubmit,
  slots,
  takenSlotUnixSet,
}: BookingFormProps) {
  return (
    <form
      onSubmit={onSubmit}
      style={{ marginBottom: '2rem', padding: '1rem', background: '#fff', borderRadius: 8 }}
    >
      <h2>Book an appointment at {CLINIC_NAME}</h2>
      {formError && <p style={{ color: 'crimson' }}>{formError}</p>}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem 1rem', maxWidth: 400 }}>
        <label>Doctor</label>
        <select value={doctorId} onChange={(e) => setDoctorId(e.target.value)}>
          {doctors.map((d) => (
            <option key={d.id} value={d.id}>{d.name}</option>
          ))}
        </select>
        <label>Date</label>
        <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        <label>Time slot</label>
        <div style={{ gridColumn: '1 / -1' }}>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(80px, 1fr))',
              gap: 4,
              maxHeight: 140,
              overflowY: 'auto',
              padding: 4,
              border: '1px solid #ccc',
              borderRadius: 4,
              background: '#fafafa',
            }}
          >
            <label
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 6,
                padding: '4px 6px',
                borderRadius: 4,
                cursor: 'pointer',
                background: slot === '' ? '#e0e7ff' : 'transparent',
              }}
            >
              <input
                type="radio"
                name="slot"
                value=""
                checked={slot === ''}
                onChange={() => setSlot('')}
                style={{ margin: 0 }}
              />
              <span style={{ color: '#666' }}>Select...</span>
            </label>
            {slots.map((s) => {
              const taken = takenSlotUnixSet.has(s.value)
              return (
                <label
                  key={s.value}
                  title={taken ? 'This time slot is taken.' : undefined}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 6,
                    padding: '4px 6px',
                    borderRadius: 4,
                    cursor: taken ? 'default' : 'pointer',
                    opacity: taken ? 0.85 : 1,
                    background: slot === String(s.value) ? '#e0e7ff' : 'transparent',
                  }}
                >
                  <input
                    type="radio"
                    name="slot"
                    value={s.value}
                    checked={slot === String(s.value)}
                    onChange={() => !taken && setSlot(String(s.value))}
                    disabled={taken}
                    style={{ margin: 0 }}
                  />
                  {taken && (
                    <span
                      style={{
                        width: 6,
                        height: 6,
                        borderRadius: '50%',
                        background: '#c94a4a',
                        flexShrink: 0,
                      }}
                      aria-hidden
                    />
                  )}
                  <span>{s.label}</span>
                </label>
              )
            })}
          </div>
          <p style={{ margin: '4px 0 0', fontSize: 12, color: '#666' }}>
            {takenSlotUnixSet.size > 0 && (
              <span title="This time slot is taken.">
                <span style={{ display: 'inline-block', width: 6, height: 6, borderRadius: '50%', background: '#c94a4a', verticalAlign: 'middle', marginRight: 4 }} />
                Red dot = taken
              </span>
            )}
          </p>
        </div>
        <label>Duration</label>
        <select value={duration} onChange={(e) => setDuration(Number(e.target.value))}>
          {DURATION_SLOTS.map((d) => (
            <option key={d.value} value={d.value}>{d.label}</option>
          ))}
        </select>
        <label>Patient name</label>
        <input
          type="text"
          value={patientName}
          onChange={(e) => setPatientName(e.target.value)}
          placeholder="Full name"
          required
        />
        <label>Contact</label>
        <input
          type="text"
          value={patientContact}
          onChange={(e) => setPatientContact(e.target.value)}
          placeholder="Email or phone"
        />
      </div>
      <button type="submit" style={{ marginTop: '0.5rem' }}>Book appointment</button>
    </form>
  )
}
