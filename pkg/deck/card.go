package deck

import (
	"fmt"
	"log"
	"strings"
)

type Color string

// Possible colors
var (
	RED      Color = "r"
	GREEN    Color = "g"
	BLUE     Color = "b"
	YELLOW   Color = "y"
	BLACK    Color = "x"
	CINVALID Color = "-"
)

// Color map from string representation to Color
var Colors = map[string]Color{
	"r": RED, "g": GREEN,
	"b": BLUE, "y": YELLOW,
	"x": BLACK,
}

type CardType uint16

// Possible card types
const (
	NUMBER CardType = (1 << iota)
	DRAW
	REVERSE
	SKIP
	SKIPALL // TODO
	SWAP
	SWAPALL    // TODO
	DISCARDALL // TODO

	WILD
)

func (t CardType) Has(other CardType) bool {
	return t&other != 0
}

func (t CardType) RequiresValue() bool {
	return t.Has(NUMBER | DRAW)
}

func (t CardType) String() string {
	var s = []string{}

	if t.Has(WILD) {
		s = append(s, "wild")
	}

	if t.Has(NUMBER) {
		s = append(s, "number")
	}

	if t.Has(DRAW) {
		s = append(s, "draw")
	}

	if t.Has(REVERSE) {
		s = append(s, "reverse")
	}

	if t.Has(SKIP) {
		s = append(s, "skip")
	}

	if t.Has(SKIPALL) {
		s = append(s, "skipall")
	}

	if t.Has(SWAP) {
		s = append(s, "swap")
	}

	if t.Has(SWAPALL) {
		s = append(s, "swapall")
	}

	return strings.Join(s, "-")
}

// This maps a card .String() representation to its sticker on telegram cache
var STICKER_MAP = map[string]string{
	"b_number_0": "CAACAgEAAxkBAAIBK2DJkaZl4bmgI47DRWr6xkPuR7eHAALUAQACVZ9RRkWe-hVeuGjbHwQ",
	"b_number_1": "CAACAgEAAxkBAAIBLWDJka4R1V6-bf9iLC5oWLj1gE5hAAKeAgAC3_NIRkaHLYay2YL7HwQ",
	"b_number_2": "CAACAgEAAxkBAAIBL2DJkbTyNXKRPYp0ytIzS8VJBdWkAAIpAQAClfhIRoYYku9KdxOWHwQ",
	"b_number_3": "CAACAgEAAxkBAAIBMWDJkbpEdVgieNyH48VhXnesEkN0AAI2AgAC90FRRltSyJcLdJItHwQ",
	"b_number_4": "CAACAgEAAxkBAAIBQWDJkx0t5KgFeWK1NvYFe2TgpoVYAALnAQACV95IRjV20I5uJXkTHwQ",
	"b_number_5": "CAACAgEAAxkBAAIBQ2DJkyalruAvWBjM1HvZXfarviKpAAIdAQACceFQRsfP4aO8fDcvHwQ",
	"b_number_6": "CAACAgEAAxkBAAIBRWDJkyoncULX1pPoWnQLwfSeS3TTAAKtAQACU39RRgW7CQABPQIT6B8E",
	"b_number_7": "CAACAgEAAxkBAAIBR2DJkzCceTC2adApuo-RtZlDhce7AAKQAQAC49hRRunX2WQnbcbGHwQ",
	"b_number_8": "CAACAgEAAxkBAAIBSWDJkza53FS2zEz_boKdghldZcHSAALnAAP7glBGNweu9nakAgcfBA",
	"b_number_9": "CAACAgEAAxkBAAIBS2DJkzu8eOwe9jKZ2rjZjHnX-oZdAAJBAQACR8BJRr-owX5DipOaHwQ",
	"b_draw_2":   "CAACAgEAAxkBAAIBTWDJkz_X-fPvBlkEz47pz7PLRXYcAAIcAgACCMhJRkzDHWZ91l6AHwQ",
	"b_reverse":  "CAACAgEAAxkBAAIBT2DJk0TjZx-dKLEgwKcu9x4-NvsuAAJRAQACzVlIRizomTmZIIZDHwQ",
	"b_skip":     "CAACAgEAAxkBAAIBUWDJk0j-Yffqb5G_tc4CvoZk2VdqAALDAQAC9MhIRt10vM9eHgABQB8E",
	"b_swap":     "CAACAgEAAxkBAAOjYNEPSGpz6YCwKDybWH4LQG6V3sUAAlQBAAJdMYlGcyU99GGajQEfBA",

	"g_number_0": "CAACAgEAAxkBAAIBU2DJk04u41iYCwuVXCmuO9uYydZtAAIuAgACpxxJRoemIok5h-zjHwQ",
	"g_number_1": "CAACAgEAAxkBAAIBVWDJk1TFhr69I9MLbWJF_ANl2AmVAAJ2AQACc1BIRmVGmb9h9PPKHwQ",
	"g_number_2": "CAACAgEAAxkBAAIBV2DJk1jt9qJswZRQ2WSkn9XLYTWCAAJZAQACiQlIRnfuP0nko2BLHwQ",
	"g_number_3": "CAACAgEAAxkBAAIBWWDJk15eMoj93gGQoR5cPH33dHIkAAL-AQACVLFQRjBzDs4m7c0EHwQ",
	"g_number_4": "CAACAgEAAxkBAAIBW2DJk2ImLRZjTM4z74GAkUAcdb3pAAIxAQAC2otRRms00bPG3964HwQ",
	"g_number_5": "CAACAgEAAxkBAAIBXWDJk2dfw4W5mWdaHf1rMX5666jaAAKlAQACqzRJRlj0j9Gkhv-BHwQ",
	"g_number_6": "CAACAgEAAxkBAAIBX2DJk23a_Y37kPNDciuY0zI4balSAAJ4AQAC-5FRRhTSXaRS9BExHwQ",
	"g_number_7": "CAACAgEAAxkBAAIBYWDJk_rIZIGBg0zfDKfPA_3PI95VAAKzAQACWAVQRgY3aAi3bUwbHwQ",
	"g_number_8": "CAACAgEAAxkBAAIBY2DJlAAB1y8axkbuHJT8A2Hl55DLyAACKAEAArTPSUad7df7hmUVSx8E",
	"g_number_9": "CAACAgEAAxkBAAIBZWDJlAZ0lwYK9EyQ5J2SwOISPm41AAIPAQACF7tRRlVlkuxS5nZAHwQ",
	"g_draw_2":   "CAACAgEAAxkBAAIBZ2DJlAsCLArxmM_67olqms11rSueAAJjAQACvN9JRmPqeQ4alEVXHwQ",
	"g_reverse":  "CAACAgEAAxkBAAIBaWDJlA9pLTnBDXgEzoDnpAxsk_vqAAKLAQACbflJRlba1tDRoDVdHwQ",
	"g_skip":     "CAACAgEAAxkBAAIBa2DJlBMOGMXLy4sl6Y4s283ELF2-AAJnAQACrlVRRqWSzLCaD1O7HwQ",
	"g_swap":     "CAACAgEAAxkBAAOnYNEPUkhVdjzt4MiqLr52AAHwVGexAAJHAgACdaWQRlt0_um-jmOlHwQ",

	"r_number_0": "CAACAgEAAxkBAAIBbWDJlBinj6xQOcYsLhxZQ9S0zXw_AALfAQACRqlJRnltaloUAuTnHwQ",
	"r_number_1": "CAACAgEAAxkBAAIBb2DJlB7p0ezsTVg0qajiyK1l9sOeAAINAgACG3xIRnloTzWOMQ_PHwQ",
	"r_number_2": "CAACAgEAAxkBAAIBcWDJlCI2uJxpnhyatl6gA04d8IpfAAKeAQACy-pJRt34wTItDfcwHwQ",
	"r_number_3": "CAACAgEAAxkBAAIBc2DJlCdbYZF6sBmWp1QPELsDQzFVAAJLAQACiaxIRsjf14Pa7JE0HwQ",
	"r_number_4": "CAACAgEAAxkBAAIBdWDJlC1w4RzSuPuvqkNYhG2ebZgfAALKAAN8uVBGun5zFLayrccfBA",
	"r_number_5": "CAACAgEAAxkBAAIBd2DJlDFhH4m4YdTMDKjKHlgEHRNCAAIdAQACTExJRoDDOiaGGiklHwQ",
	"r_number_6": "CAACAgEAAxkBAAIBeWDJlDecY-JSWChdtqHIwGpme7R5AAIKAQACsOlIRhhpDALIv4yCHwQ",
	"r_number_7": "CAACAgEAAxkBAAIBe2DJlDtN3sy49BIbdaNjdzXJ7-BfAAKTAQACemFQRnhb2gyM-DxpHwQ",
	"r_number_8": "CAACAgEAAxkBAAIBfWDJlEAoSyhTaB66gOaFSVpEgRysAAI2AQACoxtIRsT4s7eXkmLOHwQ",
	"r_number_9": "CAACAgEAAxkBAAIBf2DJlERqSzW_eyq43UzwAAFb7iNiFAAC3gEAAi5RSEY1iLeWp29aWh8E",
	"r_draw_2":   "CAACAgEAAxkBAAIBgWDJlEil9RmSEBHcyRSZOqW_xO19AALSAgAClXhJRsfqg1hETSkCHwQ",
	"r_reverse":  "CAACAgEAAxkBAAIBg2DJlEx9xLUjgXprtG70Lv9oKIwJAAJYAQACBuRJRg4HpS_jNAy8HwQ",
	"r_skip":     "CAACAgEAAxkBAAIBhWDJlFRtbejkpH613o7JgCxJUGmNAAJVAQACZWhIRjJ_m3iP-DNCHwQ",
	"r_swap":     "CAACAgEAAxkBAAOpYNEPVFCvdpcQLFqYHVvcNbU--lUAAoMBAAJ_TIlG85v8gye00A0fBA",

	"y_number_0": "CAACAgEAAxkBAAIBh2DJlFueKdFyfrypJ0S8oK1BogKGAAJsAQACQEBJRuuk3JI3IBleHwQ",
	"y_number_1": "CAACAgEAAxkBAAIBiWDJlGDkxS5xFXftifc5jMOaE_KgAAJFAQAC3P1JRi30z7KEqiF-HwQ",
	"y_number_2": "CAACAgEAAxkBAAIBi2DJlGQv6WhssPIAAe9qqSsVmbY4tgACtgEAAtNiSEbnLCA4cEj10x8E",
	"y_number_3": "CAACAgEAAxkBAAIBjWDJlGhz40CHHExO0kBMoKlsCcUEAAKUAQACz0RIRmgv9AoIS6OWHwQ",
	"y_number_4": "CAACAgEAAxkBAAIBj2DJlG4Lku31Lt-7zYLUkShwgMgGAAKUAQACZkRIRiz7LzpEbkegHwQ",
	"y_number_5": "CAACAgEAAxkBAAIBkWDJlHthEtjYuWWg1cPh68yqIkwlAAKSAQACHPFJRmkw8fIbg_6yHwQ",
	"y_number_6": "CAACAgEAAxkBAAIBk2DJlIWgA6AgkJ1rZtIbpe42TJJ5AALxAAMj2khGTCbSLLGiBi4fBA",
	"y_number_7": "CAACAgEAAxkBAAIBlWDJlJOO2PAmDRjz8NCGoHrYZQG9AAKrAQACoIRQRlQy-OC-6krPHwQ",
	"y_number_8": "CAACAgEAAxkBAAIBl2DJlJgp9G2ev3X5nf9wi1d0yh3PAAJoAgACeCBQRm_VWGlbkysRHwQ",
	"y_number_9": "CAACAgEAAxkBAAIBmWDJlJ0fYN1oOmfVyppnq5QoSZ-oAAIwAQACrjBJRnU3tRy_Qs4tHwQ",
	"y_draw_2":   "CAACAgEAAxkBAAIBm2DJlKKai_3QIZ4_uvPXFs3LkgAB-QACJAEAAqvEUEZEIMIQ5r7zDx8E",
	"y_reverse":  "CAACAgEAAxkBAAIBnWDJlKZJuc-JYNd4Os6Nr57-4V_tAAJWAQACjJdIRid1t1zGmr9SHwQ",
	"y_skip":     "CAACAgEAAxkBAAIBn2DJlKzeon9DCvZgqA8RbrCfNFH6AAIzAQACQ6lIRp819aloSOsvHwQ",
	"y_swap":     "CAACAgEAAxkBAAOrYNEPVJOaY9qS9H2ExWiTaindFNwAAoQBAAL5WYlGcazfMc55UekfBA",

	"x_wild":        "CAACAgEAAxkBAAIBoWDJlLADWQRStHlsVSe9-T3TEg0gAAInAQACyi9JRr_nqlKaxdDTHwQ",
	"x_wild-draw_4": "CAACAgEAAxkBAAIBo2DJlLH6hFL4xYJSsrkHI2GJhl2jAAJOAQACP1RJRt3qldF9Fq6VHwQ",
}

// This maps a card .String() representation to its faded sticker on telegram cache
var FADED_STICKER_MAP = map[string]string{
	"b_number_0": "CAACAgEAAxkBAAIBpWDJlL2_QfqBatTVVcbTpA0EYsiiAAJQAQACZSpRRum0YZjn4IVbHwQ",
	"b_number_1": "CAACAgEAAxkBAAIBp2DJlMN-QgZC_hP3qJEV4ktlUTokAAL1AQAC7HBQRuXo-IDujvubHwQ",
	"b_number_2": "CAACAgEAAxkBAAIBqWDJlM1evtX0EZ7U_IxJkkvqH7LHAAIKAgACO9xQRpltHeqh6RLQHwQ",
	"b_number_3": "CAACAgEAAxkBAAIBq2DJlNGR3psx0w3EvFN0QkPMvsBtAALNAQACBulJRndUEKcdSi-8HwQ",
	"b_number_4": "CAACAgEAAxkBAAIBrWDJlNbg4APkA_qzVaOUVzDain7zAAJ9AgAC6slIRgHpd7FAZ1o4HwQ",
	"b_number_5": "CAACAgEAAxkBAAIBr2DJlNvo-FTpHtDqjgielUd7QOCFAAJmAgACeP9JRm9FAf9AdqDzHwQ",
	"b_number_6": "CAACAgEAAxkBAAIBsWDJlOAmn_0-KQq44Ht34UTGVj8HAAJuAQACN0RIRtkqRl-akE3-HwQ",
	"b_number_7": "CAACAgEAAxkBAAIBs2DJlORhKuGtxr3iOYUG_MqOopZ0AAK2AQACC6dJRkVtPGQhjTp9HwQ",
	"b_number_8": "CAACAgEAAxkBAAIBtWDJlOgdqQM44OV4oUr5oPuVhSAaAAIzAQACr_5JRoVigEqyO3C7HwQ",
	"b_number_9": "CAACAgEAAxkBAAIBt2DJlO24DQKj9Tz32viXPzbxMSCCAAJgAQACJeZQRtP6q14ihsbnHwQ",
	"b_draw_2":   "CAACAgEAAxkBAAIBuWDJlPHIQN8k9LV8d2v3CxDvH1xNAAKpAgACOzZRRvbHvnZlll6fHwQ",
	"b_reverse":  "CAACAgEAAxkBAAIBu2DJlPUc_rMHmYRBBYOkzo8h9qtyAALAAQACfkhIRkg2nKZFHMurHwQ",
	"b_skip":     "CAACAgEAAxkBAAIBvWDJlPoJN3QKju94pxcirgdriWkPAAJKAQACFwFQRjjpYXTP9fAGHwQ",
	"b_swap":     "CAACAgEAAxkBAAICoWEE50RVW6k_8sX_osec8WxAnghQAAJ4AQACNqCJRr5U662AZZpqIAQ",

	"g_number_0": "CAACAgEAAxkBAAIBv2DJlP_ngxbrshSmKuk_tFn-aXzzAAJEAgACDWZIRm6-8cDsHGs3HwQ",
	"g_number_1": "CAACAgEAAxkBAAIBwWDJlQS9Hn8xzopyImv0S9ssvC9RAALJAQACCypJRj0MLS52SePTHwQ",
	"g_number_2": "CAACAgEAAxkBAAIBw2DJlQjZ0giWJCCfrvrQEg4QjCzYAAIMAgACnSVIRoL3pSNKQRxLHwQ",
	"g_number_3": "CAACAgEAAxkBAAIBxWDJlQ2nfu7-lv6FiA48YfPxGYWUAAL6AAOXwUhGPtiNDAi4lkIfBA",
	"g_number_4": "CAACAgEAAxkBAAIBx2DJlRItabZ1YpMLeAywQE-x4auyAAJ1AQAC4XdIRgm8y7JFS_n2HwQ",
	"g_number_5": "CAACAgEAAxkBAAIByWDJlRZ-yrCSVJrC774-RDFSYAl8AAItAgAC0lNIRiTNxFicU2LkHwQ",
	"g_number_6": "CAACAgEAAxkBAAIBy2DJlRpYVfOvzGY-nVivfzo4TV4NAAJ2AQACJsxJRuL9I6_y9_eRHwQ",
	"g_number_7": "CAACAgEAAxkBAAIBzWDJlR5PKoHAl3NCdZV2meHZEBX7AAJmAQACm8lQRrzDNxHUTO7sHwQ",
	"g_number_8": "CAACAgEAAxkBAAIBz2DJlSLqzf3-2_SDEcc65RjT6N9UAAKfAQACVI1IRuCuqlemf4GLHwQ",
	"g_number_9": "CAACAgEAAxkBAAIB0WDJlSc5sg5W7IMn_1FeTwFa36zKAAIhAgACSLlIRivuF-FZed4zHwQ",
	"g_draw_2":   "CAACAgEAAxkBAAIB02DJlSuXsShEmNeNQ0RsODUmck-2AAJWAQACwotJRoKjHcy4zJohHwQ",
	"g_reverse":  "CAACAgEAAxkBAAIB1WDJlS-Ms0TG0zCZWdgQ3JUNy0H4AAJuAQACmoJJRpt76k44qpitHwQ",
	"g_skip":     "CAACAgEAAxkBAAIB12DJlTW9lYEtAAHejWdWV1vYhdlmLAACWQEAAjCRSUYl_0l7UiHu2x8E",
	"g_swap":     "CAACAgEAAxkBAAICo2EE50oEyww4psftPH_UGlleLB8MAAKIAwACjROQRmjs9-iA66BgIAQ",

	"r_number_0": "CAACAgEAAxkBAAIB2WDJlTw6EZKkFU-XiD7dWcMFENeMAALfAQACGV1JRtOcHU5WfPfyHwQ",
	"r_number_1": "CAACAgEAAxkBAAIB3WDJlUFcUfG7vVYTnjYCeVnBkkUtAAJmAQACqMRIRq4ofAKjuzu_HwQ",
	"r_number_2": "CAACAgEAAxkBAAIB32DJlUWRnrITp-wHVVlTaVlCkPPSAAL_AQACLrNJRlGL0AttwGgkHwQ",
	"r_number_3": "CAACAgEAAxkBAAIB4WDJlUnAJkwaN8tB6vsDSuL6_Z0nAAKMAQACnlFJRk9tBiyaEbtIHwQ",
	"r_number_4": "CAACAgEAAxkBAAIB42DJlU0nnhF4knsgwU1yl91qbtiPAAJNAQACFO9QRmf3Hu-byD2nHwQ",
	"r_number_5": "CAACAgEAAxkBAAIB5WDJlVEnF_Kn_mCtDdOg_zP_Nd1YAALaAQACbrtIRod0o_6RaFU5HwQ",
	"r_number_6": "CAACAgEAAxkBAAIB52DJlVVoJvbU1U8wt3FP7D3j8IVZAAJAAQACXTpJRpM92LXQmqcxHwQ",
	"r_number_7": "CAACAgEAAxkBAAIB6WDJlVlyK8eKt7le1xfjNDLoc8itAAI7AQACwc5IRhfWHin3mlzOHwQ",
	"r_number_8": "CAACAgEAAxkBAAIB62DJlV44YCYTLa1No0BPzaok1XtZAAKBAQACsjpRRlyfe_RSkWvEHwQ",
	"r_number_9": "CAACAgEAAxkBAAIB7WDJlWJ5aMvVtUcks7EtAAEhRrSRuQAC8wEAAmRMSEZrY9-FrwbAmR8E",
	"r_draw_2":   "CAACAgEAAxkBAAIB72DJlWarXMlu5iFT1VunNjF1yNBYAAIwAQACfnJQRlqAgEOsNHAvHwQ",
	"r_reverse":  "CAACAgEAAxkBAAIB8WDJlWtbEvJZh8NS9NdJsDZiQTCHAAJkAQAC0E5JRv3RAAF-VCTFkx8E",
	"r_skip":     "CAACAgEAAxkBAAIB82DJlW_Kf8KzwE8_AlPy9Rl8YtU1AAJXAQACFr5JRqG2ixkAAQ2HBR8E",
	"r_swap":     "CAACAgEAAxkBAAICpWEE50-ykh26sXFw_mSQXlDt0Tm5AAI_AgACC46JRjphbIH7FlqIIAQ",

	"y_number_0": "CAACAgEAAxkBAAIB-WDJlX611n6YhL3sjF_zZfyUosKoAAJ8AQACAiBJRtX_WfVLl8P7HwQ",
	"y_number_1": "CAACAgEAAxkBAAIB-2DJlYPOKLg9L9GF80FZtf-4-LBYAAKIAQACMfhIRunEX8ou07pIHwQ",
	"y_number_2": "CAACAgEAAxkBAAIB_WDJlYjRp1YBM8kP3fnS_A_FpsKwAAIfAQACzehQRrlTcjcva-pQHwQ",
	"y_number_3": "CAACAgEAAxkBAAIB_2DJlY7A7FAWG5Y6U9KnkjlnosEwAAJNAQACZ8FIRkoGvlIy1vrbHwQ",
	"y_number_4": "CAACAgEAAxkBAAICAWDJlZPJettYeUH53_yQfv7t9NjKAAKdAQACG2tJRrbXPi1zZF28HwQ",
	"y_number_5": "CAACAgEAAxkBAAICA2DJlZir5YbvmJ1Tzdi81Vr1R3qzAAK1AQACiJpIRiGzdQAB7HDGcx8E",
	"y_number_6": "CAACAgEAAxkBAAICB2DJlZ9pwoTaY2NkUlAaLJPS925aAAIvAQACCqhJRh5p3JZb9wthHwQ",
	"y_number_7": "CAACAgEAAxkBAAICCWDJlaNH-vWpzK8aDAo3cfVO8lPkAAI0AQACA1pQRoK2aSf83af-HwQ",
	"y_number_8": "CAACAgEAAxkBAAICC2DJlahhmrY8FXuQco5DEC9AkWIzAAKQAQACcy9IRrz2b-U7o_nMHwQ",
	"y_number_9": "CAACAgEAAxkBAAICDWDJla1-QMlKeFDOZlU9smZw5xzsAAIYAQAC5uJRRq1l4DD9sWktHwQ",
	"y_draw_2":   "CAACAgEAAxkBAAICD2DJla7Ve4____QGyroujc8aCc14AAJ6AQACawNQRgwpXu2oOu60HwQ",
	"y_reverse":  "CAACAgEAAxkBAAICEWDJla-o4TFmY9cCzNOLuDmS_u7xAAKIAQACWaRIRmL7F6AjiAiiHwQ",
	"y_skip":     "CAACAgEAAxkBAAICE2DJlbB0jhy-zVsN3wlzAAFJnULYJQACXAEAAoscSEavb27H3m42pR8E",
	"y_swap":     "CAACAgEAAxkBAAICp2EE51Sv1aXwSqRvTFfJmfRCmX6GAALcAQAC79uQRm4E_cP4BLSXIAQ",

	"x_wild":        "CAACAgEAAxkBAAIB9WDJlXOhd2SeQbf5BZRwnL_B534yAAIhAQACbPBJRgW0z399peF_HwQ",
	"x_wild-draw_4": "CAACAgEAAxkBAAIB92DJlXQ3uZjL0y1E7FG6VZGkuzVUAAJfAgACaJdJRsO9dzLwppwrHwQ",
}

// Color icons for textual representation
var COLOR_ICONS = map[Color]string{
	RED:    "🟥",
	BLUE:   "🟦",
	GREEN:  "🟩",
	YELLOW: "🟨",
	BLACK:  "⬛",
}

// Card represents a single card in a UNO deck
type Card struct {
	Color Color
	Type  CardType
	Value int
}

// NewCard instatiates a new card, it makes sure that only one of value or special are valid values
// Defaults to Black 0
func NewCard(color Color, cardType CardType, value int) *Card {
	if color == CINVALID {
		log.Println("Trying to create a card without a color, defaulting to black")
		color = BLACK
	}

	if cardType.RequiresValue() && value < 0 {
		log.Println("Trying to create a card that requires a value without a valid one, defaulting to 0")
		value = 0
	}

	if !cardType.RequiresValue() && value >= 0 {
		log.Println("Trying to create a card that doesn't require a value with one, setting to -1")
		value = -1
	}

	return &Card{
		Color: color,
		Type:  cardType,
		Value: value,
	}
}

// SetColor changes the card color
// This is used when choosing a special card color for the next player
func (c *Card) SetColor(clr Color) {
	if c.IsSpecial() {
		c.Color = clr
	}
}

// Returns the simplified string representation of the card
func (c *Card) String() string {
	var color = c.Color

	if c.IsSpecial() {
		color = BLACK
	}

	if c.Type.RequiresValue() {
		return fmt.Sprintf("%s_%s_%d", color, c.Type.String(), c.Value)
	}

	return fmt.Sprintf("%s_%s", c.Color, c.Type.String())
}

// Returns the string representation of the card, with color icon and emojis
func (c *Card) StringPretty() string {
	var s = fmt.Sprintf("%s %s", COLOR_ICONS[c.Color], c.Type.String())

	if c.Type.RequiresValue() {
		s += fmt.Sprintf(" %d", c.Value)
	}

	return s
}

// IsSpecial checks if a card is special
func (c *Card) IsSpecial() bool {
	return c.Type.Has(WILD)
}

// GetValue returns the card value
func (c *Card) GetValue() CardType {
	return c.Type
}

// Sticker returns the card's sticker FileID
func (c *Card) Sticker() string {
	s, ok := STICKER_MAP[c.String()]

	if !ok {
		fmt.Printf("Couldn't find card sticker {%s}\n", c)
	}

	return s
}

// StickerNotAvailable returns the card's faded sticker FileID
func (c *Card) StickerNotAvailable() string {
	s, ok := FADED_STICKER_MAP[c.String()]

	if !ok {
		fmt.Printf("Couldn't find card faded sticker {%s}\n", c)
	}

	return s
}

type StackConfig struct {
	CanStackDraws  bool // Draws can be stacked at all, this overrides everything else
	CanStackWild   bool // Wild cards can be stacked
	CanStackBigger bool // Cards with a bigger draw value can also be stacked, otherwise the value must be the same
}

// CanPlayOnTop checks if c can be played on top of c2
// If a draw is pending, only other DRAW cards can be played
// Color must always match (or be a wild card). Other rules depends on config
func (c *Card) CanPlayOnTop(c2 *Card, draw_pending bool, config StackConfig) bool {
	if !draw_pending {
		return c.IsSpecial() || c.Color == c2.Color || (c.Type == c2.Type && c.Value == c2.Value)
	}

	if !config.CanStackDraws || !c.Type.Has(DRAW) {
		return false
	}

	if !config.CanStackWild && c2.Type.Has(WILD) {
		return false
	}

	if config.CanStackBigger {
		return (c.Type.Has(WILD) || c.Color == c2.Color) && c.Value >= c2.Value
	}

	return (c.Type.Has(WILD) || c.Color == c2.Color) && c.Value == c2.Value
}

// UID returns the card Unique Identifier
// We use the cards pointer address here as it's always unique and it
// won't change during a game lifetime as we are moving the same card around Deck, Game and Players' hands
func (c *Card) UID() string {
	return fmt.Sprintf("%p", c)
}

// Score returns the card score value according to this table:
//
// | Card         | Value            |
// | ------------ | ---------------- |
// | Number Cards | Face Value (0-9) |
// | Draw 2       | 20               |
// | Reverse      | 20               |
// | Skip         | 20               |
// | Wild         | 50               |
// | Draw Four    | 50               |
func (c *Card) Score() int {
	return 0

	// Deprecated
	// if c.IsSpecial() {
	// 	return 50
	// }

	// if c.Type > DRAW {
	// 	return 20
	// }

	// return int(c.Type)
}
