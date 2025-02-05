// Copyright (C) 2022 Satoshi Konno All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"time"

	"github.com/trueblacker/go-redis/redis"
	"github.com/trueblacker/go-redis/redis/regexp"
)

func (server *Server) Del(conn *redis.Conn, keys []string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	removedCount := 0
	for _, key := range keys {
		err := db.RemoveRecord(key)
		if err == nil {
			removedCount++
		}
	}
	return redis.NewIntegerMessage(removedCount), nil
}

func (server *Server) Exists(conn *redis.Conn, keys []string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	existCount := 0
	for _, key := range keys {
		_, ok := db.GetRecord(key)
		if ok {
			existCount++
		}
	}
	return redis.NewIntegerMessage(existCount), nil
}

func (server *Server) Expire(conn *redis.Conn, key string, opt redis.ExpireOption) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return redis.NewIntegerMessage(0), nil
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewIntegerMessage(0), nil
	}
	now := time.Now()
	record.TTL = opt.Time.Sub(now)
	return redis.NewIntegerMessage(1), nil
}

func (server *Server) Type(conn *redis.Conn, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewStringMessage("none"), nil
	}
	switch record.Data.(type) {
	case string:
		return redis.NewStringMessage("string"), nil
	case Hash:
		return redis.NewStringMessage("hash"), nil
	case *List:
		return redis.NewStringMessage("list"), nil
	case *Set:
		return redis.NewStringMessage("set"), nil
	case *ZSet:
		return redis.NewStringMessage("zset"), nil
	}
	return redis.NewStringMessage("none"), nil
}

func (server *Server) Keys(conn *redis.Conn, pattern string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	r := regexp.MustCompile(pattern)
	matchKeys := []string{}
	for _, key := range db.Keys() {
		if !r.MatchString(key) {
			continue
		}
		matchKeys = append(matchKeys, key)
	}
	return redis.NewStringArrayMessage(matchKeys), nil
}

func (server *Server) Rename(conn *redis.Conn, key string, newkey string, opt redis.RenameOption) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	if opt.NX {
		if _, ok := db.GetRecord(newkey); ok {
			return redis.NewIntegerMessage(0), nil
		}
	}
	err = db.RenameRecord(key, newkey)
	if err != nil {
		return nil, err
	}
	if opt.NX {
		return redis.NewIntegerMessage(1), nil
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) TTL(conn *redis.Conn, key string) (*redis.Message, error) {
	const ttlRecordNotFound int = -2
	const ttlRecordNotSet int = -1
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewIntegerMessage(ttlRecordNotFound), nil
	}
	if record.TTL <= 0 {
		return redis.NewIntegerMessage(ttlRecordNotSet), nil
	}
	now := time.Now()
	ttl := record.Timestamp.Add(record.TTL).Sub(now)
	if ttl < 0 {
		return redis.NewIntegerMessage(ttlRecordNotFound), nil
	}
	return redis.NewIntegerMessage(int(ttl / time.Second)), nil
}
