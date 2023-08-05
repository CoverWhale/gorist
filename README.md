# gorist

Gorist is a Grist client for Golang. It is in alpha phase so things will change.

## Getting Records

All records returned from Grist follow this structure:

```go
type Records struct {
	Records []Record `json:"records"`
}

type Record struct {
	ID     string       `json:"id"`
	Fields RecordFields `json:"fields"`
}

type RecordFields struct {
	Open  bool   `json:"open"`
	Name string  `json:"name"`
	IDs  []int   `json:"ids"`
}

```

RecordFields is a struct containing the columns and their types from your dataset. 

## Returning Lists In A Record 

Grist exhibits unique behaviors when it comes to returning records that contain arrays which is especially noticable in a statically typed language. Notably, when a field contains a list, Grist adds a letter character to describe the data type even if the array is only integers.

One way to handle this is to satisfy the UnmarshalJSON interface for that record type and manually ignore the first item. 
Then you need to loop through multiple temporary maps to exract and type assert the values.

For example:

```go

type Records struct {
	Records []Record `json:"records"`
}

type Record struct {
	ID     string       `json:"id"`
	Fields RecordFields `json:"fields"`
}

type RecordFields struct {
	Open  bool   `json:"open"`
	Name string  `json:"name"`
	IDs  []int   `json:"ids"`
}

func (r *Record) UnmarshalJSON(b []byte) error {
	type tempMap map[string]json.RawMessage

	var data map[string]json.RawMessage

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for k := range data {
		if k != "fields" {
			continue
		}

		var t tempMap
		if err := json.Unmarshal(data[k], &t); err != nil {
			return err
		}

		var ids []interface{}
		if err := json.Unmarshal(t["ids"], &ids); err != nil {
			return err
		}

		for _, v := range ids {
			_, ok := v.(string)
			if ok {
				continue
			}

			i, ok := v.(float64)
			if !ok {
				continue
			}

			r.Fields.IDs = append(r.Fields.IDs, int(i))
		}
	}

	return nil
}

```

