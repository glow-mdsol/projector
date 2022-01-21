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

## Pattern

1. Get the list of URLs that match the pattern
2. For each matching URL
    1. Get the Subject Count -> write to sheet
    2. Get the Projects and for each Project (sorted by name)
        1. Get the ProjectVersions and for each ProjectVersion (sorted by CRFVersion)
            1. Get the Unused Edits
                1. With OpenQuery Action -> write to sheet
                2. Without OpenQuery Action -> write to sheet
            3. Get the Edit check counts
                1. write to sheet, ordered by CRF Version
        2. For the last ProjectVersion
            1. Write the summarised counts to the last project sheet
        3. Generate aggregate and average counts for last versions -> write to sheet
