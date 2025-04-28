package catalogue

func formatTypeTitle(productType string) string {
	switch productType {
	case "audio":
		return "Аудіо"
	case "midi":
		return "MIDI"
	case "samples":
		return "Семпли"
	default:
		return productType
	}
}
