package filesystem

import (
	"net/http"
)

var fs = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

// @Summary Отдать файл
// @Description Отдать файл с диска
// @ID get-static
// @Param PathToFile path string true "Путь к файлу"
// @Success 200 "Файл найден"
// @Failure 301 "Редирект, если имя папки не заканчивается на /"
// @Failure 403 "Нет прав (сервер)"
// @Failure 404 "Файл не найден"
// @Failure 500 "Внутренняя ошибка"
// @Router /static/{path/to/file} [GET]
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fs.ServeHTTP(w, r)
}
