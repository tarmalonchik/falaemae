//go:generate go-enum -f=$GOFILE --nocase --sqlnullint
package telegram

// ParseMode ENUM(html, markdown)
type ParseMode int64
