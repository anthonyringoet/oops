## oops

> A toy implementation of a (fake) ransomware-style program

If you're looking for some kind of evil piece of software, i'll have to dissappoint you.
It's really not actual ransomware. Not by far.

The file selection is currently hardcoded to select all files in a directory named `target` in the current working directory.

## Run

```bash
go run main.go -mode=encrypt -dryrun
go run main.go -mode=decrypt -dryrun

# help
go run main.go -h
-dryrun
      Set to true for dry run.
-mode string
      Set mode to either 'encrypt' or 'decrypt'. (default "encrypt")
```

## Build

```bash
# TBC
```

## Test

```bash
# run tests for all packages
go test ./...

# run single test case
go test -run TestGenerateKeyPair

# run multiple cases using regex
go test -run TestGenerate.*
```

## Manual test

```bash
./seed.sh
go run main.go

# output below
Starting oops...
Generating keypair
Keys already exist, continuing...
Getting file list
Total files: 70
⏳⏳⏳ START Processing mode=encrypt dryRun=false

⏳ Processing file=/Users/someone/oops/target/dir1/file7.txt
⏳ Processing file=/Users/someone/oops/target/dir2/file7.txt
⏳ Processing file=/Users/someone/oops/target/dir2/file8.txt
⏳ Processing file=/Users/someone/oops/target/dir1/file3.txt
🔐 Encrypting file=/Users/someone/oops/target/dir2/file8.txt
🔐 Encrypting file=/Users/someone/oops/target/dir1/file3.txt
✅ Processed /Users/someone/oops/target/dir2/file8.txt in 1.174709ms
⏳ Processing file=/Users/someone/oops/target/dir4/file9.txt
🔐 Encrypting file=/Users/someone/oops/target/dir4/file9.txt
✅ Processed /Users/someone/oops/target/dir4/file8.txt in 985.5µs
...
Processed all files in 9.099709ms
Processed 70 of 70 files. Failed: 0


# afterwards to decrypt
go run main.go -mode=decrypt
```
