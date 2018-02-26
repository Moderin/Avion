package utilites

import "strings"

// Replace < > &  with its HTML names
func ReplaceHTMLSymbols(data *string) {
    *data = strings.Replace(*data, "&", "&amp;", -1)
	*data = strings.Replace(*data, "<", "&lt;", -1)
	*data = strings.Replace(*data, ">", "&gt;", -1)
}