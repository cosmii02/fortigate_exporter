package probe

func versionAtLeast(meta *TargetMetadata, major, minor int) bool {
	if meta == nil {
		return false
	}
	if meta.VersionMajor > major {
		return true
	}
	if meta.VersionMajor < major {
		return false
	}
	return meta.VersionMinor >= minor
}
