package main

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Load sample config included with repo
	//
	// There aren't many surprises in this functionality
	// and so there aren't many edges to worry about

	c, err := LoadConfig()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for _, test := range []struct {
		name   string
		value  interface{}
		expect interface{}
	}{
		{"Correct domain name", c.DomainName, "prattle"},
		{"Correct listen address", c.ListenAddr, "localhost:8080"},
		{"Correct redis", c.RedisAddr, "localhost:6379"},
		{"MaxTokens are correct", c.MaxTokens, 10},
		{"MaxKeys are correct", c.MaxKeys, 10},
		{"Revalidate Frequency is correctly parsed to a time.Duration", c.RevalidateFrequency, time.Second},
		{"There is 1 federated peer", len(c.Federations), 1},
		{"peer[0] has correct ConnectionString", c.Federations["prittle"].ConnectionString, "localhost:8081"},
		{"peer[0] has correct PSK", c.Federations["prittle"].PSK, "F6oqn7CuhrUF0/V2v1+NsO194cQ993t4vrjm4XdqrVS5v+e440ZpQ96hqUCpPmnFvICyfugQc/O1PiPQ4g5tNg=="},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.expect != test.value {
				t.Errorf("expected %v, received %v", test.expect, test.value)
			}
		})
	}
}
