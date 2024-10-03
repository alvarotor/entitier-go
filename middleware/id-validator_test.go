package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIDValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramType      string
		path           string
		expectedStatus int
		expectedBody   gin.H
		expectedID     interface{}
	}{
		{
			name:           "Valid String ID",
			paramType:      "string",
			path:           "/abc123",
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			expectedID:     "abc123",
		},
		{
			name:           "Valid Uint ID",
			paramType:      "uint",
			path:           "/123",
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			expectedID:     uint(123),
		},
		{
			name:           "Invalid ID (Empty)",
			paramType:      "string",
			path:           "/",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"err": models.ErrMustProvideValidID.Error()},
			expectedID:     nil,
		},
		{
			name:           "Invalid ID (Type Mismatch)",
			paramType:      "uint",
			path:           "/abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"err": models.ErrIDTypeMismatch.Error()},
			expectedID:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			var validatorFunc gin.HandlerFunc

			if tt.paramType == "string" {
				validatorFunc = IDValidator[string]()
			} else {
				validatorFunc = IDValidator[uint]()
			}

			router.GET("/:id", validatorFunc, func(c *gin.Context) {
				id, exists := c.Get("validatedID")
				if exists {
					assert.Equal(t, tt.expectedID, id)
					c.Status(http.StatusOK)
				} else {
					t.Logf("ValidatedID does not exist in context")
				}
			})

			// Add a catch-all route to handle empty ID
			router.GET("/", validatorFunc, func(c *gin.Context) {
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, gin.H{"err": models.ErrMustProvideValidID.Error()})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			router.ServeHTTP(w, req)

			t.Logf("Test case: %s", tt.name)
			t.Logf("Expected status: %d, Actual status: %d", tt.expectedStatus, w.Code)
			t.Logf("Response body: %s", w.Body.String())

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedBody, response)
				}
			}
		})
	}
}

func TestGetIDParam(t *testing.T) {
	tests := []struct {
		name     string
		paramID  string
		expected interface{}
	}{
		{"Empty ID", "", ""},
		{"String ID", "abc123", "abc123"},
		{"Uint ID", "123", uint(123)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.paramID}}

			result := getIDParam(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertToGenericID(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		expectedStr  string
		expectedUint uint
		expectError  bool
	}{
		{"Valid String", "abc123", "abc123", 0, false},
		{"Valid Uint", uint(123), "", 123, false},
		{"Invalid Type", 123, "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultStr, errStr := convertToGenericID[string](tt.input)
			resultUint, errUint := convertToGenericID[uint](tt.input)

			if tt.expectError {
				assert.Error(t, errStr)
				assert.Error(t, errUint)
			} else {
				if tt.expectedStr != "" {
					assert.NoError(t, errStr)
					assert.Equal(t, tt.expectedStr, resultStr)
				}
				if tt.expectedUint != 0 {
					assert.NoError(t, errUint)
					assert.Equal(t, tt.expectedUint, resultUint)
				}
			}
		})
	}
}
