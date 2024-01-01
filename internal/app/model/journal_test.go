package model

import (
	"testing"
)

func TestJournal_GetDate(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"2018-05-10", "Thursday May 10, 2018"},
		{"200-00-00", ""},
		{"", ""},
		{"0000-00-00", ""},
	}

	for _, table := range tables {
		j := Journal{Date: table.input}
		actual := j.GetDate()
		if actual != table.output {
			t.Errorf("Expected GetDate() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}

func TestJournal_GetEditableDate(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"2018-05-10", "2018-05-10"},
		{"2018-05-10EXTRATHINGS", "2018-05-10"},
		{"200-00-00", ""},
		{"", ""},
		{"0000-00-00", "0000-00-00"},
	}

	for _, table := range tables {
		j := Journal{Date: table.input}
		actual := j.GetEditableDate()
		if actual != table.output {
			t.Errorf("Expected GetEditableDate() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}

func TestJournal_GetExcerpt(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"<p>Some simple text</p>", "Some simple text"},
		{"<p>Multiple</p><p>paragraphs, some with</p><p>multiple words</p>", "Multiple paragraphs, some with multiple words"},
		{"", ""},
		{"<p></p><p></p>", " "},
		{"<p>a b c d e f g h i j k l m n o p q r s t u v w x y z a b c d e f g h i j k l m n o p q r s t u v w x y z</p>", "a b c d e f g h i j k l m n o p q r s t u v w x y z a b c d e f g h i j k l m n o p q r s t u v w x..."},
	}

	for _, table := range tables {
		j := Journal{Content: table.input}
		actual := j.GetExcerpt()
		if actual != table.output {
			t.Errorf("Expected GetExcerpt() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}

func TestSlugify(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"A SIMPLE TITLE", "a-simple-title"},
		{"already-slugified", "already-slugified"},
		{"   ", "---"},
		{"lower cased", "lower-cased"},
		{"Special!!!Characters@$%^&*(", "special---characters-------"},
	}

	for _, table := range tables {
		actual := Slugify(table.input)
		if actual != table.output {
			t.Errorf("Expected Slugify() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}
