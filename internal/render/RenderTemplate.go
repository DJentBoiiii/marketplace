package render

import (
	"fmt"
	"path"
	"sync"

	"github.com/flosch/pongo2/v6"
	"github.com/gofiber/fiber/v2"
)

type TemplateData map[string]interface{}

var (
	tplCache    = make(map[string]*pongo2.Template)
	tplCacheMtx sync.RWMutex
)

// RenderTemplate рендерить шаблон з переданими даними
func RenderTemplate(c *fiber.Ctx, templateName string, data TemplateData) error {
	// Збереження шаблону в кеш
	tplCacheMtx.RLock()
	tpl, exists := tplCache[templateName]
	tplCacheMtx.RUnlock()

	if !exists {
		// Якщо шаблон не в кеші, загружаємо його
		tplCacheMtx.Lock()
		defer tplCacheMtx.Unlock()

		// Завантажуємо шаблон з файлу
		var err error
		tpl, err = pongo2.FromFile(path.Join("/marketplace/web/static/templates/", templateName))
		if err != nil {
			fmt.Printf("Error loading template %s: %v\n", templateName, err)
			return c.Status(500).SendString("Помилка завантаження шаблону")
		}
		tplCache[templateName] = tpl
	}

	// Виконання шаблону
	out, err := tpl.Execute(pongo2.Context(data))
	if err != nil {
		fmt.Printf("Error rendering template %s: %v\n", templateName, err)
		return c.Status(500).SendString("Помилка рендерингу шаблону")
	}

	// Відправка відповіді
	return c.Type("html", "utf-8").SendString(out)
}
