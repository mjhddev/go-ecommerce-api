package cache

import (
	"encoding/json"
	"time"

	"github.com/mjhddev/go-ecommerce-api/configs"
)

func Get(key string, dest interface{}) (bool, error) {
	val, err := configs.Redis.Get(configs.Ctx, key).Result()
	if err != nil {
		return false, nil
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil
}

func Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return configs.Redis.Set(
		configs.Ctx,
		key,
		data,
		ttl,
	).Err()
}

func Delete(key string) error {
	return configs.Redis.Del(
		configs.Ctx,
		key,
	).Err()
}

func DeleteByPattern(pattern string) error {
	keys, err := configs.Redis.Keys(configs.Ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return configs.Redis.Del(configs.Ctx, keys...).Err()
}
