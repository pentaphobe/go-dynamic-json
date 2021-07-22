# Dynamic JSON types in Golang

- [Dynamic JSON types in Golang](#dynamic-json-types-in-golang)
	- [Example implementation](#example-implementation)
		- [Marshalling / Unmarshalling](#marshalling--unmarshalling)
		- [Next steps](#next-steps)
	- [Background](#background)
		- [Previous attempts](#previous-attempts)
			- [Hide behind an interface](#hide-behind-an-interface)
			- [Lots of empty pointers](#lots-of-empty-pointers)
- [References](#references)

Ok, so you've got some JSON or YAML like the following:

```yaml
- type: person
  name: fred
- type: company
  abn: 123-456
- type: person
  name: james
```

And you want to load it in Golang but with _some_ semblance of type safety?

Well, good luck - because doing so unequivocally sucks.

> :warning: _And it still sucks with this library, though hopefully slightly less_

## Example implementation

Here are our basic Go structures as defined by the above YAML example

```go
type Person struct {
	Name string `json:"name"`
}

type Company struct {
	ABN string `json:"abn"`
}
```

In order to make them dynamically loadable we need two basic bits of boilerplate:

1. Implement the `reflected.Payload` interface for all subtypes
	 ```go
	 func (t *Person) TypeName() string {
			return "person"
	 }

	 func (t *Person) Cast(i interface{}) Payload {
			return i.(*Person)
	 }



	 func (t *Company) TypeName() string {
			return "company"
	 }

	 func (t *Company) Cast(i interface{}) Payload {
			return i.(*Company)
	 }
	 ```
2. Register your types as "handlers" (I know, I know - working on it...)
	```go
	 func init() {
			reflected.RegisterHandler(&Person{})
			reflected.RegisterHandler(&Company{})
	 }
	```

### Marshalling / Unmarshalling

In order to load/save these types dynamically you just use the standard `Marshal()` and `Unmarshal()` functions (`json.`, or `yaml.`(TBD)) as usual, but using the `TypedContainer` struct

```go
jsonData := `
{
	"type": "person",
	"name": "Freddy Mercury"
}   	
`
var t reflected.TypedContainer
_ := json.Unmarshal([]byte(jsonData), &t)

fmt.Println(t.Type)
fmt.Println(t.Payload)

// And here's the caveat I didn't tell you about earlier:
person := t.Payload.(*Person)
fmt.Println(person.Name)  

// This can be made a little nicer with some helpers, but ultimately
// seems to be unavoidable - reflection hides a lot of the pain but
// a type assertion appears necessary _at some point_
```

> :warning: If you don't read comments in code blocks: _everything is fine, there are no caveats_


### Next steps

Now that this works, it can likely be further simplified to avoid the need for `Payload`, instead dynamically updating fields on the user structures and still magically inserting the `Type` - however it seems there's no ideal way to avoid the final type assertion happening _somewhere_ outside the library.

One option is to allow the end user to define their own structs equivalent to
the `TypedContainer` and just magic them into existence via reflection.  This way they can provide Getters for each of their potential embedded types.

Have yet to try this out, it might be a nice improvement - or may be too ugly to live.

## Background

There are quite a few approaches to this pattern in Go, all of which are nonoptimal and either involve sacrificing readability, lots of boilerplate, or whacky things like just keeping everything in a `map[string]interface{}`

So far I've yet to find anything which is properly satisfying, but this is my latest attempt.

### Previous attempts

#### Hide behind an interface

```go
type Envelope struct {
	Type string `json:"type"`
}

type EntitySubtype interface{}

type Entity struct {
	Envelope
	EntitySubtype
}

// TODO: copy some code in here, because I don't want to have to write it again
```

> "You were supposed to be the chosen one!"

Alas, this one looked somewhat promising - or at least was preferable to the alternatives I'd seen elsewhere (see [references](#references))

#### Lots of empty pointers

```go
type Monolith struct {
	Type string `json:"type"`
	*Person
	*Company
}

type Company struct {
	ABN string
}

type Person struct {
	Name string
}
```

This is obviously quite ugly to look at, and involves lots of `nil` checks all over your business logic code.  But aside from that it's actually the nicest, and easiest to use approach...

*However* the problem with this one is mostly down to clashing fields on embedded structs.  And the necessity to update the `Type` field manually (or make boilerplate to handle setting it dynamically)



# References

A few of the articles, gists, and panicked blog posts I read in the hope that someone had found the Holy Grail