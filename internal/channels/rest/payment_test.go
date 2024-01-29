package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/channels/rest"
	"tech-challenge-payment/internal/mocks"
	"tech-challenge-payment/internal/service"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	errorProcessingID = "PAYMENT_ERROR_PROCESSING"
)

func TestRegisterGroup(t *testing.T) {
	endpoint := "/payment"

	type Given struct {
		group          *echo.Group
		paymenyService service.PaymentService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid group, should register endpoints successfully": {
			given: Given{
				group:          echo.New().Group("/payment"),
				paymenyService: &mocks.PaymentServiceMock{},
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tests {
		p := rest.NewPaymentChannel(tc.given.paymenyService)
		p.RegisterGroup(tc.given.group)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, endpoint+"/123", nil)
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("123")

		e.ServeHTTP(rec, req)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)
	}
}

func TestCreate(t *testing.T) {
	endpoint := "/payment"

	type Given struct {
		request        *http.Request
		paymenyService service.PaymentService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json income must process normally": {
			given: Given{
				request:        createJsonRequest(http.MethodPost, endpoint, rest.PaymentRequest{}),
				paymenyService: mockPaymentServiceForCreate(canonical.Payment{}, canonical.Payment{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given wrong format must return error": {
			given: Given{
				request:        createRequest(http.MethodPost, endpoint),
				paymenyService: mockPaymentServiceForCreate(canonical.Payment{}, canonical.Payment{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalid data, must return application error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentRequest{
					PaymentType: 0,
					Status:      0,
					OrderID:     "asdasdasdasd",
				}),
				paymenyService: mockPaymentServiceForCreate(canonical.Payment{}, canonical.Payment{}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		err := rest.NewPaymentChannel(tc.given.paymenyService).Create(echo.New().NewContext(tc.given.request, rec))
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestCallback(t *testing.T) {
	endpoint := "/payment/callback"

	type Given struct {
		request        *http.Request
		paymenyService service.PaymentService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given normal json with status ok income must process normally as ok": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentCallback{
					PaymentID: "1234",
					Status:    "OK",
				}),
				paymenyService: mockPaymentServiceForCallback("1234", canonical.PAYMENT_PAYED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given normal json with status error income must process normally as error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentCallback{
					PaymentID: "1234",
					Status:    "ERROR",
				}),
				paymenyService: mockPaymentServiceForCallback("1234", canonical.PAYMENT_FAILED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given normal json with empty status income must process normally as error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentCallback{
					PaymentID: "1234",
					Status:    "",
				}),
				paymenyService: mockPaymentServiceForCallback("1234", canonical.PAYMENT_FAILED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given normal json with unkown status income must process normally as error": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentCallback{
					PaymentID: "1234",
					Status:    "asdasdasd",
				}),
				paymenyService: mockPaymentServiceForCallback("1234", canonical.PAYMENT_FAILED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given application error, must return statuscode 500": {
			given: Given{
				request: createJsonRequest(http.MethodPost, endpoint, rest.PaymentCallback{
					PaymentID: errorProcessingID,
					Status:    "",
				}),
				paymenyService: mockPaymentServiceForCallback("", canonical.PAYMENT_FAILED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusInternalServerError,
			},
		},
		"given invalid data, must return bad request": {
			given: Given{
				request:        createJsonRequest(http.MethodPost, endpoint, rest.PaymentRequest{}),
				paymenyService: mockPaymentServiceForCallback("", canonical.PAYMENT_FAILED),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		err := rest.NewPaymentChannel(tc.given.paymenyService).Callback(echo.New().NewContext(tc.given.request, rec))
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func TestGetByID(t *testing.T) {
	endpoint := "/payment/"

	type Given struct {
		request        *http.Request
		pathParamID    string
		paymenyService service.PaymentService
	}
	type Expected struct {
		err        assert.ErrorAssertionFunc
		statusCode int
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid id returns valid payment and status 200": {
			given: Given{
				request:     createRequest(http.MethodGet, endpoint),
				pathParamID: "1234",
				paymenyService: mockPaymentServiceForGetByID("1234", &canonical.Payment{
					ID: "1234",
				}),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusOK,
			},
		},
		"given empty id returns no payment and status 400": {
			given: Given{
				request:        createRequest(http.MethodGet, endpoint),
				pathParamID:    "",
				paymenyService: mockPaymentServiceForGetByID("1234", nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusBadRequest,
			},
		},
		"given invalic id returns no payment and status 404": {
			given: Given{
				request:        createRequest(http.MethodGet, endpoint),
				pathParamID:    errorProcessingID,
				paymenyService: mockPaymentServiceForGetByID("1234", nil),
			},
			expected: Expected{
				err:        assert.NoError,
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		e := echo.New().NewContext(tc.given.request, rec)
		e.SetPath("/:id")
		e.SetParamNames("id")

		e.SetParamValues(tc.given.pathParamID)
		err := rest.NewPaymentChannel(tc.given.paymenyService).GetByID(e)
		statusCode := rec.Result().StatusCode

		assert.Equal(t, tc.expected.statusCode, statusCode)

		tc.expected.err(t, err)
	}
}

func createRequest(method, endpoint string) *http.Request {
	req := createJsonRequest(method, endpoint, nil)
	req.Header.Del("Content-Type")
	return req
}

func createJsonRequest(method, endpoint string, request interface{}) *http.Request {
	json, _ := json.Marshal(request)
	req := httptest.NewRequest(method, endpoint, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func mockPaymentServiceForCreate(paymentReceived, paymentReturned canonical.Payment) *mocks.PaymentServiceMock {
	mockPaymentSvc := new(mocks.PaymentServiceMock)
	mockPaymentSvc.On("Create", mock.Anything, paymentReceived).Return(&paymentReturned, nil)
	mockPaymentSvc.On("Create", mock.Anything, canonical.Payment{
		OrderID: "asdasdasdasd",
	}).Return(&paymentReturned, errors.New(""))
	return mockPaymentSvc
}

func mockPaymentServiceForCallback(paymentID string, paymentStatus canonical.PaymentStatus) *mocks.PaymentServiceMock {
	mockPaymentSvc := new(mocks.PaymentServiceMock)

	mockPaymentSvc.
		On("Callback", mock.Anything, paymentID, paymentStatus).
		Return(nil)
	mockPaymentSvc.
		On("Callback", mock.Anything, errorProcessingID, mock.Anything).
		Return(errors.New(""))

	return mockPaymentSvc
}

func mockPaymentServiceForGetByID(paymentID string, paymentReturned *canonical.Payment) *mocks.PaymentServiceMock {
	mockPaymentSvc := new(mocks.PaymentServiceMock)

	mockPaymentSvc.
		On("GetByID", mock.Anything, paymentID).
		Return(paymentReturned, nil)
	mockPaymentSvc.
		On("GetByID", mock.Anything, errorProcessingID).
		Return(paymentReturned, errors.New(""))

	return mockPaymentSvc
}
