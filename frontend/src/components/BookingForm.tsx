import type { Doctor } from '../api'
import { DURATION_SLOTS, CLINIC_NAME } from '../constants'
import type { SlotOption } from '../utils/slotUtils'
import './BookingForm.css'

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
    <form onSubmit={onSubmit} className="bookingForm">
      <h2>Book an appointment at {CLINIC_NAME}</h2>
      {formError && <p className="formError">{formError}</p>}
      <div className="formGrid">
        <label>Doctor</label>
        <select value={doctorId} onChange={(e) => setDoctorId(e.target.value)}>
          {doctors.map((d) => (
            <option key={d.id} value={d.id}>{d.name}</option>
          ))}
        </select>
        <label>Date</label>
        <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        <label>Time slot</label>
        <div className="slotPickerWrap">
          <div className="slotPicker">
            <label
              className={`slotOption ${slot === '' ? 'slotOptionSelected' : ''}`}
            >
              <input
                type="radio"
                name="slot"
                value=""
                checked={slot === ''}
                onChange={() => setSlot('')}
                className="radio"
              />
              <span className="selectPlaceholder">Select...</span>
            </label>
            {slots.map((s) => {
              const taken = takenSlotUnixSet.has(s.value)
              return (
                <label
                  key={s.value}
                  title={taken ? 'This time slot is taken.' : undefined}
                  className={`slotOption ${slot === String(s.value) ? 'slotOptionSelected' : ''} ${taken ? 'slotOptionTaken' : ''}`}
                >
                  <input
                    type="radio"
                    name="slot"
                    value={s.value}
                    checked={slot === String(s.value)}
                    onChange={() => !taken && setSlot(String(s.value))}
                    disabled={taken}
                    className="radio"
                  />
                  {taken && (
                    <span className="slotDot" aria-hidden />
                  )}
                  <span className="slotTime">{s.label}</span>
                </label>
              )
            })}
          </div>
          <p className="slotLegend">
            {takenSlotUnixSet.size > 0 && (
              <span title="This time slot is taken.">
                <span className="slotLegendDot" />
                Black dot = taken
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
      <button type="submit" className="submitBtn">Book appointment</button>
    </form>
  )
}
