import { SLOT_MINUTES } from '../constants'

export function formatSlot(unix: number): string {
  return new Date(unix * 1000).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' })
}

export interface SlotOption {
  value: number
  label: string
}

const SLOT_SECONDS = SLOT_MINUTES * 60

export function slotOptions(date: Date): SlotOption[] {
  const out: SlotOption[] = []
  const d = new Date(date)
  d.setHours(8, 0, 0, 0)
  const end = new Date(d)
  end.setHours(18, 0, 0, 0)
  while (d < end) {
    const unix = Math.floor(d.getTime() / 1000)
    out.push({
      value: unix,
      label: d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' }),
    })
    d.setMinutes(d.getMinutes() + SLOT_MINUTES)
  }
  return out
}

/** Appointment-like shape to compute taken slots without importing api */
export interface AppointmentSlotSpan {
  start_time_unix: number
  duration_slots: number
}

/** Returns the set of slot start unix timestamps that are occupied on the given date. */
export function getTakenSlotUnixSet(
  appointments: AppointmentSlotSpan[],
  date: Date
): Set<number> {
  const dayStart = new Date(date)
  dayStart.setHours(0, 0, 0, 0)
  const dayEnd = new Date(date)
  dayEnd.setHours(23, 59, 59, 999)
  const dayStartUnix = Math.floor(dayStart.getTime() / 1000)
  const dayEndUnix = Math.floor(dayEnd.getTime() / 1000)
  const taken = new Set<number>()
  for (const a of appointments) {
    const aStart = a.start_time_unix
    const aEnd = a.start_time_unix + a.duration_slots * SLOT_SECONDS
    if (aEnd <= dayStartUnix || aStart > dayEndUnix) continue
    for (let t = aStart; t < aEnd; t += SLOT_SECONDS) {
      if (t >= dayStartUnix && t <= dayEndUnix) taken.add(t)
    }
  }
  return taken
}

export function slotDurationMinutes(slots: number): number {
  return slots * SLOT_MINUTES
}
