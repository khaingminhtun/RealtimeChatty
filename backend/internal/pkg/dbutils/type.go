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
		// If parsing fails (e.g., empty string or invalid input), return nil
		// so Postgres safely registers it as a NULL database column.
		return nil
	}
	return &addr
}
