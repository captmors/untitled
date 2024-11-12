package ctl_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/onsi/ginkgo/v2"
	log "github.com/sirupsen/logrus"

	"untitled/internal/testutils"
	"untitled/internal/users/mdl"

	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Auth API", func() {
	var logFile *os.File

	ginkgo.BeforeEach(func() {
		db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

		registerJSON := `{"name":"John Doe","email":"john@example.com","password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		resp := httptest.NewRecorder()
		authRouter.ServeHTTP(resp, req)

		gomega.Expect(resp.Code).To(gomega.Equal(http.StatusCreated))
	})

	ginkgo.AfterEach(func() {
		dockerTest.Close()
		logFile.Close()
	})

	ginkgo.Context("POST /auth/login", func() {
		ginkgo.It("should authenticate a user and return a JWT", func() {
			logFile := testutils.InitTestLogging("auth succ")
			logFiles = append(logFiles, logFile)

			loginJSON := `{"name":"John Doe","password":"password"}`
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginJSON))
			resp := httptest.NewRecorder()

			log.Info("Sending login request")
			authRouter.ServeHTTP(resp, req)

			log.Infof("Response Code: %d", resp.Code)
			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusOK))

			var response mdl.UserLoginResponse
			err := json.Unmarshal(resp.Body.Bytes(), &response)
			if err != nil {
				log.Errorf("Failed to unmarshal response: %v", err)
			}
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response.Token).To(gomega.Not(gomega.BeEmpty()))
			log.Info("User authenticated successfully, JWT returned")
		})

		ginkgo.It("should return unauthorized for invalid credentials", func() {
			logFile := testutils.InitTestLogging("auth unsucc")
			logFiles = append(logFiles, logFile)

			loginJSON := `{"name":"invalid_user","password":"wrongpassword"}`
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginJSON))
			resp := httptest.NewRecorder()

			log.Info("Sending login request with invalid credentials")
			authRouter.ServeHTTP(resp, req)

			log.Infof("Response Code: %d", resp.Code)
			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusUnauthorized))
			log.Info("Unauthorized response received as expected")
		})
	})
})
