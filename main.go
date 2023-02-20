package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// ProductValidator echo validator for product
type ProductValidator struct {
	validator *validator.Validate
}

// Validate validates product request body
func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

func main() {
	// export MY_APP_PORT=4000
	port := os.Getenv("MY_APP_PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	v := validator.New()

	products := []map[int]string{{1: "mobile"}, {2: "laptop"}, {3: "macbook"}, {4: "iphone"}}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})
	e.POST("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "POST METHOD")
	})

	e.GET("/products/:id", func(c echo.Context) error {
		var product map[int]string
		for _, p := range products {
			for k := range p {
				pID, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					return err
				}
				if pID == k {
					product = p
				}
			}
		}
		if product == nil {
			return c.JSON(http.StatusNotFound, "product not found")
		}
		return c.JSON(http.StatusOK, product)
	})
	e.GET("/products", func(c echo.Context) error {
		return c.JSON(http.StatusOK, products)
	})
	e.POST("/products", func(c echo.Context) error {
		type body struct {
			Name string `json:"product_name" validate:"required,min=4"`
			// Vendor          string `json:"vendor" validate:"min=5,max=10`
			// Email           string `json:"email" validate:""required_with=Vendor,email`
			// Website         string `json:"website" validate:"url`
			// Country         string `json:"country" validate:"len=2`
			// DefaultDeviceIp string `json:"default_device_id" validate:"ip`
		}
		var reqBody body

		e.Validator = &ProductValidator{validator: v}

		if err := c.Bind(&reqBody); err != nil {
			return err
		}

		// if err := v.Struct(reqBody); err != nil {
		// 	return err
		// }

		if err := c.Validate(reqBody); err != nil {
			return err
		}

		product := map[int]string{
			len(products) + 1: reqBody.Name,
		}
		products = append(products, product)
		return c.JSON(http.StatusOK, products)
	})

	e.Logger.Printf("Listning on port: %v", port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%s", port)))
}
