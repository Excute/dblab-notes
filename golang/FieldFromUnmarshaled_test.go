package test

func TestFieldFromUnmarshaled(t *testing.T) {
	Resp := &CmResponses{}
	Reso := &CmResources{}
	initData(Resp, Reso)

	var v interface{}

	// if err := json.Unmarshal([]byte(Reso.Map["Pod"].Manifest), &v); err != nil {
	if err := json.Unmarshal([]byte(`
    {
      "containers": [
        {
          "volumes":[
            {"path":"A"},
            {"path":"B"},
            {"path":"C"}
          ],
          "image":"nginx"
        },{
          "volumes":[
            {"path":"D"},
            {"path":"E"},
            {"path":"F"}
          ],
          "image":"mySQL"
        }
      ]
    }
    `), &v); err != nil {
		t.Error(err)
	}

	// target := "kind"
	target := "containers.volumes.path"

	is, err := FieldFromUnmarshaled(v, target)
	if err != nil {
		panic(err)
	}
	for _, i := range is {
		t.Logf("target: %s, type: %s, value: %s", target, reflect.ValueOf(i).Kind(), i)
	}
}

func FieldFromUnmarshaled(input interface{}, target string) ([]interface{}, error) {
	// fmt.Printf("FFU: Called\n\ttarget:%s\n\tinput:\n%s\n", target, input)
	if target == "" || len(target) < 1 {
		return []interface{}{input}, nil
	}

	// 입력을 reflect
	rM := reflect.ValueOf(input)

	// fmt.Printf("FFU: Type of input: %s\n", rM.Kind())

	// 현재 찾아야 할 키와 나중에 찾을 키를 구분
	var (
		currentTarget string // 현재 찾아야 할 키
		nextTarget    string // 다음 재귀호출의 target 키 인수
	)
	if i := strings.Index(target, "."); i > 0 {
		currentTarget = target[:i]
		nextTarget = target[i+1:]
	} else {
		currentTarget = target
		nextTarget = ""
	}

	switch rM.Kind() {
	case reflect.Map:
		// 입력이 Map 일 경우

		i := rM.MapRange()
		tmps := []interface{}{}
		for i.Next() {
			if i.Key().String() == currentTarget {
				// 매니페스트의 인수가 interface{}이여야 해서 reflect된 매니페스트(rM)을 그대로 사용하지 못하고
				// 매니페스트 m을 맵 타입으로 parse 하여 키로 값을 찾아 인수로 전달
				tmp, err := FieldFromUnmarshaled(input.(map[string]interface{})[currentTarget], nextTarget)
				if err != nil {
					return nil, err
				}
				tmps = append(tmps, tmp...)
			}
		}
		return tmps, nil
	case reflect.Slice:
		// 입력이 slice인 경우
		res := []interface{}{}
		for i := 0; i < rM.Len(); i++ {
			is, err := FieldFromUnmarshaled(input.([]interface{})[i], target)
			if err != nil {
				return nil, err
			}
			res = append(res, is...)
		}
		return res, nil
	default:
		// 더이상 하위 구조가 없는 경우
		// 더 찾아야 하는 키가 있는 경우 에러
		if len(target) > 0 {
			return nil, errors.New("structure depth is not that deep")
		} else {
			return []interface{}{input}, nil
		}
	}