package locales

import (
	"encoding/json"
	"go-echo/helpers"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizer *i18n.Localizer
var bundle *i18n.Bundle

func LocalesGenerate() *i18n.Bundle {
	// Load i18n
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("locales/en.json")
	bundle.LoadMessageFile("locales/id.json")

	return bundle
}

func LocalesSet() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check request header for accepting language for the app
			acceptLanguage := c.Request().Header["Accept-Language"]
			if len(acceptLanguage) > 0 {
				lang := acceptLanguage[0]
				if acceptLanguage[0] == "" {
					return helpers.ErrorServer("Please fill Accept-Language header value")
				}

				localizer = i18n.NewLocalizer(bundle, lang)

				return next(c)
			} else {
				return helpers.ErrorServer("You must add Accept-Language header")
			}
		}
	}
}

func LocalesGet(key string) string {
	localizeConfig := i18n.LocalizeConfig{
		MessageID: key,
	}
	localization, _ := localizer.Localize(&localizeConfig)

	return localization
}
