package ctx

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/hpifu/go-kit/hconf"
)

type Ctx struct {
	mutex sync.RWMutex
	store *hconf.InterfaceStorage
}

func NewCtx() *Ctx {
	return &Ctx{
		store: hconf.NewInterfaceStorage(map[string]interface{}{}),
	}
}

func (c *Ctx) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.store.Get(key)
}

func (c *Ctx) Set(key string, val interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.store.Set(key, val)
}

var expRegex = regexp.MustCompile(`(.+?)\((.*?)\)`)

func (c *Ctx) EvaluateExpr(expr string) (interface{}, error) {
	vals := expRegex.FindStringSubmatch(expr)
	if len(vals) == 0 {
		return c.Get(expr)
	}
	//
	//fun := vals[1]
	//params := vals[2]

	//strings.Split(params, ",")

	return nil, nil
}

func (c *Ctx) Evaluate(data *interface{}) (interface{}, error) {
	switch (*data).(type) {
	case string:
		v := (*data).(string)
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			text := v[2 : len(v)-2]
			if strings.HasPrefix(text, "type.int64(") {
				val, err := strconv.ParseInt(text[11:len(text)-1], 10, 64)
				if err != nil {
					return nil, err
				}
				*data = val
			} else {
				key := strings.TrimSpace(text)
				val, err := c.Get(key)
				if err != nil {
					return nil, err
				}
				*data = val
			}
		}
	case map[string]interface{}:
		for k, v := range (*data).(map[string]interface{}) {
			if val, err := c.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).(map[string]interface{})[k] = val
			}
		}
	case map[interface{}]interface{}:
		for k, v := range (*data).(map[interface{}]interface{}) {
			if val, err := c.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).(map[interface{}]interface{})[k] = val
			}
		}
	case []string:
		for i, v := range (*data).([]string) {
			var vi interface{}
			vi = v
			if val, err := c.Evaluate(&vi); err != nil {
				return nil, err
			} else {
				(*data).([]interface{})[i] = val
			}
		}
	case []interface{}:
		for i, v := range (*data).([]interface{}) {
			if val, err := c.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).([]interface{})[i] = val
			}
		}
	default:
	}

	return *data, nil
}
