#!/bin/bash

# Define the dummy text
dummy_text="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas bibendum."

# Create the target directory, deleting it first if it already exists
rm -rf target
mkdir target

# Change to the target directory
cd target

# Create 10 text files
for i in {1..10}; do
    echo $dummy_text > "file$i.txt"
done

# Create 5 directories
for i in {1..5}; do
    mkdir "dir$i"

    # Change to the new directory
    cd "dir$i"

    # Create 10 text files in this directory
    for j in {1..10}; do
        echo $dummy_text > "file$j.txt"
    done

    # Create another directory inside this one
    mkdir "subdir"

    # Change to the new subdirectory
    cd "subdir"

    # Create 2 text files in this subdirectory
    for k in {1..2}; do
        echo $dummy_text > "file$k.txt"
    done

    # Go back to the parent directory (i.e., target)
    cd ../..

done
