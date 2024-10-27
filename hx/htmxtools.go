package hx

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func IsHxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func IsPartialUpdate(c echo.Context) bool {
	return IsHxRequest(c) && c.Request().Header.Get("HX-Boosted") != "true"
}

func Redirect(c echo.Context, location string) error {
	if IsHxRequest(c) {
		c.Response().Header().Set("HX-Redirect", string(templ.URL(location)))
		return c.NoContent(http.StatusOK)
	}
	return c.Redirect(http.StatusSeeOther, location)
}

func Trigger(c echo.Context, triggers ...string) {
	c.Response().Header().Set("HX-Trigger", strings.Join(triggers, ", "))
}

func TriggerAfterSwap(c echo.Context, triggers ...string) {
	c.Response().Header().Set("HX-Trigger-After-Swap", strings.Join(triggers, ", "))
}

func ReplaceUrl(c echo.Context, url string) {
	c.Response().Header().Set("HX-Replace-Url", string(templ.URL(url)))
}

func PushUrl(c echo.Context, url string) {
	c.Response().Header().Set("HX-Push-Url", string(templ.URL(url)))
}

func Retarget(c echo.Context, selector string) {
	c.Response().Header().Set("HX-Retarget", selector)
}

func Reswap(c echo.Context, swapMode string) {
	c.Response().Header().Set("HX-Reswap", swapMode)
}

func GetHxTrigger(c echo.Context) string {
	return c.Request().Header.Get("Hx-Trigger")
}

func GetCurrentUrl(c echo.Context) *url.URL {
	hxCurrentUrl := c.Request().Header.Get("HX-Current-URL")
	u, err := url.Parse(hxCurrentUrl)
	if err != nil {
		return nil
	}
	return u
}
