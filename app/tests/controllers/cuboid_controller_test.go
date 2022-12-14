package controllers_test

import (
	. "cuboid-challenge/app/models"
	"cuboid-challenge/app/tests/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cuboid Controller", func() {
	testutils.LoadEnv()
	testutils.ConnectDB()
	testutils.ClearDB()

	AfterEach(func() {
		testutils.ClearDB()
	})

	var w *httptest.ResponseRecorder
	var bag *Bag

	BeforeEach(func() {
		bag = &Bag{
			Title:  "A bag",
			Volume: 5,
			Cuboids: []Cuboid{
				{Width: 1, Height: 1, Depth: 1},
				{Width: 1, Height: 1, Depth: 2},
			},
		}
		testutils.AddRecords(bag)
	})

	Describe("List", func() {
		BeforeEach(func() {
			w = testutils.MockRequest(http.MethodGet, "/cuboids", nil)
		})

		It("Response HTTP status code 200", func() {
			Expect(w.Code).To(Equal(200))
		})

		It("Fetch all cuboids", func() {
			l, _ := testutils.DeserializeList(w.Body.String())
			Expect(len(l)).To(Equal(2))
			for _, m := range l {
				Expect(m["width"]).ToNot(BeNil())
				Expect(m["height"]).ToNot(BeNil())
				Expect(m["depth"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["bagId"]).To(BeEquivalentTo(bag.ID))
			}
		})
	})

	Describe("Get", func() {
		var cuboidID uint

		JustBeforeEach(func() {
			w = testutils.MockRequest(http.MethodGet, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})

		Context("When the cuboid is present", func() {
			BeforeEach(func() {
				cuboidID = bag.Cuboids[0].ID
			})

			It("Response HTTP status code 200", func() {
				Expect(w.Code).To(Equal(200))
			})

			It("Get the cuboid", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["id"]).To(BeEquivalentTo(bag.Cuboids[0].ID))
				Expect(m["width"]).ToNot(BeNil())
				Expect(m["height"]).ToNot(BeNil())
				Expect(m["depth"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["bagId"]).ToNot(BeNil())
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 9999
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
			})
		})
	})

	Describe("Create", func() {
		cuboidPayload := map[string]interface{}{}

		BeforeEach(func() {
			cuboidPayload = map[string]interface{}{
				"width":  1,
				"height": 1,
				"depth":  1,
				"bagId":  bag.ID,
			}
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(cuboidPayload)
			w = testutils.MockRequest(http.MethodPost, "/cuboids", &body)
		})

		It("Response HTTP status code 201", func() {
			Expect(w.Code).To(Equal(201))
		})

		It("Returns the created cuboid", func() {
			m, _ := testutils.Deserialize(w.Body.String())
			Expect(m["width"]).ToNot(BeNil())
			Expect(m["height"]).ToNot(BeNil())
			Expect(m["depth"]).ToNot(BeNil())
			Expect(m["volume"]).ToNot(BeNil())
			Expect(m["bagId"]).To(BeEquivalentTo(bag.ID))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				cuboidPayload["width"] = 3
			})

			It("Does not create the Cuboid", func() {
				Expect(w.Code).To(Equal(400))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When the bag is disabled", func() {
			BeforeEach(func() {
				bag.SetDisabled(true)
				bag.Cuboids = []Cuboid{}
				testutils.UpdateRecords(&bag)
			})

			It("Does not admit new cuboids", func() {
				Expect(w.Code).To(Equal(400))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Bag is disabled"))
			})
		})
	})

	// DO NOT modify the tests ABOVE
	// IMPLEMENT the tests BELOW

	Describe("Update", func() {
		var cuboidID uint

		BeforeEach(func() {
			cuboidID = bag.Cuboids[0].ID
		})

		updateCuboidPayload := map[string]interface{}{}

		BeforeEach(func() {
			updateCuboidPayload = map[string]interface{}{
				"width":  1,
				"height": 1,
				"depth":  2,
				"bagId":  bag.ID,
			}
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(updateCuboidPayload)
			w = testutils.MockRequest(http.MethodPut, "/cuboids/"+fmt.Sprintf("%v", cuboidID), &body)
		})

		It("Response HTTP status code 200", func() {
			Expect(w.Code).To(Equal(200))
		})

		It("Returns the updated cuboid", func() {
			m, _ := testutils.Deserialize(w.Body.String())
			Expect(m["width"]).To(BeEquivalentTo(updateCuboidPayload["width"]))
			Expect(m["height"]).To(BeEquivalentTo(updateCuboidPayload["height"]))
			Expect(m["depth"]).To(BeEquivalentTo(updateCuboidPayload["depth"]))
			Expect(m["volume"]).To(BeEquivalentTo(2))
			Expect(m["bagId"]).To(BeEquivalentTo(updateCuboidPayload["bagId"]))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				updateCuboidPayload["width"] = 3
			})

			It("Does not update the Cuboid", func() {
				Expect(w.Code).To(Equal(400))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 10
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("not Found"))
			})
		})
	})

	Describe("Delete", func() {

		var cuboidID uint
		BeforeEach(func() {
			cuboidID = bag.Cuboids[0].ID
		})
		JustBeforeEach(func() {
			w = testutils.MockRequest(http.MethodDelete, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})
		Context("When the cuboid is present", func() {
			It("Response HTTP status code 200", func() {
				Expect(w.Code).To(Equal(200))
			})

			It("Remove the cuboid", func() {
				w = testutils.MockRequest(http.MethodGet, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
				Expect(w.Code).To(Equal(404))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 10
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("not Found"))
			})
		})
	})
})
