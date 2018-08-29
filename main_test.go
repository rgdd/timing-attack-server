package main

import (
	"testing"
)

func TestParseArg(t *testing.T) {
	for _, table := range []struct {
		delay, user, tag string
		ok               bool
	}{
		// valid
		{"500", "alice", "12345678", true},
		{"1", "alice", "12345678", true},
		{"1000", "alice", "12345678", true},
		{"+1", "alice", "12345678", true},
		{"500", "a", "12345678", true},
		{"500", "12345678", "12345678", true},
		// bad delay
		{"0", "alice", "12345678", false},
		{"1001", "alice", "12345678", false},
		{"a", "alice", "12345678", false},
		{"1+", "alice", "12345678", false},
		{"1a", "alice", "12345678", false},
		{"0.", "alice", "12345678", false},
		{"0.1", "alice", "12345678", false},
		// bad user
		{"500", "", "12345678", false},
		{"500", "12345678e", "12345678", false},
		// bad tag
		{"500", "alice", "", false},           // too short
		{"500", "alice", "123456", false},     // too short
		{"500", "alice", "1234567", false},    // not even length
		{"500", "alice", "1234567_", false},   // not hex
		{"500", "alice", "1234567890", false}, // too long
	} {
		_, _, _, err := parseArgs(table.delay, table.user, table.tag)
		if table.ok && err != nil {
			t.Errorf("expected %v is ok, got error => %v ", table, err.Error())
		}
		if !table.ok && err == nil {
			t.Errorf("expected %v is bad, got ok", table)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	for _, table := range []struct {
		delay, user, tag string
		ok               bool
	}{
		// valid
		{"1", "alice", "1eda43dd", true},
		// invalid
		{"2", "alice", "1eda43dd", false},
		{"1", "alice", "0eda43dd", false},
		{"1", "alice", "1eda43de", false},
		{"1", "alice", "1edb43dd", false},
	} {
		delay, user, tag, err := parseArgs(table.delay, table.user, table.tag)
		if err != nil {
			t.Fatalf("bad test data => %v", table)
		}

		err = authenticate(delay, user, tag)
		if table.ok && err != nil {
			t.Errorf("expected %v is ok, got error => %v", table, err.Error())
		}
		if table.ok && err != nil {
			t.Errorf("expected %v is ok, got error => %v", table, err.Error())
		}
	}
}
