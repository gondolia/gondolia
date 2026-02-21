package service

import "strings"

// formatOptionLabel generates i18n labels from option codes.
// e.g. "1_5kw" -> {"de": "1,5 kW", "en": "1.5 kW"}, "230v" -> {"de": "230V", "en": "230V"}
func formatOptionLabel(code string) map[string]string {
	// Common unit suffixes
	units := map[string]string{
		"kw": " kW", "kva": " kVA", "v": "V", "a": "A",
		"mm": " mm", "cm": " cm", "m": " m", "kg": " kg", "g": " g",
		"l": " L", "ml": " mL", "bar": " bar", "rpm": " rpm",
	}

	lower := strings.ToLower(code)
	for suffix, formatted := range units {
		if strings.HasSuffix(lower, suffix) {
			numPart := code[:len(code)-len(suffix)]
			numPart = strings.ReplaceAll(numPart, "_", ".")
			deLbl := strings.ReplaceAll(numPart, ".", ",") + formatted
			enLbl := numPart + formatted
			return map[string]string{"de": deLbl, "en": enLbl}
		}
	}

	label := strings.ReplaceAll(code, "_", " ")
	if len(label) > 0 {
		label = strings.ToUpper(label[:1]) + label[1:]
	}
	return map[string]string{"de": label, "en": label}
}

// formatAxisLabel generates i18n labels from axis attribute codes.
// e.g. "power_rating" -> {"de": "Leistung", "en": "Power Rating"}
func formatAxisLabel(code string) map[string]string {
	knownAxes := map[string]map[string]string{
		"power_rating":  {"de": "Leistung", "en": "Power Rating"},
		"voltage":       {"de": "Spannung", "en": "Voltage"},
		"mounting_type": {"de": "Bauform", "en": "Mounting Type"},
		"size":          {"de": "Grösse", "en": "Size"},
		"color":         {"de": "Farbe", "en": "Color"},
		"material":      {"de": "Material", "en": "Material"},
		"weight":        {"de": "Gewicht", "en": "Weight"},
		"length":        {"de": "Länge", "en": "Length"},
		"width":         {"de": "Breite", "en": "Width"},
		"height":        {"de": "Höhe", "en": "Height"},
	}

	if labels, ok := knownAxes[code]; ok {
		return labels
	}

	// Fallback: humanize the code
	label := strings.ReplaceAll(code, "_", " ")
	words := strings.Fields(label)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	humanized := strings.Join(words, " ")
	return map[string]string{"de": humanized, "en": humanized}
}

// formatOptionLabelSimple returns a simple display label (German preferred).
func formatOptionLabelSimple(code string) string {
	labels := formatOptionLabel(code)
	if de, ok := labels["de"]; ok {
		return de
	}
	if en, ok := labels["en"]; ok {
		return en
	}
	return code
}
