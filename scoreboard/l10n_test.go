/**
 * @file l10n_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 * @brief test localization helpers
 */

package scoreboard

import (
	"net/http/httptest"
	"testing"
)

func TestL10n(*testing.T) {
	r := httptest.NewRequest("GET", "http://localhost", nil)

	// Must be translate
	r.Header = map[string][]string{
		"Accept-Language": {"ru-ru"},
	}

	for key, value := range l10nMap {
		if l10n(r, key) != value {
			panic("Wrong l10n")
		}
	}

	for _, value := range l10nMap {
		if l10n(r, "RANDOM_STRING_c504fe0cc6abb24") == value {
			panic("Wrong l10n")
		}
	}

	// Must not translate
	r.Header = map[string][]string{
		"Accept-Language": {"en-Us"},
	}

	for key, value := range l10nMap {
		if l10n(r, key) == value {
			panic("Wrong l10n")
		}
	}

	for _, value := range l10nMap {
		if l10n(r, value) != value {
			panic("Wrong l10n")
		}
	}

}
