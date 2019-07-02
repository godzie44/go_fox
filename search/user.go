package search

type UserInfo struct {
	userId        int
	ipAddressList []string
}

func (u *UserInfo) isDuplicatedWith(anotherUser *UserInfo) bool {
	u2IpAddressMap := make(map[string]interface{})

	for i := range anotherUser.ipAddressList {
		u2IpAddressMap[anotherUser.ipAddressList[i]] = struct{}{}
	}

	notUniqIdCount := 0

	for _, ipAddr := range u.ipAddressList {

		_, ipDuplicated := u2IpAddressMap[ipAddr]

		if ipDuplicated {
			notUniqIdCount++
		}

		if notUniqIdCount >= 2 {
			return true
		}
	}

	return false
}
