package main

import (
	"encoding/json"
	"fmt"

	"github.com/pentaphobe/go-dynamic-json/pkg/reflected"
)

type Person struct {
	Name string `json:"name"`
}

type Company struct {
	ABN string `json:"abn"`
}

func (t *Person) TypeName() string {
	return "person"
}

func (t *Person) Cast(i interface{}) reflected.Payload {
	return i.(*Person)
}

func (t *Company) TypeName() string {
	return "company"
}

func (t *Company) Cast(i interface{}) reflected.Payload {
	return i.(*Company)
}

func init() {
	reflected.RegisterHandler(&Person{})
	reflected.RegisterHandler(&Company{})
}

func main() {
	jsonData := `
{
	"type": "person",
	"name": "Freddy Mercury"
}   	
`
	var t reflected.TypedContainer
	_ = json.Unmarshal([]byte(jsonData), &t)

	fmt.Println(t.Type)

	dump, _ := json.MarshalIndent(&t, "", "  ")
	fmt.Println(string(dump))

	// And here's the caveat I didn't tell you about earlier:
	person := t.Payload.(*Person)
	fmt.Println(person.Name)

	// This can be made a little nicer with some helpers, but ultimately
	// seems to be unavoidable - reflection hides a lot of the pain but
	// a type assertion appears necessary _at some point_

}
