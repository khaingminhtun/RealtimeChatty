package dbutils

import (
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func NewTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func ParseIPAddress(ipStr string) *netip.Addr {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil
	}
	return &addr
}

// NewDate converts a time.Time object directly into a valid pgtype.Date
func NewDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

// ParseDateString takes a "YYYY-MM-DD" string and converts it to pgtype.Date.
// If the string is empty or invalid, it gracefully returns an invalid (NULL) pgtype.Date.
func ParseDateString(dateStr string) pgtype.Date {
	if dateStr == "" {
		return pgtype.Date{Valid: false}
	}

	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  parsedTime,
		Valid: true,
	}
}
