package the_redis_one

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"unicode"

	"github.com/hdt3213/rdb/parser"
)

type TheRedisOne struct{}

type Data struct {
	Rdb          string `json:"rdb"`
	Requirements struct {
		CheckTypeOf string `json:"check_type_of"`
	} `json:"requirements"`
}

func (d TheRedisOne) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = makeRdb(data.Rdb)
	if err != nil {
		return nil, err
	}
	rdbFile, err := os.Open("/tmp/redis/dump.rdb")
	if err != nil {
		return nil, err
	}
	defer rdbFile.Close()

	decoder := parser.NewDecoder(rdbFile)
	dbSet := make(map[int]bool)

	err = decoder.Parse(func(o parser.RedisObject) bool {
		switch o.GetType() {
		case parser.StringType:
			str := o.(*parser.StringObject)
			dbSet[str.DB] = true
			if !isASCII(str.Key) {
				output["emoji_key_value"] = string(str.Value)
			}
			if str.Expiration != nil {
				output["expiry_millis"] = str.Expiration.UnixMilli()
			}
		case parser.ListType:
			list := o.(*parser.ListObject)
			dbSet[list.DB] = true
			if list.Key == data.Requirements.CheckTypeOf {
				output[list.Key] = "list"
			}
		case parser.HashType:
			hash := o.(*parser.HashObject)
			dbSet[hash.DB] = true
			if hash.Key == data.Requirements.CheckTypeOf {
				output[hash.Key] = "hash"
			}
		case parser.SetType:
			set := o.(*parser.SetObject)
			dbSet[set.DB] = true
			if set.Key == data.Requirements.CheckTypeOf {
				output[set.Key] = "set"
			}
		}
		return true
	})
	output["db_count"] = len(dbSet)

	return output, err
}

func makeRdb(rdb string) error {
	rawDecodedText, err := base64.StdEncoding.DecodeString(rdb)
	if err != nil {
		return err
	}

	os.Mkdir("/tmp/redis", os.ModePerm)

	f, err := os.Create("/tmp/redis/dump.rdb")
	if err != nil {
		return err
	}

	defer f.Close()
	redisBytes := []byte{0x52, 0x45, 0x44, 0x49, 0x53}
	copy(rawDecodedText[:5], redisBytes)

	_, err = f.Write(rawDecodedText)

	return err
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
