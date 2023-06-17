package main

import (
	"fmt"
	"strings"

	"github.com/Reidond/vue-template-compiler-wasm/internal/compiler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ParseBody is helper function for parsing the body.
// Is any error occurs it will panic.
// Its just a helper function to avoid writing if condition again n again.
func parseBody(ctx *fiber.Ctx, body interface{}) *fiber.Error {
	if err := ctx.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	return nil
}

var validate = validator.New()

// Validate validates the input struct
func validatePayload(payload interface{}) *fiber.Error {
	err := validate.Struct(payload)

	if err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				fmt.Sprintf("`%v` with value `%v` doesn't satisfy the `%v` constraint", err.Field(), err.Value(), err.Tag()),
			)
		}

		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: strings.Join(errors, ","),
		}
	}

	return nil
}

// ParseBodyAndValidate is helper function for parsing the body.
// Is any error occurs it will panic.
// Its just a helper function to avoid writing if condition again n again.
func parseBodyAndValidate(ctx *fiber.Ctx, body interface{}) *fiber.Error {
	if err := parseBody(ctx, body); err != nil {
		return err
	}

	return validatePayload(body)
}

type SfcCodeRequest struct {
	MountID string `json:"mountId" validate:"required"`
	SfcCode string `json:"sfcCode" validate:"required"`
}

type SfcCodeResponse struct {
	AppCode   string   `json:"appCode"`
	StyleCode []string `json:"styleCode"`
}

func main() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		b := new(SfcCodeRequest)

		if err := parseBodyAndValidate(c, b); err != nil {
			return err
		}

		out, err := compiler.CompileSfcCode(b.SfcCode, b.MountID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(&SfcCodeResponse{AppCode: out.AppCode, StyleCode: out.StyleCode})
	})

	app.Listen(":3001")
}
