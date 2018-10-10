# IMDB Titles

Simple fetcher for get target tittles with basic infos from [IMDb]("http://www.imdb.com/")

## Requirements

- docker or go

Use make for docker and go for go

## Build
For get imdb binary in folder, you can use:
```sh
make build
```

## Run
By default will get the 500 most rating titles from each genre:

```sh
make run
```

For run with options you can check with:


### Run with options
Options for pass as var for the fetcher change default action.


for pass args to makefile you can set flags in var args:
```
make run args='--limit=1000 --sort="num_votes,asc" --debug --adult=false'
```


you too can check options with
```sh
imdb -h
```

### Change limit
for change max limit, you can use flag limit
```sh
imdb --limit=1000
```

### Remove adult results
for prevent adult titles you can use flag adult
```sh
imdb --adult=false
```

### Change sort
if you need change sort, you can use flag sort
```sh
imdb --sort="num_votes,asc"
```

### Change items per request
if you need split request in parts, you can use flag count
```sh
imdb --count=100
```

### Debug mode
for activate debug mode, you can use flag debug
```sh
imdb --debug
```

## Output

Output will be in [jsonlines](http://jsonlines.org)
and writed inside a results folder in this path

You could use:
```sh
ls results/*.jsonl
````
for see results

### Example 

```json
{
  "Name": "Apenas Um Show",
  "Episode": "A Regular Epic Final Battle",
  "Year": "(2009–2017)",
  "Genres": [
    "animation",
    "action",
    "adventure"
  ],
  "Rating": {
    "Value": 9.9,
    "Best": 10,
    "Count": 571,
    "Position": 1
  }
}
```

