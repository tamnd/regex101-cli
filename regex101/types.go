package regex101

// Pattern is a shared regex pattern from regex101.com.
type Pattern struct {
	Rank        int    `json:"rank"         csv:"rank"         tsv:"rank"`
	Title       string `json:"title"        csv:"title"        tsv:"title"`
	Flavor      string `json:"flavor"       csv:"flavor"       tsv:"flavor"`
	Author      string `json:"author"       csv:"author"       tsv:"author"`
	Upvotes     int    `json:"upvotes"      csv:"upvotes"      tsv:"upvotes"`
	Downvotes   int    `json:"downvotes"    csv:"downvotes"    tsv:"downvotes"`
	DateCreated string `json:"date_created" csv:"date_created" tsv:"date_created"`
	URL         string `json:"url"          csv:"url"          tsv:"url"`
}
