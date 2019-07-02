package search

import (
	"testing"
)

type TestCase struct {
	UserInfo1    *UserInfo
	UserInfo2    *UserInfo
	isDuplicated bool
}

func TestDuplicateFinding(t *testing.T) {
	user1 := &UserInfo{1, []string{"127.0.0.1", "127.0.0.2"}}
	user2 := &UserInfo{2, []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}}
	user3 := &UserInfo{3, []string{"127.0.0.3", "127.0.0.1"}}
	user4 := &UserInfo{4, []string{"127.0.0.1"}}

	cases := []TestCase{
		{user1, user2, true},
		{user1, user3, false},
		{user2, user1, true},
		{user3, user2, true},
		{user1, user4, false},
		{user3, user1, false},
		{user1, user1, true},
	}

	for _, item := range cases {

		res := item.UserInfo1.isDuplicatedWith(item.UserInfo2)

		if res != item.isDuplicated {
			t.Errorf("invalid duplicate prediction: %t (must be %t), %#v, %#v", res, item.isDuplicated, item.UserInfo1, item.UserInfo2)
		}

	}

}
