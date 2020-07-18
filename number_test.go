package cgminer

import (
	"strings"
	"testing"
)

func TestNumber_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		want Number
		err  string
	}{
		`"32"`: {
			want: Number(32),
		},
		`32`: {
			want: Number(32),
		},
		"":     {},
		`""`:   {},
		"null": {},
		`{"fff": 32}`: {
			err: "Number.UnmarshalJSON: value is not a number",
		},
		`[1,2,3]`: {
			err: "Number.UnmarshalJSON: value is not a number",
		},
	}

	for input, c := range cases {
		t.Run(input, func(t *testing.T) {
			var got Number
			if err := got.UnmarshalJSON([]byte(input)); err != nil {
				if c.err == "" {
					t.Fatal("unexpected error:", err)
				}

				if msg := err.Error(); !strings.Contains(msg, c.err) {
					t.Fatalf("error message doesn't contains %q\nin: %s", c.err, msg)
				}
			}

			if got != c.want {
				t.Fatalf("value mismatch: %f != %f", got, c.want)
			}
		})
	}
}
