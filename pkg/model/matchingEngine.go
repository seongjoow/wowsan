package model

func MatchingEngineSRT(
	srtItem *SubscriptionRoutingTableItem,
	reqSrtItem *SubscriptionRoutingTableItem,
) bool {
	if srtItem.Advertisement.Subject == reqSrtItem.Advertisement.Subject &&
		srtItem.Advertisement.Operator == reqSrtItem.Advertisement.Operator &&
		srtItem.Advertisement.Value == reqSrtItem.Advertisement.Value {
		return true
	}
	return false
}
