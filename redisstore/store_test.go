// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package redisstore

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
)

var testClient *redis.Client

func TestMain(m *testing.M) {
	testClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := testClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestPrefixOption(t *testing.T) {
	s := &store{}
	prefix := "foo"
	Prefix(prefix)(s)
	if s.prefix != prefix {
		t.Errorf("expected prefix %s, got %s", prefix, s.prefix)
	}
}
func TestGetKey(t *testing.T) {
	prefix := "foo"
	s := &store{prefix: prefix}
	key := "bar"
	if s.getKey(key) != prefix+":"+key {
		t.Errorf("expected key %s, got %s", prefix+":"+key, s.getKey(key))
	}
}

func TestNew(t *testing.T) {
	expiration := 10 * time.Minute
	prefix := "foo"
	s, _ := New(testClient, Expiration(expiration), Prefix(prefix)).(*store)
	if s.expiration != expiration {
		t.Errorf("expected expiration %v, got %v", expiration, s.expiration)
	}
	if s.prefix != prefix {
		t.Errorf("expected prefix %s, got %s", prefix, s.prefix)
	}
}

func TestStoreGet(t *testing.T) {
	s := New(testClient)
	_, err := s.Get("foo", true)
	if err == nil {
		t.Error("expected a non-nil error, got nil")
	}

	err = s.Set("foo", "bar")
	if err != nil {
		t.Fatalf("failed to set: %s", err)
	}
	for _, clear := range []bool{false, true} {
		value, err := s.Get("foo", clear)
		if err != nil {
			t.Fatalf("expected non error, got %s", err)
		}
		if value != "bar" {
			t.Errorf("expected value %q, got %q", "bar", value)
		}
	}

	_, err = s.Get("foo", true)
	if err == nil {
		t.Error("expected a non-nil error, got nil")
	}
}

func TestStoreSet(t *testing.T) {
	// TBD
}
