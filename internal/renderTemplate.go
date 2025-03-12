package internal

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gofiber/fiber/v2"
)

func RenderTemplate(c *fiber.Ctx, templateName string, data map[string]interface{}) error {
	tpl, err := pongo2.FromFile("/marketplace/web/static/templates/" + templateName)
	if err != nil {
		return c.Status(500).SendString("Помилка завантаження шаблону")
	}

	out, err := tpl.Execute(pongo2.Context(data))
	if err != nil {
		return c.Status(500).SendString("Помилка рендерингу шаблону")
	}

	return c.Type("html", "utf-8").SendString(out)
}
