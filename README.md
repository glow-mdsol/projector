# Reporting tool for BodyCheck Application

Pass in a URL pattern and get out a spreadsheet dump of results.  

Intended for use with BodyCheck derived data

```shell
➜  projector git:(feature/multiple_patterns) ✗ ./projector -pattern googleplex
2017/07/24 01:34:56 Processing URL Pattern  googleplex
2017/07/24 01:34:56 Retrieving Subject Counts
2017/07/24 01:34:56 Retrieving Unfired Edits
2017/07/24 01:34:58 Retrieving URL Metrics
2017/07/24 01:35:16 Generated metrics for  2 URLs
2017/07/24 01:35:16 Writing Subject Counts
2017/07/24 01:35:16 Writing Unfired Edits
2017/07/24 01:35:16 Writing Study Metrics
```