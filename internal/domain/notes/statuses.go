package notes

type Status int

const (
	New = iota
	Active
	Inactive
	Deleted
)

var Statuses = []string{"New", "Active", "Inactive", "Deleted"}

func (s Status) String() string {
	return Statuses[s]
}

func ParseStatus(status string) Status {
	switch status {
	case "New":
		return New
	case "Active":
		return Active
	case "Inactive":
		return Inactive
	case "Deleted":
		return Deleted
	default:
		return -1
	}
}
