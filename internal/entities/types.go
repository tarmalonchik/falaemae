//go:generate go-enum -f=$GOFILE --nocase --sqlnullint
package entities

// DirectionType ENUM(invalid, to_vladikavkaz, to_tskhinvali)
type DirectionType int64
