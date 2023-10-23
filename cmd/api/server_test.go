package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const port = 80

func TestCartCreationAndUpdate(t *testing.T) {
	server := NewServer(port)
	content, err := json.Marshal(CartRequest{
		UserID: "123",
	})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/api/v1/cart", bytes.NewReader(content))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusCreated, recorder.Code)
	cartResp := CartResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &cartResp)
	require.NoError(t, err)
	assert.Equal(t, "123", cartResp.UseID)

	// and given
	request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/cart/%s", cartResp.CartID), bytes.NewReader(content))
	require.NoError(t, err)
	recorder = httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusOK, recorder.Code)
	cartRespGet := CartResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &cartRespGet)
	require.NoError(t, err)
	assert.Equal(t, cartResp.UseID, cartRespGet.UseID)
	assert.Equal(t, cartResp.CartID, cartRespGet.CartID)

	// and given
	content, err = json.Marshal(productsUpdateRequest{
		ProductID: "1",
		Quantity:  2,
	})
	require.NoError(t, err)

	request, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/cart/%s", cartResp.CartID), bytes.NewReader(content))
	require.NoError(t, err)
	recorder = httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusOK, recorder.Code)
	err = json.Unmarshal(recorder.Body.Bytes(), &cartResp)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cartResp.Products))
	assert.Equal(t, "1", cartResp.Products[0].ProductID)
	assert.Equal(t, int32(2), cartResp.Products[0].Quantity)

	// and given
	content, err = json.Marshal(productsUpdateRequest{
		ProductID: "1",
		Quantity:  7,
	})
	require.NoError(t, err)

	request, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/cart/%s/product/1", cartResp.CartID), bytes.NewReader(content))
	require.NoError(t, err)
	recorder = httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusOK, recorder.Code)
	err = json.Unmarshal(recorder.Body.Bytes(), &cartResp)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cartResp.Products))
	assert.Equal(t, "1", cartResp.Products[0].ProductID)
	assert.Equal(t, int32(7), cartResp.Products[0].Quantity)
}
