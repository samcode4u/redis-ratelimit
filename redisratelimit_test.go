package redisratelimit

import (
	"testing"
	"time"
)

// https://www.programmersought.com/article/7162782496/

func TestGetRoundRobin(t *testing.T) {
	rrl := new(RedisRateLimit)

	rrl.InitClient()

	t.Log("EMPTY")
	rrl.ReSetRoundRobin("RRKEY")
	for i := 0; i < 10; i++ {
		t.Log(rrl.GetRoundRobin("RRKEY"))
	}

	t.Log("ONE")
	rrl.ReSetRoundRobin("RRKEY")
	rrl.SetRoundRobin("RRKEY", "node1")
	for i := 0; i < 10; i++ {
		t.Log(rrl.GetRoundRobin("RRKEY"))
	}

	t.Log("TWO")
	rrl.ReSetRoundRobin("RRKEY")
	rrl.SetRoundRobin("RRKEY", "node1")
	rrl.SetRoundRobin("RRKEY", "node2")
	for i := 0; i < 10; i++ {
		t.Log(rrl.GetRoundRobin("RRKEY"))
	}

	t.Log("SIX")
	rrl.ReSetRoundRobin("RRKEY")
	rrl.SetRoundRobin("RRKEY", "node1")
	rrl.SetRoundRobin("RRKEY", "node2")
	rrl.SetRoundRobin("RRKEY", "node3")
	rrl.SetRoundRobin("RRKEY", "node4")
	rrl.SetRoundRobin("RRKEY", "node5")
	rrl.SetRoundRobin("RRKEY", "node6")
	for i := 0; i < 10; i++ {
		t.Log(rrl.GetRoundRobin("RRKEY"))
	}
}

func TestCheckRateLimit(t *testing.T) {

	rrl := new(RedisRateLimit)

	rrl.InitClient()

	LIMIT := int64(40)

	// reset and start testing
	rrl.ResetRateLimit("key1")

	isAllowed, remaining, ttl := rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT {
		t.Log("TEST1 PASS")
	} else {
		t.Fail()
	}

	rrl.ResetRateLimit("key1")

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT {
		t.Log("TEST2 PASS")
	} else {
		t.Fail()
	}

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, false, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT-1 {
		t.Log("TEST3 PASS")
	} else {
		t.Fail()
	}

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, false, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT-1 {
		t.Log("TEST4 PASS")
	} else {
		t.Fail()
	}

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT-1 {
		t.Log("TEST5 PASS")
	} else {
		t.Fail()
	}

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT-2 {
		t.Log("TEST6 PASS")
	} else {
		t.Fail()
	}

	// reset and start testing
	rrl.ResetRateLimit("key1")

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT {
		t.Log("TEST7 PASS")
	} else {
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	}

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	t.Log(isAllowed, remaining, ttl)

	time.Sleep(time.Second * 1)
	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	t.Log(isAllowed, remaining, ttl)
	time.Sleep(time.Second * 1)
	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	t.Log(isAllowed, remaining, ttl)
	time.Sleep(time.Second * 1)
	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	t.Log(isAllowed, remaining, ttl)
	time.Sleep(time.Second * 1)
	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)
	t.Log(isAllowed, remaining, ttl)
	time.Sleep(time.Second * 1)

	isAllowed, remaining, ttl = rrl.CheckRateLimit("key1", LIMIT, true, time.Second*5)

	t.Log(isAllowed, remaining, ttl)
	if isAllowed && remaining == LIMIT {
		t.Log("TEST8 PASS")
	} else {
		t.Fail()
	}

	// r := int64(LIMIT)

	// for i := int64(0); i < 120; i++ {
	// 	time.Sleep(time.Second)
	// 	isAllowed, remaining = rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)
	// 	t.Log(i)
	// 	if isAllowed && i < LIMIT {
	// 		t.Log("RATELIMIT ALLOWD PASS")
	// 	} else {
	// 		t.Fail()
	// 	}

	// 	if isAllowed == false && i > LIMIT {
	// 		t.Log("RATELIMIT REACHED PASS")
	// 	} else {
	// 		t.Fail()
	// 	}

	// 	if remaining == r && i < LIMIT {
	// 		t.Log("REMAINING MATCHES PASS")
	// 	} else {
	// 		t.Fail()
	// 	}
	// 	r--
	// }

	// time.Sleep(time.Second)

	// isAllowed, remaining = rrl.CheckRateLimit("key1", LIMIT, true, time.Minute)

	// t.Log(remaining, isAllowed)

	// if isAllowed {
	// 	t.Log("RATELIMIT ALLOWD PASS")
	// } else {
	// 	t.Fail()
	// }

	// if remaining == LIMIT {
	// 	t.Log("REMAINING PASS")
	// } else {
	// 	t.Fail()
	// }
}
