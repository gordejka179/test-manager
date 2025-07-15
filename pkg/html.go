package pkg

import (
	"fmt"
	"strings"
)

func GenerateForm(config map[string]interface{}, prefix string, w *strings.Builder) {
	for key, value := range config {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}: // Вложенная структура
			w.WriteString(fmt.Sprintf("<fieldset><legend>%s</legend>\n", key))
			GenerateForm(v, fullKey, w)
			w.WriteString("</fieldset>\n")

		default: // Обычное поле
			w.WriteString(fmt.Sprintf(
				"<label for='%s'>%s</label>\n"+
					"<input type='text' id='%s' name='%s' value='%v'><br>\n",
				fullKey, key, fullKey, fullKey, v,
			))
		}
	}
}

func CreateHTMLForm(config map[string]interface{}) string {
	var html strings.Builder
	html.WriteString("<form method='POST'>\n")

	GenerateForm(config, "", &html)

	html.WriteString("<input type='submit' value='Save'>\n</form>")
	return html.String()
}
