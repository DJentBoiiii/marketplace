package internal

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gofiber/fiber/v2"
)

func RenderTemplate(c *fiber.Ctx, templateName string, args ...[2]interface{}) error {

	data := make(map[string]interface{})
	if len(args) > 0 {
		for _, arg := range args {
			key, value := arg[0].(string), arg[1]
			if key != "" && value != nil {
				data[key] = value
			}
		}
	}

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
