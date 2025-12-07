package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
	mocks "github.com/divijg19/physiolink/backend/tests/__mocks__"
)

func addChiURLParam(req *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func withUser(ctx context.Context, id string, role string) context.Context {
	ctx = context.WithValue(ctx, middleware.UserIDKey, id)
	if role != "" {
		ctx = context.WithValue(ctx, middleware.UserRoleKey, role)
	}
	return ctx
}

func TestCreateAvailability_Created(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{}
	handlers.InitAppointments(svc)

	body := map[string]interface{}{"slots": []map[string]string{{"startTime": "2025-12-05T10:00:00Z", "endTime": "2025-12-05T10:30:00Z"}}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/appointments/availability", bytes.NewReader(b))
	req = req.WithContext(withUser(req.Context(), "11111111-1111-1111-1111-111111111111", "pt"))
	rr := httptest.NewRecorder()

	handlers.CreateAvailability(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestGetTherapistAvailability_OK(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{SlotsResp: []service.Slot{{ID: uuid.New(), TherapistID: uuid.New(), StartTs: "2025-12-05T10:00:00Z", EndTs: "2025-12-05T10:30:00Z", Status: "open"}}}
	handlers.InitAppointments(svc)

	tid := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/api/appointments/availability?ptId="+tid, nil)
	rr := httptest.NewRecorder()
	handlers.GetTherapistAvailability(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetTherapistAvailability_MissingPtId_Returns400(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{}
	handlers.InitAppointments(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/appointments/availability", nil)
	rr := httptest.NewRecorder()
	handlers.GetTherapistAvailability(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetTherapistAvailability_BadUUID_Query_Returns400(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{}
	handlers.InitAppointments(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/appointments/availability?ptId=not-a-uuid", nil)
	rr := httptest.NewRecorder()
	handlers.GetTherapistAvailability(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetTherapistAvailability_BadUUID_Path_Returns400(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{}
	handlers.InitAppointments(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/appointments/availability/not-a-uuid", nil)
	req = addChiURLParam(req, "ptId", "not-a-uuid")
	rr := httptest.NewRecorder()
	handlers.GetTherapistAvailability(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestBookAppointment_OK(t *testing.T) {
	apptID := uuid.New()
	svc := &mocks.AppointmentServiceMock{BookResp: apptID}
	handlers.InitAppointments(svc)

	slotID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, "/api/appointments/"+slotID+"/book", nil)
	req = addChiURLParam(req, "id", slotID)
	req = req.WithContext(withUser(req.Context(), "11111111-1111-1111-1111-111111111111", "patient"))
	rr := httptest.NewRecorder()

	handlers.BookAppointment(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["id"] != apptID.String() {
		t.Fatalf("expected id %s, got %s", apptID.String(), resp["id"])
	}
}

func TestUpdateAppointmentStatus_OK(t *testing.T) {
	apptID := uuid.New()
	brief := service.AppointmentBrief{ID: apptID.String(), Status: "confirmed"}
	svc := &mocks.AppointmentServiceMock{UpdateResp: brief}
	handlers.InitAppointments(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/appointments/"+apptID.String()+"/status", bytes.NewReader([]byte(`{"status":"confirmed"}`)))
	req = addChiURLParam(req, "id", apptID.String())
	req = req.WithContext(withUser(req.Context(), "11111111-1111-1111-1111-111111111111", "pt"))
	rr := httptest.NewRecorder()

	handlers.UpdateAppointmentStatus(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetMyAppointments_OK(t *testing.T) {
	svc := &mocks.AppointmentServiceMock{ListResp: []service.AppointmentBrief{{ID: uuid.New().String(), Status: "booked"}}}
	handlers.InitAppointments(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/appointments/me", nil)
	req = req.WithContext(withUser(req.Context(), "11111111-1111-1111-1111-111111111111", "patient"))
	rr := httptest.NewRecorder()
	handlers.GetMyAppointments(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
