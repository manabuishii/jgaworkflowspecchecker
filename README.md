# specchecker

Spec checker

# Setup

```
GO111MODULE=on
```

# Develop

## Test

### Simply test

```
go test ./...
```

### Coverage

```
go test ./... -race -coverprofile=coverage.txt -covermode=atomic 
go tool cover -html=coverage.txt -o cover.html
```

```
open cover.html
```

# schematyper

```
schematyper simple-schema.json
```

# Execute

```
go run main.go simple-schema-file-md5.json simple-data-file-md5.json
```

## Execute case 1

```
go run main.go schema/samplesheet_schema.json schema/samplesheet_data.json
```

### 20210702 reseult

```console
$ go run main.go schema/samplesheet_schema.json schema/samplesheet_data.json
The document is valid
Load sample
Load end
Name is [aaabbb]
Fq1 is []
Md5 is []
index: 0, SampleId: XX10000
index: 0, SampleId: hello runid1
result=false
index: 1, SampleId: hello runid2
result=false
index: 2, SampleId: hello runid3
result=false
index: 1, SampleId: XX10001
index: 0, SampleId: hello runid4
result=false
index: 1, SampleId: hello runid5
result=false
index: 2, SampleId: hello runid6
result=false
```

## Execute case 2

```
go run main.go testschema/sampleid_list_runid_schema.json testschema/sampleid_list_runid_data.json
```

## Execute case 3

```
go run main.go schema/samplesheet_schema.json schema/samplesheet_data.json
```

# TODO

- [ ] two PE_2 entries exists in same runid, check whether error is happens
- [ ] decide `platform` accepts free words.
- [ ] Option, estimate hash value check time estimate by file size and hash algorithm
- [ ] Option, no hash value check
- [ ] `haplotypecaller_chrX_nonPAR_ploidy_2_interval_bed` and `haplotypecaller_chrX_nonPAR_ploidy_1_interval_bed` are same file or separate files
- [ ] output version string
- [ ] outpu git commit id
- [ ] get exit status of each command
