package handlers

import (
	"bytes"
	"context"
	"errors"
	"finalAssing/internal/auth"
	"finalAssing/internal/middleware"
	"finalAssing/internal/models"
	"finalAssing/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func Test_handler_AcceptApplicant(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, services.Service)
		expectedStatusCode int
		ExpectedResponse   string
	}{
		{
			name: "Missing tracker Id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httptest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", nil)
				c.Request = httptest
				return c, rr, nil
			},
			expectedStatusCode: 500,
			ExpectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TrackerIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "Invalid Request Body",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{Invalid Body}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TrackerIdKey, "123")
				ctx = context.WithValue(ctx, auth.AuthKey, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: 400,
			ExpectedResponse:   `{"msg2":"Bad Request"}`,
		},
		{
			name: "Unccessful validation",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`
				[
					{
						"Name":           "Vikalp Tyagi",
						"JobId":          1,
						"Experience":     3,
						"Max_NP":         2,
						"Budget":         50000,
						"Locations":      [1, 2, 3],
						"Stack":          [1,2,3],
						"WorkMode":       "Full-Time",
						"Qualifications": [1],
						"Shift":          "Day"
					}
				]
				`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TrackerIdKey, "12")
				ctx = context.WithValue(ctx, auth.AuthKey, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: 400,
			ExpectedResponse:   `{"Error":"All fields are mandatory"}`,
		},
		{
			name: "Unsuccess error in getting response",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`
				[
					{
						"name": "Vikalp Tyagi",
						"job": 1,
						"experience": 3,
						"noticePeriode": 2,
						"salary": 50000,
						"locations": [1, 2, 3],
						"skills": [1, 2, 3],
						"WorkMode": "Full-Time",
						"qualification": [1],
						"Shift": "Day"
					},
					{
						"name": "Akash",
						"job": 1,
						"experience": 3,
						"noticePeriode": 2,
						"salary": 90000,
						"locations": [1, 2, 3],
						"skills": [1, 2, 3],
						"WorkMode": "Full-Time",
						"qualification": [1],
						"Shift": "Day"
					}
				]
				`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TrackerIdKey, "12")
				ctx = context.WithValue(ctx, auth.AuthKey, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().FIlterApplication(ctx, gomock.Any()).Return(nil, errors.New("test error")).Times(1)

				return c, rr, ms
			},
			expectedStatusCode: 500,
			ExpectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "Success Case",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`
				[
					{
						"name": "Vikalp Tyagi",
						"job": 1,
						"experience": 3,
						"noticePeriode": 2,
						"salary": 50000,
						"locations": [1, 2, 3],
						"skills": [1, 2, 3],
						"WorkMode": "Full-Time",
						"qualification": [1],
						"Shift": "Day"
					},
					{
						"name": "Akash",
						"job": 1,
						"experience": 3,
						"noticePeriode": 2,
						"salary": 90000,
						"locations": [1, 2, 3],
						"skills": [1, 2, 3],
						"WorkMode": "Full-Time",
						"qualification": [1],
						"Shift": "Day"
					}
				]
				`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TrackerIdKey, "12")
				ctx = context.WithValue(ctx, auth.AuthKey, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().FIlterApplication(ctx, gomock.Any()).Return([]*models.ApplicantRespo{
					{
						Name:  "Vikalp Tyagi",
						JobId: 1,
					}}, nil).Times(1)

				return c, rr, ms
			},
			expectedStatusCode: 200,
			ExpectedResponse:   `[{"Name":"Vikalp Tyagi","JobId":1}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				s: ms,
			}
			h.AcceptApplicant(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.ExpectedResponse, rr.Body.String())

		})
	}
}
