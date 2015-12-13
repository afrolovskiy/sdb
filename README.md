# Simple Database

A simple in-memory database.

Database accept the following commands:
* `SET name value` – Set the variable name to the value. Neither variable names nor values contain spaces.
* `GET name` – Print out the value of the variable name, or NULL if that variable is not set.
* `UNSET name` – Unset the variable name, making it just like that variable was never set.
* `NUMEQUALTO value` – Print out the number of variables that are currently set to value. If no variables equal that value, print 0.
* `BEGIN` – Open a new transaction block. Transaction blocks can be nested; a BEGIN can be issued inside of an existing block.
* `ROLLBACK` – Undo all of the commands issued in the most recent transaction block, and close the block. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.
* `COMMIT` – Close all open transaction blocks, permanently applying the changes made in them. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.
* `END` – Exit the program. Your program will always receive this as its last command.

## Installation

Was tested on OSX. Should also work on Unix systems.

Requires Go 1.3+. Make sure you have Go properly installed, including setting up your GOPATH. 

Create directory in your GOPATH and move folder with project there:

    $ mkdir -p $GOPATH/src/github.com/afrolovskiy
    $ mv sdb $GOPATH/src/github.com/afrolovskiy

Now you can install it:

    $ go install github.com/afrolovskiy/sdb/cmd/sdb

## Usage

Now you can use `sdb` by passing commands from command line or via stdin.

    $ echo -e "SET a 10\nGET a" | sdb
    10

    $ sdb
    GET a
    NULL
    SET a 10
    GET a
    10
    END
