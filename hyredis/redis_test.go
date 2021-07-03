package hyredis

import (
	"github.com/y1015860449/go-tools/hy-utils"
	"reflect"
	"testing"
	"time"
)

var (
	rConf = RedisConfig{
		Addrs:        "127.0.0.1:6379",
		MaxIdleConns: 1024,
		MaxOpenConns: 0,
		MaxLifeTime:  100,
	}
)

func Test_Get(t *testing.T) {

	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		key         string
		beforeStart func(key string) error
		onEnd       func(key string) error
		want        string
		wantErr     bool
	}{
		//
		{"getNotExitKey", randomStr(), nil, nil, "", true},
		//
		{"getExitKey", randomStr(), func(key string) error {
			return cli.Set(key, "valid value")
		}, func(key string) error {
			return cli.Del(key)
		}, "valid value", false},
		//
		{"getExitKeyWithEmptyValue", randomStr(), func(key string) error {
			return cli.Set(key, "")
		}, func(key string) error {
			return cli.Del(key)
		}, "", false},
		//
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeStart != nil {
				if err = tt.beforeStart(tt.key); err != nil {
					t.Errorf("Get()  run beforeStart err %v", err)
				}
			}
			got, err := cli.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if tt.onEnd != nil {
				if err = tt.onEnd(tt.key); err != nil {
					t.Errorf("Get()  run onEnd err %v", err)
				}
			}
		})
	}
}

func Test_MGet(t *testing.T) {

	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		keys        []string
		beforeStart func(keys []string) error
		onEnd       func(keys []string) error
		wants       []string
		wantErr     bool
	}{
		{"getNotExitKey", randomStrs(3), nil, nil, nil, false},
		// key1 and key3 exist, key3 does not exist, expect get []{"valid key","","valid key}
		{"getOnKeyAndOneInvalidKey", randomStrs(3), func(keys []string) error {
			clis := make(map[string]string)
			clis[keys[0]] = "valid key"
			// clis[keys[1]] // not exist
			clis[keys[2]] = "valid key"
			return cli.MSet(clis)
		}, func(keys []string) error {
			return cli.Del(keys...)
		}, []string{"valid key", "", "valid key"}, false},
		//
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeStart != nil {
				if err = tt.beforeStart(tt.keys); err != nil {
					t.Errorf("MGet()  run beforeStart err %v", err)
				}
			}
			gots, err := cli.MGet(tt.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && gots != nil && tt.wants != nil {
				if len(tt.wants) != len(gots) {
					t.Errorf("MGet get len %d, want len %d", len(gots), len(tt.wants))
				} else {
					for i, want := range tt.wants {
						if want != gots[i] {
							t.Errorf("MGet get %v, want %v", gots[i], want)
						}
					}
				}
			}
			if tt.onEnd != nil {
				if err = tt.onEnd(tt.keys); err != nil {
					t.Errorf("MGet()  run onEnd err %v", err)
				}
			}
		})
	}
}

func randomStr() string {
	return hy_utils.GetUUID()[:4]
}

func randomStrs(num int) []string {
	var keys []string
	for i := 0; i < num; i++ {
		keys = append(keys, randomStr())
	}
	return keys
}

func Test_Del(t *testing.T) {

	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}

	exitKey := randomStr()
	notExitKey := randomStr()

	// prepare test data
	if err = cli.SetEx(exitKey, "valid key", 10); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"deleteExistKey", []string{exitKey}, false},
		{"deleteNotExistKey", []string{notExitKey}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.Del(tt.args...); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Exists(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	exitKey := randomStr()
	notExitKey := randomStr()

	defer cli.Del(exitKey, notExitKey)

	// 设置测试值
	if err = cli.SetEx(exitKey, "valid key", 10); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    string
		want    bool
		wantErr bool
	}{
		{"exitKey", exitKey, true, false},
		{"notExitKey", notExitKey, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.Exists(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Expire(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}

	// non-existent key, return false
	if ok, err := cli.Expire(randomStr(), 1); err != nil || ok {
		t.Error("expect false")
	}
	// prepare test data
	exitKey := randomStr()
	exitValue := "exit value"
	defer cli.Del(exitKey)
	if err = cli.Set(exitKey, exitValue); err != nil {
		t.Fatal("expect success")
	}

	timeout := 2

	// want set success
	if ok, err := cli.Expire(exitKey, timeout); err != nil || !ok {
		t.Error("expect true")
	}

	// want to get the correct data
	if v, err := cli.Get(exitKey); err != nil || v != exitValue {
		t.Errorf("expect get %s , but get %v", exitValue, v)
	}

	// wait for expire
	time.Sleep(time.Duration(timeout) * time.Second)

	// want to get err
	if _, err := cli.Get(exitKey); err == nil {
		t.Fatal("expect get err")
	}
}

func Test_HGet(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data

	key := randomStr()
	field := "name"
	value := "bob"

	if err := cli.HSet(key, field, value); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)

	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"exitData", args{key, field}, "bob", false},
		{"notExitData", args{key, "wrong"}, "", true},
		{"notExitKey", args{"invalidKey", "wrong"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.HGet(tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_HMGet(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data

	key := randomStr()
	field := "name"
	value := "bob"

	if err := cli.HSet(key, field, value); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)

	type args struct {
		key   string
		field []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"exitData", args{key, []string{"name"}}, []string{"bob"}, false},
		{"exitAndNoExist", args{key, []string{"name", "age"}}, []string{"bob", ""}, false},
		{"notExitField", args{key, []string{"wrong"}}, nil, false},
		{"notExitKey", args{"invalidKey", []string{"wrong"}}, nil, false},
	}
	for i, tt := range tests {
		_ = i
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.HMGet(tt.args.key, tt.args.field...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("HGet() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_HMSet(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()
	defer cli.Del(key)
	field := "name"
	value := "bob"

	fieldValue := make(map[string]string)
	fieldValue[field] = value
	if err := cli.HMSet(key, fieldValue); err != nil {
		t.Fatal(err)
	}

	if get, err := cli.HMGet(key, field); err != nil {
		t.Fatal(err)
	} else {
		if get == nil {
			return
		}
		if len(get) != 1 || get[0] != value {
			t.Errorf("except []string{bob} but get %v", get)
		}
	}
}

func Test_Incr(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()
	defer cli.Del(key)
	tests := []struct {
		name    string
		key     string
		want    int
		wantErr bool
	}{
		{"baseOnNil", key, 1, false},
		{"baseOnOne", key, 2, false},
		{"baseOnTwo", key, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.Incr(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Incr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Incr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IncrBy(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()
	defer cli.Del(key)
	type args struct {
		key  string
		incr int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"plus", args{key, 10}, 10, false},
		{"subtract", args{key, -5}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.IncrBy(tt.args.key, tt.args.incr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncrBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MSet(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data

	existKey1 := randomStr()
	existValue1 := randomStr()
	existKey2 := randomStr()
	existValue2 := randomStr()

	clis := make(map[string]string)
	clis[existKey1] = existValue1
	clis[existKey2] = existValue2
	if err = cli.MSet(clis); err != nil {
		t.Fatal(err)
	}
	cli.Del(existKey1, existKey2)
}

func Test_SAdd(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()
	defer cli.Del(key)
	type args struct {
		key     string
		members []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"sadd", args{key, []string{"v1", "v2", "v3"}}, false},
		{"duplicate", args{key, []string{"v0", "v2", "v3", "v4"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.SAdd(tt.args.key, tt.args.members...); (err != nil) != tt.wantErr {
				t.Errorf("SAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_SIsMember(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()
	existValue := randomStr()
	noExistValue := randomStr()

	if err := cli.SAdd(key, existValue); err != nil {
		t.Fatal(err)
	}

	defer cli.Del(key)

	type args struct {
		key    string
		member string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"existMember", args{key, existValue}, true, false},
		{"notExistMember", args{key, noExistValue}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.SIsMember(tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("SIsMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SIsMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SMembers(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	// prepare test data
	key := randomStr()

	member1 := randomStr()
	member2 := randomStr()

	if err := cli.SAdd(key, member1, member2); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)

	tests := []struct {
		name    string
		key     string
		want    []string
		wantErr bool
	}{
		{"validMember", key, []string{member2, member1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.SMembers(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("SMembers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SRem(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	m1 := randomStr()
	m2 := randomStr()
	if err := cli.SAdd(key, m1, m2); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)
	// prepare test data
	type args struct {
		key     string
		members []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"existMember", args{key, []string{m1, m2}}, false},
		{"notExistMember", args{key, []string{randomStr()}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.SRem(tt.args.key, tt.args.members...); (err != nil) != tt.wantErr {
				t.Errorf("SRem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_SetEx(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	if err := cli.SetEx(key, randomStr(), 1); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if _, err := cli.Get(key); err == nil {
		t.Fatal("expect expired")
	}
}

func Test_SetNx(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	existKey := randomStr()
	notExistKey := randomStr()
	if err := cli.Set(existKey, randomStr()); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(existKey, notExistKey)
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"notExist", args{notExistKey, randomStr()}, true, false},
		{"Exist", args{existKey, randomStr()}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.SetNx(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetNx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ZAdd(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v1 := randomStr()
	v2 := randomStr()
	vs := make(map[string]int)
	vs[v1] = 1
	vs[v2] = 2

	defer cli.Del(key)

	type args struct {
		key    string
		values map[string]int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{key, vs}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.ZAdd(tt.args.key, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("ZAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ZRem(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v1 := randomStr()
	v2 := randomStr()
	vs := make(map[string]int)
	vs[v1] = 1
	vs[v2] = 2
	if err := cli.ZAdd(key, vs); err != nil {
		t.Fatal(err)
	}

	defer cli.Del(key)

	type args struct {
		key     string
		members []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{key, []string{v1}}, false},
		{"validWithInvalid", args{key, []string{v1, randomStr()}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.ZRem(tt.args.key, tt.args.members...); (err != nil) != tt.wantErr {
				t.Errorf("ZRem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestcliRedis_ZScore(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v1 := randomStr()
	v2 := randomStr()
	vs := make(map[string]int)
	vs[v1] = 1
	vs[v2] = 2
	if err := cli.ZAdd(key, vs); err != nil {
		t.Fatal(err)
	}

	defer cli.Del(key)
	type args struct {
		key    string
		member string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"v1", args{key, v1}, 1, false},
		{"v2", args{key, v2}, 2, false},
		{"invalid", args{key, randomStr()}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.ZScore(tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZScore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ZRange(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v0 := randomStr()
	v1 := randomStr()
	v2 := randomStr()

	clis := make(map[string]int)
	clis[v0] = 1
	clis[v1] = 2
	clis[v2] = 3
	if err := cli.ZAdd(key, clis); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)

	type args struct {
		key   string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"getAll", args{key, 0, -1}, []string{v0, v1, v2}, false},
		{"get0-1", args{key, 0, 1}, []string{v0, v1}, false},
		{"get2-1000", args{key, 2, 1000}, []string{v2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.ZRange(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ZRangeWithScores(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v0 := randomStr()
	v1 := randomStr()
	v2 := randomStr()

	clis := make(map[string]int)
	clis[v0] = 1
	clis[v1] = 2
	clis[v2] = 3
	if err := cli.ZAdd(key, clis); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)
	type args struct {
		key   string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{"valid", args{key, 0, -1}, clis, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.ZRangeWithScores(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRangeWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRangeWithScores() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ZRangeByScore(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v0 := randomStr()
	v1 := randomStr()
	v2 := randomStr()

	clis := make(map[string]int)
	clis[v0] = 1
	clis[v1] = 2
	clis[v2] = 3
	if err := cli.ZAdd(key, clis); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)

	type args struct {
		key   string
		start interface{}
		end   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"getAll", args{key, "-inf", "+inf"}, []string{v0, v1, v2}, false},
		{"get0-1", args{key, 0, 1}, []string{v0}, false},
		{"get2-1000", args{key, 2, 1000}, []string{v1, v2}, false},
		{"nokey", args{"test", 0, 3}, []string{}, false},
		{"nodata", args{key, 4, 8}, []string{}, false},
	}
	for i, tt := range tests {
		_ = i
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.ZRangeByScore(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ZRangeByScoreWithScores(t *testing.T) {
	cli, err := InitRedis(&rConf)
	if err != nil {
		t.Fatal(err)
	}
	key := randomStr()
	v0 := randomStr()
	v1 := randomStr()
	v2 := randomStr()

	clis := make(map[string]int)
	clis[v0] = 1
	clis[v1] = 2
	clis[v2] = 3
	if err := cli.ZAdd(key, clis); err != nil {
		t.Fatal(err)
	}
	defer cli.Del(key)
	type args struct {
		key   string
		start interface{}
		end   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{"getAll", args{key, "-inf", "+inf"}, clis, false},
		{"valid", args{key, 0, 4}, clis, false},
		{"invalid", args{key, 4, 5}, map[string]int{}, false},
		{"nokey", args{"test", 4, 5}, map[string]int{}, false},
	}
	for i, tt := range tests {
		_ = i
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.ZRangeByScoreWithScores(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRangeWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRangeWithScores() got = %v, want %v", got, tt.want)
			}
		})
	}
}
