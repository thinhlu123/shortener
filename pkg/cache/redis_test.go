package cache

import (
	"github.com/thinhlu123/shortener/config"
	"log"
	"testing"
)

func beforeTest() *Cache {
	if err := config.GetConfig("../../config/config_local"); err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	var cache Cache
	cache.Init()

	return &cache
}

func TestCache_Get(t *testing.T) {
	cache := beforeTest()

	TestCache_Set(t)

	get, err := cache.Get("a")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if get.(string) != "ad213@!" {
		t.Errorf("Excpect ad213@! but get %v", get)
	}
}

func TestCache_Set(t *testing.T) {
	cache := beforeTest()

	testMap := map[string]string{
		"a": "ad213@!",
		"b": "avagcdq2412",
		"c": "ACNJSndjnqw",
		"d": "ANCSJ123sdn",
		"e": "AAAAA123abc!@",
		"f": "!@#$%",
		"g": "$!@ssfacsaADMJWEQ'",
		"h": "bưdjbajkdsbjkaaskdnakndkandkasndkmasndksnjkdnjkasndjkansdjkbajghvdhwqbjdjkabskjd",
	}

	for key, val := range testMap {
		err := cache.Set(key, val)
		if err != nil {
			t.Error(err)
		}
	}
}
