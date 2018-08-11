package fhir

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/minhajuddinkhan/fhir/models"
	"github.com/minhajuddinkhan/gatewaygo/redox/models/common"
	"github.com/minhajuddinkhan/gatewaygo/redox/models/scheduling"
)

func NewAppointment(RedoxPayload scheduling.New) (interface{}, error) {

	provider, err := func() (common.Provider, error) {

		fmt.Println("WHAT", (RedoxPayload.Visit.AttendingProvider.ID))
		if len(RedoxPayload.Visit.AttendingProvider.ID) > 0 {
			return RedoxPayload.Visit.AttendingProvider, nil
		}
		if len(RedoxPayload.Visit.ConsultingProvider.ID) > 0 {
			return RedoxPayload.Visit.ConsultingProvider, nil
		}
		if len(RedoxPayload.Visit.ReferringProvider.ID) > 0 {
			return RedoxPayload.Visit.ReferringProvider, nil
		}

		if len(RedoxPayload.Visit.VisitProvider.ID) > 0 {
			return RedoxPayload.Visit.VisitProvider, nil
		}

		return common.Provider{}, errors.New("No Provider")
	}()

	if err != nil {
		panic("Practitioner not found.")
	}

	fhirDateLayout := "2006-01-02T15:04:05.999999999Z07:00"
	var startTime time.Time

	a := models.Appointment{
		Identifier: []models.Identifier{
			{
				Value: RedoxPayload.Visit.VisitNumber,
			},
		},
		Participant: []models.AppointmentParticipantComponent{
			{
				Actor: &models.Reference{
					Reference: "Patient/" + RedoxPayload.Patient.Identifiers[0].ID,
					Display:   "EPI",
				},
			},
			{
				Actor: &models.Reference{
					Reference: "Practitioner/" + provider.ID,
				},
			},
		},
		Status: strings.ToLower(RedoxPayload.Visit.Status),
		Start: &models.FHIRDateTime{
			Time: func() time.Time {
				t, err := time.Parse(fhirDateLayout, RedoxPayload.Visit.VisitDateTime)
				if err != nil {
					fmt.Println("err\n", err.Error())
				}
				startTime = t
				return t
			}(),
		},
		MinutesDuration: RedoxPayload.Visit.Duration,
		End: &models.FHIRDateTime{
			Precision: "15",
			Time: func() time.Time {
				duration, _ := time.ParseDuration(fmt.Sprint(RedoxPayload.Visit.Duration))
				startTime.Add(time.Minute * duration)
				return startTime
			}(),
		},
		ServiceType: &models.CodeableConcepts{
			{
				Text: RedoxPayload.Visit.Location.Facility,
			},
		},
		AppointmentType: &models.CodeableConcepts{
			{
				Text: strings.Split(RedoxPayload.Visit.Reason, ":")[0],
			},
		},
		ServiceCategory: &models.CodeableConcept{
			Text: RedoxPayload.Visit.Location.Department,
		},
	}

	return a.GetBSON()
}
