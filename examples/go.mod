module lewjuh/example

go 1.24.0

replace github.com/lewjuh/diskusage => ../diskusage

require (
	github.com/lewjuh/diskusage v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.33.0
)

require howett.net/plist v1.0.1 // indirect
