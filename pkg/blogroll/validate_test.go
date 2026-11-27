package blogroll

import "testing"

func TestValidateFeedDocument(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want string
	}{
		{"RSS", "<rss version=\"2.0\"><channel><title>Example</title></channel></rss>", "rss"},
		{"Atom", "<feed xmlns=\"http://www.w3.org/2005/Atom\"><title>Example</title></feed>", "atom"},
		{"JSON Feed", `{"version":"https://jsonfeed.org/version/1.1","items":[]}`, "json"},
		{"HTML", "<html><head><title>Not a feed</title></head></html>", ""},
		{"XML declaration with HTML", "<?xml version=\"1.0\"?><html/>", ""},
		{"Missing RSS channel", "<rss/>", ""},
		{"Malformed XML", "<rss><channel></rss>", ""},
		{"Unrelated JSON", `{"error":"not found"}`, ""},
		{"Empty", "", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateFeedDocument([]byte(tt.body))
			if got != tt.want || (err != nil) != (tt.want == "") {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
