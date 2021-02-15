package main

import (
	"flag"
	"strings"
)

func main() {
	// add flags
	add := flag.String("add", "", "add an entry to the diary. Date format: today:, yesterday:, YYYY-MM-DD")
	remove := flag.String("remove", "", "remove an entry from the diary. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	view := flag.String("view", "", "view an entry or all entries from the diary. Use all to see all. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	searchkeywords := flag.String("searchkeywords", "", "search entries by keyword")
	searchtags := flag.String("searchtags", "", "search entries by tags")
	searchfields := flag.String("searchfields", "", "search entries by fields")
	plaintext := flag.Bool("plaintext", false, "show as plaintext")
	flag.Parse()

	// no commands were provided and no text was written
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.PrintDefaults()
		// now exit
		return
	}

	// vreate empty Journal
	j := crate_journal()
	// load from database
	e := j.load()

	if e != nil {
		print_error(e, 3)
		return
	}

	// no commands were provided but some text was recognized
	if flag.NFlag() == 0 && flag.NArg() > 0 {
		entry := strings.Join(flag.Args(), " ")
		j.createEntry(entry)
	} else if *add != "" {
		// get text provided by the flag
		// get remainder text
		// concantenate them
		entry := string(*add) + " " + strings.Join(flag.Args(), " ")
		j.createEntry(entry)
	} else if *remove != "" {
		e := j.removeEntry(*remove)
		if e != nil {
			print_error(e, 2)
		}
	} else if *view != "" {
		// get entry by date
		// check if parameter is "all"
		if strings.ToLower(*view) == "all" {
			entries, e := j.getAllEntries()
			if e != nil {
				print_error(e, 1)
			} else {
				for _, entry := range entries {
					print_entry(entry, *plaintext)
				}
			}
		} else {
			// check if the parameter is some kind of date
			entry, e := j.viewEntry(*view)
			if e != nil {
				print_error(e, 1)
			} else {
				print_entry(entry, *plaintext)
			}
		}
	} else if *searchkeywords != "" {
		var keywords []string
		// concantenate all the keywords
		keywords = append(keywords, *searchkeywords)
		keywords = append(keywords, flag.Args()...)
		entries, e := j.searchKeywords(keywords)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *plaintext)
			}
		}
	} else if *searchtags != "" {
		var tags []string
		// concantenate all the keywords
		tags = append(tags, *searchtags)
		tags = append(tags, flag.Args()...)
		entries, e := j.searchTags(tags)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *plaintext)
			}
		}
	} else if *searchfields != "" {
		var keys []string
		// concantenate all the fields keys
		keys = append(keys, *searchfields)
		keys = append(keys, flag.Args()...)
		entries, e := j.searchFields(keys)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *plaintext)
			}
		}
	}

	e = j.save()

	if e != nil {
		print_error(e, 1)
	}
}
