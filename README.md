# rambler [![build](https://app.wercker.com/status/b645428b6f548288d71d3ba83cc1a783/s/master "wercker status")](https://app.wercker.com/project/bykey/b645428b6f548288d71d3ba83cc1a783) [![Coverage Status](https://coveralls.io/repos/elwinar/rambler/badge.svg?branch=master&service=github)](https://coveralls.io/github/elwinar/rambler?branch=master)

A simple and language-independent SQL schema migration tool

## Installation

You can download the latest release on the [release page](https://github.com/elwinar/rambler/releases) of the project.
The standard release include executables for linux, mac and windows on both 386 and amd64 architectures.
The compilation is done using the wonderful [GoXC](https://github.com/laher/goxc) tool, let me know if you want your particular architecture added to the standard set.

Go users can also simply compile it from source and install it as a go executable using the following command :

```
go install github.com/elwinar/rambler
```

## Usage

### Migrations

In rambler, migrations are kept in the simplest form possible: a migration is a list of sections (`up` and `down`), each section being an SQL statement. Example:

```sql
-- rambler up

CREATE TABLE foo (
	id INTEGER UNSIGNED AUTO_INCREMENT,
	bar VARCHAR(60),
	PRIMARY KEY (id)
);

-- rambler down

DROP TABLE foo;
```

Sections are delimited by SQL comments sufixed by the rambler marker (white-spaces sensitive). While applying a migration, rambler will execute each `up` section in order, and while reversing it it will execute each `down` section in reverse order.

Migrations are executed in alphabetical order, thus a versionning scheme of the form `version_description.sql` is highly recommended, version being an integer value, and description an underscored string. Examples:

* `201409272258_Added_table_foo.sql`
* `01_First_migration.sql`

### Configuration

Rambler configuration is lightweight: just dump the credentials of your database and the path to your migrations' directory into a JSON file, and you're done. Here is an example or JSON configuration file with the default values of rambler:

```json
{
	"driver": "mysql",
	"protocol": "tcp",
	"host": "localhost",
	"port": 3306,
	"user": "root",
	"password": "",
	"database": "",
	"directory": "."
}
```

When running, rambler will try to find a configuration file in the working directory and use its values to connect to the managed database.

#### Drivers

Rambler supports actually 3 drivers:

- `mysql`
- `postgresql`
- `sqlite`

Don't hesitate to get in touch if you want to see another one supported, provided a golang `database/sql` driver exist for your database vendor.

### Applying a migration

To apply a migration, use the `apply` command.

```
rambler apply
```

Rambler will compare the migrations already applied and the available migrations in increasing order to find the next migration to apply, then execute all its `up` sections' statements in order. 

### Reversing a migration

To reverse a migration, use the `reverse` command.

```
rambler reverse
```

Rambler will compare the migrations already applied and the available migrations in decreasing order to find the last applied migrations, then execute all its `down` sections' statements in reverse order.

### Options

You can tell rambler to repeat the process while there is a migration to apply (or reverse) with the `all` flag (or its shorthand, `a`).

### Errors

To ensure database schema consistency, rambler will complain and stop when encountering a new migration in the middle of the already existing ones or if it can't find a migration already applied.

### Environments

An environment is an additionnal configuration that is given a name, and can be used to create multiple configurations for a single application (for example, to differenciate production, testing, etc).

Environments are defined in the configuration file, under the `environments` item.
Each environment is defined as an attribute of this item, the key being the name and the value being the configuration options.

Environments configuration are derived from the default configuration of rambler (at the configuration file's root), so you only need to override the needed options:

```json
{
	"driver": "mysql",
	"protocol": "tcp",
	"port": 3306,
	"user": "root",
	"password": "",
	"database": "rambler_default",
	"directory": "migrations",
	"environments": {
		"development": {
			"database": "rambler_development"
		},
		"testing": {
			"database": "rambler_testing"
		}
	}
}
```

Here we have three environments defined:
- `default`, will use the `rambler_default` database,
- `development`, will use the `rambler_development` database,
- `testing`, will use the `rambler_testing` database;

## CONTRIBUTORS

- [cjhubert](https://github.com/cjhubert)

## CHANGELOG

- **3.2**
	- Add GoXC configuration file to automate
- **3.1**
	- Add continuous integration tooling
	- Add the [Unlicense](http://unlicense.org/) LICENSE file
	- Add ego-growing badges on the top of the README file
- **3.0**
	- Major refactor to improve the unit testing of the tool
	- Fix the bugs introduced by the mergo breaking change
- **2.1**
	- Added the postgresql driver from the good work of [cjhubert](https://github.com/cjhubert)
- **2.0**
	- Complete rewrite to add unit-testing
	- Removed the command-line options to override the configuration
	- Fixed behavior of the transactions
	- Switched from spf13/cobra & spf13/viper to codegangsta/cli and encoding/json : removed the posibility of using alternatives configuration file types (namely YAML and TOML), changed the command-line usage
- **1.1.0**
	- Added environments handling
- **1.0.2**
	- Fixed a bug about the migration paths building whilte scanning
	- Made a real documentation
	- Fixed misconceptions about spf13/viper & spf13/cobra (thanks @spf13 for the pointers)
- **1.0.1**
	- Fixed imports paths of the internal packages

## Feedback and contributions

Feel free to give feedback, make pull requests or simply open issues if you find a bug or have an idea.
