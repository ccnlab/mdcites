# mdcites

Markdown citation extractor, reads pandoc-citeproc citations, exports .bib file with those.

The markdown citation format is `[@RefYY]` and the RefYY is the bibtex cite key that is looked up in the source .bib file

```
Usage of mdcites:
  -bib string
    	required full path to .bib file containing all references that could be cited
  -dir string
    	optional directory containing .md files to process (default "./")
  -out string
    	filename for output .bib file containing only the cited references (default "references.bib")
   
 ```


