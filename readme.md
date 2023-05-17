## oops

> A toy implementation of a (fake) ransomware-style program

## features

Here is a high-level list of steps that the toy ransomware-style encrypting program could perform on startup:

- **Setup**: Initialize necessary variables, configuration settings, and load any required packages or modules.
- **Key Pair Generation**: Generate a public/private key pair using a cryptographic library. This will be used for the encryption and decryption of files.
- **Report Status**: Notify that the program has started and the key pair has been generated successfully. This could be a simple print statement to stdout.
- **Directory Selection**: Specify the directory that will be the target for encryption. Make sure this is a safe directory that doesn't contain any important files.
- **Recursive File Search**: Traverse the specified directory and its subdirectories to find all files. Store these filenames in a list or similar data structure for later use.
- **Encrypt Files**: Iterate through the list of filenames, encrypting each one using the public key. Each successful encryption should remove the original file.
- **Encryption Report**: After all files have been encrypted, report the status of the operation. This could include the number of files encrypted, any errors encountered, etc.
- **Start** Decryption Listener: Start a function that constantly checks for a signal to decrypt the files. This could be the presence of a specific file in the directory.

## Run

```bash
# TBC
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

## CLI

```bash
  -dryrun
        Set to true for dry run.
  -mode string
        Set mode to either 'encrypt' or 'decrypt'. (default "encrypt")

# examples
go run main.go -mode=encrypt -dryrun
go run main.go -mode=decrypt -dryrun
```