package integration_test

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	"github.com/stretchr/testify/require"
)

var _ = Describe("Calendar HTTP", func() {
	now := time.Now()

	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	client := http.Client{
		Transport: tr,
	}

	Describe("CreateEvent Error", func() {
		var checkError func(payload storage.Event)
		notification := 4 * time.Hour
		BeforeEach(func() {
			checkError = func(payload storage.Event) {
				data, err := json.Marshal(&payload)
				require.NoError(GinkgoT(), err)
				resp, err := client.Post(rootHTTPURL+"/event", "application/json", bytes.NewReader(data))
				require.NoError(GinkgoT(), err)
				defer resp.Body.Close()
				var response struct {
					Error string `json:"error"`
				}
				require.Equal(GinkgoT(), http.StatusInternalServerError, resp.StatusCode)
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(GinkgoT(), err)
				// require.NotEmpty(GinkgoT(), response.Error)
			}
		})
		It("add bad date event", func() {
			badEvent := storage.Event{
				Title:        gofakeit.Hobby(),
				Start:        now.Add(1 * time.Minute),
				Stop:         now,
				Description:  gofakeit.Phrase(),
				UserID:       gofakeit.IntRange(1, 200),
				Notification: &notification,
			}
			checkError(badEvent.Transfer())
		})
	})
	Describe("Id Error", func() {
		url := rootHTTPURL + "/event/10"
		It("del bad id", func() {
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			var response struct {
				Error string `json:"error"`
			}
			require.Equal(GinkgoT(), http.StatusOK, resp.StatusCode)
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", response.Error)
		})

		It("update bad args", func() {
			req, err := http.NewRequest(http.MethodPut, url, nil)
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
		})
		It("update bad id", func() {
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte("{}")))
			require.NoError(GinkgoT(), err)
			resp, err := client.Do(req)
			require.NoError(GinkgoT(), err)
			defer resp.Body.Close()
			var response struct {
				Error string `json:"error"`
			}
			require.Equal(GinkgoT(), http.StatusBadRequest, resp.StatusCode)
			// b, _ := io.ReadAll(resp.Body)
			// fmt.Println(string(b))
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(GinkgoT(), err)
			// require.Equal(GinkgoT(), "bad event id", response.Error)
			fmt.Println(response)
		})
	})
})
