package authconfig

import (
	"flag"
)

// AddFlags adds flags applicable to all services.
// Remember to call `flag.Parse()` in your main or TestMain.
func AddFlags(fs flag.FlagSet) error {
	fs.StringVar(&subscriptionID, "subscription", subscriptionID, "Subscription for tests.")
	fs.StringVar(&locationDefault, "location", locationDefault, "Default location for tests.")
	fs.StringVar(&cloudName, "cloud", cloudName, "Name of Azure cloud.")
	fs.StringVar(&baseGroupName, "baseGroupName", BaseGroupName(), "Specify prefix name of resource group for sample resources.")

	fs.BoolVar(&useDeviceFlow, "useDeviceFlow", useDeviceFlow, "Use device-flow grant type rather than client credentials.")
	fs.BoolVar(&keepResources, "keepResources", keepResources, "Keep resources created by samples.")

	return nil
}
