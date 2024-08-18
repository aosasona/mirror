package generator

import "testing"

type Addr struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
}

type Student struct {
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Age         int            `json:"age"`
	Email       string         `json:"email"`
	Contact     []string       `json:"contact"`
	Grades      map[string]int `json:"grades"`
	unexported  string
	Ignored     string `json:"-"`
	IgnoredTS   string `mirror:"-"`
	IgnoredBoth string `json:"-" mirror:"-"`
	Addr        `json:"address"`
}

func Test_ObjectTypeGeneration(t *testing.T) {
	tests := []struct {
		Name     string
		Source   any
		Expected string
	}{
		{
			Name:   "student",
			Source: Student{},
			Expected: `export type Student = {
    first_name: string;
    last_name: string;
    age: number;
    email: string;
    contact: string[];
    grades: Record<string, number>;
    address: Addr;
}`,
		},
		{
			Name:   "address",
			Source: Addr{},
			Expected: `export type Addr = {
    street: string;
    city: string;
    postcode: string;
}`,
		},
	}

	g := NewGenerator(Opts{
		UseTypeForObjects: true,
	})

	for _, tt := range tests {
		got := g.Generate(tt.Source)
		if got != tt.Expected {
			t.Errorf("`%s`: got\n%v, want\n%v", tt.Name, got, tt.Expected)
		}
	}
}

type Entry struct {
	Title     string `json:"title" mirror:"name:entry_name"`
	Link      string `json:"link"`
	isPrivate bool   `json:"is_private"`
}

type BoxedEntry struct {
	Name  string
	entry Entry
}

func Test_PrivateFieldExport(t *testing.T) {
	tests := []struct {
		Name     string
		Source   any
		Expected string
	}{
		{
			Name:   "entry",
			Source: Entry{},
			Expected: `export type Entry = {
    entry_name: string;
    link: string;
    is_private: boolean;
}`,
		},
		{
			Name:   "boxed_entry",
			Source: BoxedEntry{},
			Expected: `export type BoxedEntry = {
    Name: string;
    entry: Entry;
}`,
		},
	}

	g := NewGenerator(Opts{
		AllowUnexportedFields: true,
		UseTypeForObjects:     true,
	})

	for _, tt := range tests {
		got := g.Generate(tt.Source)
		if got != tt.Expected {
			t.Errorf("`%s`: got\n%v, want\n%v", tt.Name, got, tt.Expected)
		}
	}
}
