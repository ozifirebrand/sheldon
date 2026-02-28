package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
	"scheduler/backend/internal/broadcast"
	"scheduler/backend/internal/constants"
	"scheduler/backend/internal/service"
)

const apiPrefixAppointments = "/api/appointments/"

type Gateway struct {
	schedulerService *service.Service
	broadcaster      *broadcast.Broadcaster
}

func New(schedulerService *service.Service, broadcaster *broadcast.Broadcaster) *Gateway {
	return &Gateway{schedulerService: schedulerService, broadcaster: broadcaster}
}

func (g *Gateway) Handler() http.Handler {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/api/appointments", g.handleAppointments)
	serveMux.HandleFunc("/api/appointments/stream", g.handleStream)
	serveMux.HandleFunc("/api/appointments/", g.handleAppointmentByID)
	serveMux.HandleFunc("/api/doctors", g.handleDoctors)
	return CORS(serveMux)
}

func (g *Gateway) handleDoctors(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(responseWriter, http.StatusOK, []map[string]string{{"id": constants.DefaultDoctorID, "name": constants.DefaultDoctorName}})
}

func (g *Gateway) handleAppointments(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		g.listAppointments(responseWriter, request)
	case http.MethodPost:
		g.createAppointment(responseWriter, request)
	default:
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) listAppointments(responseWriter http.ResponseWriter, request *http.Request) {
	list, err := g.schedulerService.ListAppointments(request.Context())
	if err != nil {
		writeError(responseWriter, request, err)
		return
	}
	if list == nil {
		list = []*schedulerv1.Appointment{}
	}
	writeJSON(responseWriter, http.StatusOK, list)
}

func (g *Gateway) createAppointment(responseWriter http.ResponseWriter, request *http.Request) {
	var body struct {
		DoctorID       string `json:"doctor_id"`
		DoctorName     string `json:"doctor_name"`
		StartTimeUnix  int64  `json:"start_time_unix"`
		DurationSlots  int32  `json:"duration_slots"`
		PatientName    string `json:"patient_name"`
		PatientContact string `json:"patient_contact"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(responseWriter, "invalid JSON", http.StatusBadRequest)
		return
	}
	appointment, err := g.schedulerService.CreateAppointment(request.Context(), &schedulerv1.CreateAppointmentRequest{
		DoctorId:       body.DoctorID,
		DoctorName:     body.DoctorName,
		StartTimeUnix:  body.StartTimeUnix,
		DurationSlots:  body.DurationSlots,
		PatientName:    body.PatientName,
		PatientContact: body.PatientContact,
	})
	if err != nil {
		writeError(responseWriter, request, err)
		return
	}
	writeJSON(responseWriter, http.StatusCreated, appointment)
}

func (g *Gateway) deleteAppointment(responseWriter http.ResponseWriter, request *http.Request) {
	id := appointmentIDFromPath(request.URL.Path, apiPrefixAppointments)
	if id == "" || strings.Contains(id, "/") {
		http.Error(responseWriter, "missing or invalid appointment id", http.StatusBadRequest)
		return
	}
	if err := g.schedulerService.DeleteAppointment(request.Context(), id); err != nil {
		writeError(responseWriter, request, err)
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}

func (g *Gateway) handleAppointmentByID(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(responseWriter, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	g.deleteAppointment(responseWriter, request)
}

func (g *Gateway) handleStream(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	channel, unsubscribe := g.broadcaster.Subscribe()
	defer unsubscribe()

	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")
	flusher, ok := responseWriter.(http.Flusher)
	if !ok {
		http.Error(responseWriter, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	for event := range channel {
		eventMap := eventToMap(event)
		data, err := json.Marshal(eventMap)
		if err != nil {
			log.Printf("gateway: stream marshal: %v", err)
			continue
		}
		if _, err := responseWriter.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
			log.Printf("gateway: stream write: %v", err)
			return
		}
		flusher.Flush()
	}
}

func eventToMap(event *schedulerv1.AppointmentEvent) map[string]interface{} {
	out := map[string]interface{}{
		"appointment_id": event.GetAppointmentId(),
		"action":         event.GetAction(),
	}
	if appointment := event.GetAppointment(); appointment != nil {
		out["appointment"] = map[string]interface{}{
			"id":              appointment.GetId(),
			"doctor_id":       appointment.GetDoctorId(),
			"doctor_name":     appointment.GetDoctorName(),
			"start_time_unix": appointment.GetStartTimeUnix(),
			"duration_slots":  appointment.GetDurationSlots(),
			"patient_name":    appointment.GetPatientName(),
			"patient_contact": appointment.GetPatientContact(),
		}
	}
	return out
}
