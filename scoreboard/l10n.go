/**
 * @file l10n.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 * @brief localization helpers
 */

package scoreboard

import (
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

var l10nMap = map[string]string{
	"contest":                "турнир",
	contestStateNotAvailable: "состояние неизвестно",
	contestNotStarted:        "остановлен",
	contestRunning:           "запущен",
	contestCompleted:         "завершен",

	"Scoreboard": "Турнирная таблица",
	"Tasks":      "Задачи",
	"News":       "Новости",
	"Sponsors":   "Спонсоры",

	"<th>Team</th>":  "<th>Команда</th>",
	"<th>Score</th>": "<th>Счет</th>",

	"Access token": "Токен доступа",
	"Sign in":      "Войти",

	"Solved":       "Флаг принят",
	"Invalid flag": "Неправильный флаг",

	`btn-submit">Submit</button`: `btn-submit">Отправить</button`,
	`placeholder="Flag"`:         `placeholder="Флаг"`,
}

var supported = []language.Tag{
	language.AmericanEnglish,
	language.Russian,
}

func getLanguage(r *http.Request) (lang language.Tag) {
	h := r.Header.Get("Accept-Language")
	t, _, err := language.ParseAcceptLanguage(h)
	if err != nil {
		return supported[0]
	}

	m := language.NewMatcher(supported)
	lang, _, _ = m.Match(t...)

	return
}

func isAcceptRussian(r *http.Request) bool {
	return getLanguage(r) == language.Russian
}

func toRussian(in string) (out string) {
	out = in
	for key, value := range l10nMap {
		out = strings.Replace(out, key, value, -1)
	}
	return
}

func l10n(r *http.Request, html string) (s string) {
	if isAcceptRussian(r) {
		return toRussian(html)
	}

	return html
}
