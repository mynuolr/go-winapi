package winapi

type DLLNAME int

const (
	KERNEL32 DLLNAME = iota
	USER32
)
