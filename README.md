Mini Text Indexer
=================

Description
-----------
Mini Text Indexer is an application which scans and indexes user definable text sequences in text-based files. You configure a set of paths and file patterns, then the Mini Text Indexer will scan these files looking for these patterns. When it finds them it creates an in-memory index which allows for quick searching for these matched patterns.

For example, let's pretend you have pointed Mini Text Indexer to a directory full of JavaScript code and configured it to look at files with a **.js** extension. Then let's pretend you have configured Mini Text Indexer to look for the following text patterns.

```js
\$\("#(.*?)"\) // jQuery ID selectors
new\s+(.*?) // Instantiating objects with new
```

When configured with these patterns Mini Text Indexer will capture the above text and create an in-memory index. Mini Text Indexer provides HTTP endpoints to perform searches against this index to find what files have text that match your search.

Configuration
-------------
Mini Text Indexer is configured via a JSON file. You must tell this application about three things.

0. Where to look for files
0. What file patterns to index
0. What patterns to look for

The basic shell of the configuration file looks like this.

```json
{
	"paths": [],
	"filePatterns": [],
	"textPatterns": []
}
```

### Paths
Paths tell Mini Text Indexer where to look for files. This is a simple array of string directory paths.

```json
{
	"paths": [
		"/code/js/project/models",
		"/code/js/project/services"
	]
}
```

### File Patterns
This section tells Mini Text Indexer what file patterns to consider. File patterns are a simple *contains* match. If there is more than one file pattern defined it will match on any of the provided patterns. In the examples below we are saying we want to match anything that contains **.js** in the file name.

```json
{
	"filePatterns": [
		".js",
		".hbs"
	]
}
```

### Text Patterns
Text Patterns tells Mini Text Indexer what patterns to look for and index. A pattern consists of the following elements.

* Regex pattern with zero or more capture groups
* An index to the capture group which is to be used as the key for the index tree

When a regex pattern is matched it is stored in the index tree. The value that is stored as the key, and used in searches across the tree, should be an index to a capture group in the regular expression. A value of zero (0) tells Mini Text Indexer to use the whole capture as the key.

```json
{
	"textPatterns": [
		{
			"pattern": "\\$\\(\"#(.*?)\"\\)",
			"key": 1
		}
	]
}
```

### Startup Configuration
Mini Text Indexer is a command line server application. It has several command line flags that can control and customize its behavior.

* **ip** - Address to bind the HTTP server to
* **port** - Port to bind the HTTP server to
* **loglevel** - Detail level of logging: *debug*, *info*

HTTP Interface
--------------
Mini Text Indexer provides an HTTP interface to perform searches against the index tree. Below are the endpoints available.

### Search

#### GET /search?term=[searchTerm]
Performs a search against the index tree. This will return an array of terms that matches the specified search term.

The matching tree node contains a key which is the match to the provided search term. It then has an array of documents where the term is found. Each document has a name, followed by an array of match locations. Each location has the matched text, captured groups from the regular expression, and the starting location of the text in the file.

##### Parameters
* **term** - Term to search for

##### Response
```json
[
	{
		"key": "contentDiv",
		"documents": [
			{
				"documentName": "HomeController.js",
				"matches": [
					{
						"location": 100,
						"match": "$(\"#contentDiv\")",
						"captures": [
							"contentDiv"
						]
					}
				]
			}
		]
	},
	{
		"key": "contentDivabc",
		"documents": [
			{
				"documentName": "TestController.js",
				"matches": [
					{
						"location": 10,
						"match": "$(\"#contentDivabc\")",
						"captures": [
							"contentDivabc"
						]
					}
				]
			}
		]
	}
]
```

#### GET /getterm?term=[searchTerm]
Performs a search against the index tree. This will return a specific term that matches the specified search term.

The matching tree node contains a key which is the match to the provided search term. It then has an array of documents where the term is found. Each document has a name, followed by an array of match locations. Each location has the matched text, captured groups from the regular expression, and the starting location of the text in the file.

##### Parameters
* **term** - Term to search for

##### Response
```json
{
	"key": "contentDiv",
	"documents": [
		{
			"documentName": "HomeController.js",
			"matches": [
				{
					"location": 100,
					"match": "$(\"#contentDiv\")",
					"captures": [
						"contentDiv"
					]
				}
			]
		}
	]
}
```

License
-------

The MIT License (MIT)

Copyright (c) 2015 Adam Presley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
