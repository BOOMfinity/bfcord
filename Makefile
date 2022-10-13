field-alignments:
	find . -not -path '*/.*' -type d -exec go run test/analyzers/main.go -fieldalignment -fix {} \;