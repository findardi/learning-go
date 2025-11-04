package formattime

import (
	"database/sql"
	"time"
)

// Now mengembalikan waktu sekarang dalam UTC sebagai sql.NullTime
func Now() sql.NullTime {
	return sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
}

// NullTime membuat sql.NullTime dari time.Time
func NullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// EmptyNullTime mengembalikan sql.NullTime yang NULL
func EmptyNullTime() sql.NullTime {
	return sql.NullTime{
		Valid: false,
	}
}

// ParseNullTime mengkonversi string ke sql.NullTime
func ParseNullTime(layout, value string) (sql.NullTime, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return EmptyNullTime(), err
	}
	return NullTime(t), nil
}

// FormatNullTime mengkonversi sql.NullTime ke string yang readable
func FormatNullTime(nt sql.NullTime, layout string) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format(layout)
}

// FormatRFC3339 mengkonversi sql.NullTime ke format RFC3339 (ISO 8601)
func FormatRFC3339(nt sql.NullTime) string {
	return FormatNullTime(nt, time.RFC3339)
}

// TimeToString mengkonversi time.Time ke string RFC3339
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}
