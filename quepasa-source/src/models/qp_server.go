package models

import "time"

/*
<summary>
	Database representation for whatsapp controller service
</summary>
*/
type QpServer struct {
	// Public token
	Token string `db:"token" json:"token" validate:"max=100"`

	// Whatsapp session id
	WId             string `db:"wid" json:"wid" validate:"max=255"`
	Verified        bool   `db:"verified" json:"verified"`
	Devel           bool   `db:"devel" json:"devel"`
	HandleGroups    bool   `db:"handlegroups" json:"handlegroups,omitempty"`
	HandleBroadcast bool   `db:"handlebroadcast" json:"handlebroadcast,omitempty"`

	User      string    `db:"user" json:"user,omitempty" validate:"max=36"`
	Timestamp time.Time `db:"timestamp" json:"timestamp,omitempty"`
}

func (source *QpServer) GetWId() string {
	return source.WId
}
