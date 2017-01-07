# tar-stream-merge

A native Go CLI tool for merging two tar file streams (STDIN and named pipe) and streaming the result to STDOUT

## Notes

1. tar-stream-merge merges records from the input streams as each becomes available
2. tar-stream-merge will buffer at most one file from each input stream into memory at a time. This may change in the future to a bufferless implementation.
3. tar-stream-merge will merge archive entries with names that already exist in the archive - this is purposeful.
5. tar-stream-merge only writes the tarball to STDOUT. Error messages will be sent on STDERR.

## Examples

Suppose you have a tar file that you'd like to merge into a tar stream. 

    # setup some files
    echo first file > 1.txt
    echo second file > 2.txt
    tar cf input.tar 1.txt 2.txt
    echo third file > 3.txt
    tar cf piped.tar 3.txt

    # merge the tars
    tar-stream-merge <(cat piped.tar) < input.tar > merged.tar

    # build a pipeline
    echo fourth file > 4.txt
    tar cf piped2.tar 4.txt
    cat input.tar | tar-stream-merge <(cat piped.tar) | tar-stream-merge <(cat piped2.tar) > pipelined.tar

### Authors and Copyright Holders

1. Jeff Nickoloff "jeff@allingeek.com"

### License

This project is released under the MIT license.
