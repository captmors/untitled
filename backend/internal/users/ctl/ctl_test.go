package ctl_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"untitled/internal/users/ctl"
	"untitled/internal/users/mdl"
	"untitled/internal/users/repo"
	"untitled/internal/users/svc"
	"untitled/pkg/dockermocker"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var (
	router     *gin.Engine
	db         *gorm.DB
	userCtl    *ctl.UserCtl
	dockerTest *dockermocker.DockerTest
)

var _ = ginkgo.BeforeSuite(func() {
	dockerTest = dockermocker.NewDockerTest("1488", 120)
	db = dockerTest.OpenDatabaseConnection()
	db.AutoMigrate(&mdl.User{})

	userRepo := repo.NewUserRepo(db)
	userSvc := svc.NewUserSvc(userRepo)
	userCtl = ctl.NewUserCtl(userSvc)

	router = gin.Default()
	router.POST("/users", userCtl.CreateUser)
	router.GET("/users/:id", userCtl.GetUserByID)
})

var _ = ginkgo.AfterSuite(func() {
	dockerTest.Close()
})

var _ = ginkgo.Describe("User API", func() {
	ginkgo.Context("POST /users", func() {
		ginkgo.It("should create a user successfully", func() {
			userJSON := `{"name":"John Doe","email":"john@example.com","password":"123"}`
			req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusOK))

			var user mdl.User
			err := json.Unmarshal(resp.Body.Bytes(), &user)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(user.Name).To(gomega.Equal("John Doe"))
			gomega.Expect(user.Email).To(gomega.Equal("john@example.com"))
		})

		ginkgo.It("should return bad request on invalid input", func() {
			req, _ := http.NewRequest(http.MethodPost, "/users", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusBadRequest))
		})
	})

	ginkgo.Context("GET /users/:id", func() {
		var createdUserID int

		ginkgo.BeforeEach(func() {
			user := mdl.User{Name: "Jane Doe", Email: "jane@example.com", Password: "password123"}
			db.Create(&user)
			createdUserID = int(user.ID)
		})

		ginkgo.It("should retrieve a user by ID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(createdUserID), nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusOK))

			var user mdl.User
			err := json.Unmarshal(resp.Body.Bytes(), &user)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(user.ID).To(gomega.Equal(uint(createdUserID)))
		})

		ginkgo.It("should return not found for invalid ID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/users/9999", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			gomega.Expect(resp.Code).To(gomega.Equal(http.StatusNotFound))
		})
	})
})

func TestApi(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User API Suite")
}
