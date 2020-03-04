package netbox

// we need to convert some int64 pointers to nil in case Terraform SDK passed
// value is 0, this due to https://github.com/hashicorp/terraform-plugin-sdk/issues/90
func nilFromInt64Ptr(i *int64) *int64 {
	if *i == int64(0) {
		return nil
	}

	return i
}
