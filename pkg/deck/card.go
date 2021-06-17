package deck

import (
	"fmt"
	"log"
)

type Color string

var (
	RED      Color = "r"
	GREEN    Color = "g"
	BLUE     Color = "b"
	YELLOW   Color = "y"
	BLACK    Color = "x"
	CINVALID Color = "-"
)

var Colors = map[string]Color{
	"r": RED, "g": GREEN,
	"b": BLUE, "y": YELLOW,
	"x": BLACK,
}

type CardValue int

var (
	ZERO     CardValue = 0
	ONE      CardValue = 1
	TWO      CardValue = 2
	THREE    CardValue = 3
	FOUR     CardValue = 4
	FIVE     CardValue = 5
	SIX      CardValue = 6
	SEVEN    CardValue = 7
	EIGHT    CardValue = 8
	NINE     CardValue = 9
	DRAW     CardValue = 10
	REVERSE  CardValue = 11
	SKIP     CardValue = 12
	VINVALID CardValue = -1
)

var CardValues = map[int]CardValue{
	0: ZERO, 1: ONE, 2: TWO, 3: THREE, 4: FOUR, 5: FIVE,
	6: SIX, 7: SEVEN, 8: EIGHT, 9: NINE, 10: DRAW, 11: REVERSE, 12: SKIP,
}

type SpecialCard string

var (
	JOKER    SpecialCard = "joker"
	DFOUR    SpecialCard = "p4"
	SINVALID SpecialCard = "-"
)

var SpecialCards = map[string]SpecialCard{
	"joker": JOKER, "p4": DFOUR,
}

var STICKER_MAP = map[string]string{
	"b_0":  "CAACAgEAAxkBAAIBK2DJkaZl4bmgI47DRWr6xkPuR7eHAALUAQACVZ9RRkWe-hVeuGjbHwQ",
	"b_1":  "CAACAgEAAxkBAAIBLWDJka4R1V6-bf9iLC5oWLj1gE5hAAKeAgAC3_NIRkaHLYay2YL7HwQ",
	"b_2":  "CAACAgEAAxkBAAIBL2DJkbTyNXKRPYp0ytIzS8VJBdWkAAIpAQAClfhIRoYYku9KdxOWHwQ",
	"b_3":  "CAACAgEAAxkBAAIBMWDJkbpEdVgieNyH48VhXnesEkN0AAI2AgAC90FRRltSyJcLdJItHwQ",
	"b_4":  "CAACAgEAAxkBAAIBQWDJkx0t5KgFeWK1NvYFe2TgpoVYAALnAQACV95IRjV20I5uJXkTHwQ",
	"b_5":  "CAACAgEAAxkBAAIBQ2DJkyalruAvWBjM1HvZXfarviKpAAIdAQACceFQRsfP4aO8fDcvHwQ",
	"b_6":  "CAACAgEAAxkBAAIBRWDJkyoncULX1pPoWnQLwfSeS3TTAAKtAQACU39RRgW7CQABPQIT6B8E",
	"b_7":  "CAACAgEAAxkBAAIBR2DJkzCceTC2adApuo-RtZlDhce7AAKQAQAC49hRRunX2WQnbcbGHwQ",
	"b_8":  "CAACAgEAAxkBAAIBSWDJkza53FS2zEz_boKdghldZcHSAALnAAP7glBGNweu9nakAgcfBA",
	"b_9":  "CAACAgEAAxkBAAIBS2DJkzu8eOwe9jKZ2rjZjHnX-oZdAAJBAQACR8BJRr-owX5DipOaHwQ",
	"b_10": "CAACAgEAAxkBAAIBTWDJkz_X-fPvBlkEz47pz7PLRXYcAAIcAgACCMhJRkzDHWZ91l6AHwQ",  // Plus 2
	"b_11": "CAACAgEAAxkBAAIBT2DJk0TjZx-dKLEgwKcu9x4-NvsuAAJRAQACzVlIRizomTmZIIZDHwQ",  // Reverse
	"b_12": "CAACAgEAAxkBAAIBUWDJk0j-Yffqb5G_tc4CvoZk2VdqAALDAQAC9MhIRt10vM9eHgABQB8E", // Skip

	"g_0":  "CAACAgEAAxkBAAIBU2DJk04u41iYCwuVXCmuO9uYydZtAAIuAgACpxxJRoemIok5h-zjHwQ",
	"g_1":  "CAACAgEAAxkBAAIBVWDJk1TFhr69I9MLbWJF_ANl2AmVAAJ2AQACc1BIRmVGmb9h9PPKHwQ",
	"g_2":  "CAACAgEAAxkBAAIBV2DJk1jt9qJswZRQ2WSkn9XLYTWCAAJZAQACiQlIRnfuP0nko2BLHwQ",
	"g_3":  "CAACAgEAAxkBAAIBWWDJk15eMoj93gGQoR5cPH33dHIkAAL-AQACVLFQRjBzDs4m7c0EHwQ",
	"g_4":  "CAACAgEAAxkBAAIBW2DJk2ImLRZjTM4z74GAkUAcdb3pAAIxAQAC2otRRms00bPG3964HwQ",
	"g_5":  "CAACAgEAAxkBAAIBXWDJk2dfw4W5mWdaHf1rMX5666jaAAKlAQACqzRJRlj0j9Gkhv-BHwQ",
	"g_6":  "CAACAgEAAxkBAAIBX2DJk23a_Y37kPNDciuY0zI4balSAAJ4AQAC-5FRRhTSXaRS9BExHwQ",
	"g_7":  "CAACAgEAAxkBAAIBYWDJk_rIZIGBg0zfDKfPA_3PI95VAAKzAQACWAVQRgY3aAi3bUwbHwQ",
	"g_8":  "CAACAgEAAxkBAAIBY2DJlAAB1y8axkbuHJT8A2Hl55DLyAACKAEAArTPSUad7df7hmUVSx8E",
	"g_9":  "CAACAgEAAxkBAAIBZWDJlAZ0lwYK9EyQ5J2SwOISPm41AAIPAQACF7tRRlVlkuxS5nZAHwQ",
	"g_10": "CAACAgEAAxkBAAIBZ2DJlAsCLArxmM_67olqms11rSueAAJjAQACvN9JRmPqeQ4alEVXHwQ",
	"g_11": "CAACAgEAAxkBAAIBaWDJlA9pLTnBDXgEzoDnpAxsk_vqAAKLAQACbflJRlba1tDRoDVdHwQ",
	"g_12": "CAACAgEAAxkBAAIBa2DJlBMOGMXLy4sl6Y4s283ELF2-AAJnAQACrlVRRqWSzLCaD1O7HwQ",

	"r_0":  "CAACAgEAAxkBAAIBbWDJlBinj6xQOcYsLhxZQ9S0zXw_AALfAQACRqlJRnltaloUAuTnHwQ",
	"r_1":  "CAACAgEAAxkBAAIBb2DJlB7p0ezsTVg0qajiyK1l9sOeAAINAgACG3xIRnloTzWOMQ_PHwQ",
	"r_2":  "CAACAgEAAxkBAAIBcWDJlCI2uJxpnhyatl6gA04d8IpfAAKeAQACy-pJRt34wTItDfcwHwQ",
	"r_3":  "CAACAgEAAxkBAAIBc2DJlCdbYZF6sBmWp1QPELsDQzFVAAJLAQACiaxIRsjf14Pa7JE0HwQ",
	"r_4":  "CAACAgEAAxkBAAIBdWDJlC1w4RzSuPuvqkNYhG2ebZgfAALKAAN8uVBGun5zFLayrccfBA",
	"r_5":  "CAACAgEAAxkBAAIBd2DJlDFhH4m4YdTMDKjKHlgEHRNCAAIdAQACTExJRoDDOiaGGiklHwQ",
	"r_6":  "CAACAgEAAxkBAAIBeWDJlDecY-JSWChdtqHIwGpme7R5AAIKAQACsOlIRhhpDALIv4yCHwQ",
	"r_7":  "CAACAgEAAxkBAAIBe2DJlDtN3sy49BIbdaNjdzXJ7-BfAAKTAQACemFQRnhb2gyM-DxpHwQ",
	"r_8":  "CAACAgEAAxkBAAIBfWDJlEAoSyhTaB66gOaFSVpEgRysAAI2AQACoxtIRsT4s7eXkmLOHwQ",
	"r_9":  "CAACAgEAAxkBAAIBf2DJlERqSzW_eyq43UzwAAFb7iNiFAAC3gEAAi5RSEY1iLeWp29aWh8E",
	"r_10": "CAACAgEAAxkBAAIBgWDJlEil9RmSEBHcyRSZOqW_xO19AALSAgAClXhJRsfqg1hETSkCHwQ",
	"r_11": "CAACAgEAAxkBAAIBg2DJlEx9xLUjgXprtG70Lv9oKIwJAAJYAQACBuRJRg4HpS_jNAy8HwQ",
	"r_12": "CAACAgEAAxkBAAIBhWDJlFRtbejkpH613o7JgCxJUGmNAAJVAQACZWhIRjJ_m3iP-DNCHwQ",

	"y_0":  "CAACAgEAAxkBAAIBh2DJlFueKdFyfrypJ0S8oK1BogKGAAJsAQACQEBJRuuk3JI3IBleHwQ",
	"y_1":  "CAACAgEAAxkBAAIBiWDJlGDkxS5xFXftifc5jMOaE_KgAAJFAQAC3P1JRi30z7KEqiF-HwQ",
	"y_2":  "CAACAgEAAxkBAAIBi2DJlGQv6WhssPIAAe9qqSsVmbY4tgACtgEAAtNiSEbnLCA4cEj10x8E",
	"y_3":  "CAACAgEAAxkBAAIBjWDJlGhz40CHHExO0kBMoKlsCcUEAAKUAQACz0RIRmgv9AoIS6OWHwQ",
	"y_4":  "CAACAgEAAxkBAAIBj2DJlG4Lku31Lt-7zYLUkShwgMgGAAKUAQACZkRIRiz7LzpEbkegHwQ",
	"y_5":  "CAACAgEAAxkBAAIBkWDJlHthEtjYuWWg1cPh68yqIkwlAAKSAQACHPFJRmkw8fIbg_6yHwQ",
	"y_6":  "CAACAgEAAxkBAAIBk2DJlIWgA6AgkJ1rZtIbpe42TJJ5AALxAAMj2khGTCbSLLGiBi4fBA",
	"y_7":  "CAACAgEAAxkBAAIBlWDJlJOO2PAmDRjz8NCGoHrYZQG9AAKrAQACoIRQRlQy-OC-6krPHwQ",
	"y_8":  "CAACAgEAAxkBAAIBl2DJlJgp9G2ev3X5nf9wi1d0yh3PAAJoAgACeCBQRm_VWGlbkysRHwQ",
	"y_9":  "CAACAgEAAxkBAAIBmWDJlJ0fYN1oOmfVyppnq5QoSZ-oAAIwAQACrjBJRnU3tRy_Qs4tHwQ",
	"y_10": "CAACAgEAAxkBAAIBm2DJlKKai_3QIZ4_uvPXFs3LkgAB-QACJAEAAqvEUEZEIMIQ5r7zDx8E",
	"y_11": "CAACAgEAAxkBAAIBnWDJlKZJuc-JYNd4Os6Nr57-4V_tAAJWAQACjJdIRid1t1zGmr9SHwQ",
	"y_12": "CAACAgEAAxkBAAIBn2DJlKzeon9DCvZgqA8RbrCfNFH6AAIzAQACQ6lIRp819aloSOsvHwQ",

	"joker": "CAACAgEAAxkBAAIBoWDJlLADWQRStHlsVSe9-T3TEg0gAAInAQACyi9JRr_nqlKaxdDTHwQ",
	"p4":    "CAACAgEAAxkBAAIBo2DJlLH6hFL4xYJSsrkHI2GJhl2jAAJOAQACP1RJRt3qldF9Fq6VHwQ",
}

var FADED_STICKER_MAP = map[string]string{
	"b_0":  "CAACAgEAAxkBAAIBpWDJlL2_QfqBatTVVcbTpA0EYsiiAAJQAQACZSpRRum0YZjn4IVbHwQ",
	"b_1":  "CAACAgEAAxkBAAIBp2DJlMN-QgZC_hP3qJEV4ktlUTokAAL1AQAC7HBQRuXo-IDujvubHwQ",
	"b_2":  "CAACAgEAAxkBAAIBqWDJlM1evtX0EZ7U_IxJkkvqH7LHAAIKAgACO9xQRpltHeqh6RLQHwQ",
	"b_3":  "CAACAgEAAxkBAAIBq2DJlNGR3psx0w3EvFN0QkPMvsBtAALNAQACBulJRndUEKcdSi-8HwQ",
	"b_4":  "CAACAgEAAxkBAAIBrWDJlNbg4APkA_qzVaOUVzDain7zAAJ9AgAC6slIRgHpd7FAZ1o4HwQ",
	"b_5":  "CAACAgEAAxkBAAIBr2DJlNvo-FTpHtDqjgielUd7QOCFAAJmAgACeP9JRm9FAf9AdqDzHwQ",
	"b_6":  "CAACAgEAAxkBAAIBsWDJlOAmn_0-KQq44Ht34UTGVj8HAAJuAQACN0RIRtkqRl-akE3-HwQ",
	"b_7":  "CAACAgEAAxkBAAIBs2DJlORhKuGtxr3iOYUG_MqOopZ0AAK2AQACC6dJRkVtPGQhjTp9HwQ",
	"b_8":  "CAACAgEAAxkBAAIBtWDJlOgdqQM44OV4oUr5oPuVhSAaAAIzAQACr_5JRoVigEqyO3C7HwQ",
	"b_9":  "CAACAgEAAxkBAAIBt2DJlO24DQKj9Tz32viXPzbxMSCCAAJgAQACJeZQRtP6q14ihsbnHwQ",
	"b_10": "CAACAgEAAxkBAAIBuWDJlPHIQN8k9LV8d2v3CxDvH1xNAAKpAgACOzZRRvbHvnZlll6fHwQ",
	"b_11": "CAACAgEAAxkBAAIBu2DJlPUc_rMHmYRBBYOkzo8h9qtyAALAAQACfkhIRkg2nKZFHMurHwQ",
	"b_12": "CAACAgEAAxkBAAIBvWDJlPoJN3QKju94pxcirgdriWkPAAJKAQACFwFQRjjpYXTP9fAGHwQ",

	"g_0":  "CAACAgEAAxkBAAIBv2DJlP_ngxbrshSmKuk_tFn-aXzzAAJEAgACDWZIRm6-8cDsHGs3HwQ",
	"g_1":  "CAACAgEAAxkBAAIBwWDJlQS9Hn8xzopyImv0S9ssvC9RAALJAQACCypJRj0MLS52SePTHwQ",
	"g_2":  "CAACAgEAAxkBAAIBw2DJlQjZ0giWJCCfrvrQEg4QjCzYAAIMAgACnSVIRoL3pSNKQRxLHwQ",
	"g_3":  "CAACAgEAAxkBAAIBxWDJlQ2nfu7-lv6FiA48YfPxGYWUAAL6AAOXwUhGPtiNDAi4lkIfBA",
	"g_4":  "CAACAgEAAxkBAAIBx2DJlRItabZ1YpMLeAywQE-x4auyAAJ1AQAC4XdIRgm8y7JFS_n2HwQ",
	"g_5":  "CAACAgEAAxkBAAIByWDJlRZ-yrCSVJrC774-RDFSYAl8AAItAgAC0lNIRiTNxFicU2LkHwQ",
	"g_6":  "CAACAgEAAxkBAAIBy2DJlRpYVfOvzGY-nVivfzo4TV4NAAJ2AQACJsxJRuL9I6_y9_eRHwQ",
	"g_7":  "CAACAgEAAxkBAAIBzWDJlR5PKoHAl3NCdZV2meHZEBX7AAJmAQACm8lQRrzDNxHUTO7sHwQ",
	"g_8":  "CAACAgEAAxkBAAIBz2DJlSLqzf3-2_SDEcc65RjT6N9UAAKfAQACVI1IRuCuqlemf4GLHwQ",
	"g_9":  "CAACAgEAAxkBAAIB0WDJlSc5sg5W7IMn_1FeTwFa36zKAAIhAgACSLlIRivuF-FZed4zHwQ",
	"g_10": "CAACAgEAAxkBAAIB02DJlSuXsShEmNeNQ0RsODUmck-2AAJWAQACwotJRoKjHcy4zJohHwQ",
	"g_11": "CAACAgEAAxkBAAIB1WDJlS-Ms0TG0zCZWdgQ3JUNy0H4AAJuAQACmoJJRpt76k44qpitHwQ",
	"g_12": "CAACAgEAAxkBAAIB12DJlTW9lYEtAAHejWdWV1vYhdlmLAACWQEAAjCRSUYl_0l7UiHu2x8E",

	"r_0":  "CAACAgEAAxkBAAIB2WDJlTw6EZKkFU-XiD7dWcMFENeMAALfAQACGV1JRtOcHU5WfPfyHwQ",
	"r_1":  "CAACAgEAAxkBAAIB3WDJlUFcUfG7vVYTnjYCeVnBkkUtAAJmAQACqMRIRq4ofAKjuzu_HwQ",
	"r_2":  "CAACAgEAAxkBAAIB32DJlUWRnrITp-wHVVlTaVlCkPPSAAL_AQACLrNJRlGL0AttwGgkHwQ",
	"r_3":  "CAACAgEAAxkBAAIB4WDJlUnAJkwaN8tB6vsDSuL6_Z0nAAKMAQACnlFJRk9tBiyaEbtIHwQ",
	"r_4":  "CAACAgEAAxkBAAIB42DJlU0nnhF4knsgwU1yl91qbtiPAAJNAQACFO9QRmf3Hu-byD2nHwQ",
	"r_5":  "CAACAgEAAxkBAAIB5WDJlVEnF_Kn_mCtDdOg_zP_Nd1YAALaAQACbrtIRod0o_6RaFU5HwQ",
	"r_6":  "CAACAgEAAxkBAAIB52DJlVVoJvbU1U8wt3FP7D3j8IVZAAJAAQACXTpJRpM92LXQmqcxHwQ",
	"r_7":  "CAACAgEAAxkBAAIB6WDJlVlyK8eKt7le1xfjNDLoc8itAAI7AQACwc5IRhfWHin3mlzOHwQ",
	"r_8":  "CAACAgEAAxkBAAIB62DJlV44YCYTLa1No0BPzaok1XtZAAKBAQACsjpRRlyfe_RSkWvEHwQ",
	"r_9":  "CAACAgEAAxkBAAIB7WDJlWJ5aMvVtUcks7EtAAEhRrSRuQAC8wEAAmRMSEZrY9-FrwbAmR8E",
	"r_10": "CAACAgEAAxkBAAIB72DJlWarXMlu5iFT1VunNjF1yNBYAAIwAQACfnJQRlqAgEOsNHAvHwQ",
	"r_11": "CAACAgEAAxkBAAIB8WDJlWtbEvJZh8NS9NdJsDZiQTCHAAJkAQAC0E5JRv3RAAF-VCTFkx8E",
	"r_12": "CAACAgEAAxkBAAIB82DJlW_Kf8KzwE8_AlPy9Rl8YtU1AAJXAQACFr5JRqG2ixkAAQ2HBR8E",

	"y_0":  "CAACAgEAAxkBAAIB-WDJlX611n6YhL3sjF_zZfyUosKoAAJ8AQACAiBJRtX_WfVLl8P7HwQ",
	"y_1":  "CAACAgEAAxkBAAIB-2DJlYPOKLg9L9GF80FZtf-4-LBYAAKIAQACMfhIRunEX8ou07pIHwQ",
	"y_2":  "CAACAgEAAxkBAAIB_WDJlYjRp1YBM8kP3fnS_A_FpsKwAAIfAQACzehQRrlTcjcva-pQHwQ",
	"y_3":  "CAACAgEAAxkBAAIB_2DJlY7A7FAWG5Y6U9KnkjlnosEwAAJNAQACZ8FIRkoGvlIy1vrbHwQ",
	"y_4":  "CAACAgEAAxkBAAICAWDJlZPJettYeUH53_yQfv7t9NjKAAKdAQACG2tJRrbXPi1zZF28HwQ",
	"y_5":  "CAACAgEAAxkBAAICA2DJlZir5YbvmJ1Tzdi81Vr1R3qzAAK1AQACiJpIRiGzdQAB7HDGcx8E",
	"y_6":  "CAACAgEAAxkBAAICB2DJlZ9pwoTaY2NkUlAaLJPS925aAAIvAQACCqhJRh5p3JZb9wthHwQ",
	"y_7":  "CAACAgEAAxkBAAICCWDJlaNH-vWpzK8aDAo3cfVO8lPkAAI0AQACA1pQRoK2aSf83af-HwQ",
	"y_8":  "CAACAgEAAxkBAAICC2DJlahhmrY8FXuQco5DEC9AkWIzAAKQAQACcy9IRrz2b-U7o_nMHwQ",
	"y_9":  "CAACAgEAAxkBAAICDWDJla1-QMlKeFDOZlU9smZw5xzsAAIYAQAC5uJRRq1l4DD9sWktHwQ",
	"y_10": "CAACAgEAAxkBAAICD2DJla7Ve4____QGyroujc8aCc14AAJ6AQACawNQRgwpXu2oOu60HwQ",
	"y_11": "CAACAgEAAxkBAAICEWDJla-o4TFmY9cCzNOLuDmS_u7xAAKIAQACWaRIRmL7F6AjiAiiHwQ",
	"y_12": "CAACAgEAAxkBAAICE2DJlbB0jhy-zVsN3wlzAAFJnULYJQACXAEAAoscSEavb27H3m42pR8E",

	"joker": "CAACAgEAAxkBAAIB9WDJlXOhd2SeQbf5BZRwnL_B534yAAIhAQACbPBJRgW0z399peF_HwQ",
	"p4":    "CAACAgEAAxkBAAIB92DJlXQ3uZjL0y1E7FG6VZGkuzVUAAJfAgACaJdJRsO9dzLwppwrHwQ",
}

var COLOR_ICONS = map[Color]string{
	RED:    "🟥",
	BLUE:   "🟦",
	GREEN:  "🟩",
	YELLOW: "🟨",
	BLACK:  "⬛",
}

type Card struct {
	Color   Color
	Value   CardValue
	Special SpecialCard
}

func NewCard(color Color, value CardValue, special SpecialCard) Card {
	if color == CINVALID {
		log.Println("Trying to create a card without a color, defaulting to black")
		color = BLACK
	}

	if value != VINVALID && special != SINVALID {
		log.Println("Trying to create a card with both a value and special, value will be discarded")
		value = VINVALID
	}

	if value == VINVALID && special == SINVALID {
		log.Println("Trying to create a card with neither a value nor special, value will default to 0")
		value = ZERO
	}

	return Card{
		Color:   color,
		Value:   value,
		Special: special,
	}
}

func (c *Card) SetColor(clr Color) {
	if c.IsSpecial() {
		c.Color = clr
	}
}

func (c *Card) String() string {
	if c.IsSpecial() {
		return string(c.Special)
	}

	return fmt.Sprintf("%s_%d", c.Color, c.Value)
}

func (c *Card) StringPretty() string {
	if c.IsSpecial() {
		return fmt.Sprintf("%s %s", COLOR_ICONS[c.Color], c.Special)
	}

	if c.Value < DRAW {
		return fmt.Sprintf("%s %d", COLOR_ICONS[c.Color], c.Value)
	}

	if c.Value == DRAW {
		return fmt.Sprintf("%s +2", COLOR_ICONS[c.Color])
	}

	if c.Value == REVERSE {
		return fmt.Sprintf("%s 🔃", COLOR_ICONS[c.Color])
	}

	return fmt.Sprintf("%s 🚫", COLOR_ICONS[c.Color])
}

func (c *Card) IsSpecial() bool {
	return c.Special != SINVALID
}

func (c *Card) GetValue() CardValue {
	return c.Value
}

func (c *Card) GetSpecial() SpecialCard {
	return c.Special
}

func (c *Card) Sticker() string {
	s, ok := STICKER_MAP[c.String()]

	if !ok {
		fmt.Printf("Couldn't find card sticker {%s}\n", c)
	}

	return s
}

func (c *Card) StickerNotAvailable() string {
	s, ok := FADED_STICKER_MAP[c.String()]

	if !ok {
		fmt.Printf("Couldn't find card faded sticker {%s}\n", c)
	}

	return s
}

func (c *Card) CanPlayOnTop(c2 *Card, draw_pending bool) bool {
	if draw_pending {
		if c2.Value == DRAW {
			return c.Value == DRAW
		}

		return false
	}

	return c.IsSpecial() || c.Value == c2.Value || c.Color == c2.Color
}

func (c *Card) UID() string {
	// Each card has a unique address so we can use it as a unique identifier
	return fmt.Sprintf("%p", c)
}

func (c *Card) Score() int {
	if c.IsSpecial() {
		return 50
	}

	if c.Value > DRAW {
		return 20
	}

	return int(c.Value)
}
