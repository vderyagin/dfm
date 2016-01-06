# DFM - dotfile manager #

## Installation ##

- Have [Go](http://golang.org/doc/install) installed
- Have your Go workspace configured (as per http://golang.org/doc/code.html#GOPATH)
- `go get gitlab.com/vderyagin/dfm`

## Usage ##

```
dfm --help
```

Basic idea is to have a separate directory where all your dotfiles are stored (under version control, most likely) and symlinked from their corresponding locations in home directory. File system tree structure in dotfile storage directory is the same as in home directory, except without leading dots in file paths.

Like this:

```
|-------------------------------+------------------------------|
| in home directory             | in dotfile storage directory |
|-------------------------------+------------------------------|
| .gitconfig                    | gitconfig                    |
| .gnupg/gpg.conf               | gnupg/gpg.conf               |
| .config/fontconfig/fonts.conf | config/fontconfig/fonts.conf |
|-------------------------------+------------------------------|
```

You get the idea.

Dotfile storage directory includes only files that were explicitly put in there. Only storing regular files is allowed, i.e. you can not store a directory. If you do want to store directory, you'll have to add each of files contained there (which is easy to do, run `dfm store .dir/**/*` or something like that).

### Operations ###

- `dfm list` lists all stored dotfiles, including their statuses (linked, conflict, etc).

- `dfm store` moves files, given as arguments, into their appropriate places in storage directory and links them back to their original paths in home directory.

- `dfm restore` is the opposite of `dfm store`, it replaces symlinks with original files, which are removed from storage directory.

- `dfm link` link all stored files to their original locations in home directory. On fresh machine you can just clone a repo with your dotfiles to `~/.dotfiles` and run `dfm link --force` to set everyting up.

- `dfm delete` removes file from storage directory and link to it from home directory. Also cleans up any empty directories left after files are removed.

`store` and `link` support `--force` flag, which allows them to overwrite conflicting files when necessary.

### Host-specific dotfiles ###

Sometimes you need some configuration file to have different options on different machines, and yet it would be convenient to have all dotfiles for all machines in one repository. DFM allows to achieve that sort of thing by using host-specific dotfiles.

Host-specific dotfiles are only used on machine they are intended for and ignored everywhere else. Machines are distinguished by hostnames (run `hostname` to find out).

To store given dotfiles as host-specific just use `--host-specific` flag when invoking `store` command, like this:

``` dfm store --host-specific .xinitrc ```

It will be stored with suffix ".host-[host name]" in your dotfile storage directory. If you also happen to have generic version of that dotfile (without host-specific suffix), it will be used on other machines and ignored on machine for which host-specific file is intended. Other commands (`list`, `restore`, `link`, `delete`) are smart enough to deal with host-specific files automatically and in a way that makes sense.

### Forcing regular files instead of symlinks ###

Some application require their dotfiles to be regular files, not symlinks to regular files stored elsewhere. DFM supports this, just use `--copy` flag when invoking `store` command, like this:

``` dfm store --copy .xinitc ```

It will be stored with suffix ".force-copy" in your dotfile storage directory. If that copy happens to diverge from stored version, this file will be considered in conflict by DFM. You'll be able to do run `dfm link --force <file>` to overwrite original file or `dfm store --force <file>` to overwrite stored version of it. Other commands also work with such files in a way that makes sense.

And yes, this files can also be host-specific.


## Options ##

Dotfile storage directory defaults to `~/.dotfiles` and home directory is, well, home directory of current user. It is possible to override both with `--store` and `--home` global options or with `DOTFILES_STORE_DIR` and `DOTFILES_HOME_DIR` enviroment variables. You probably not gonna need to override home directory, but it is possible to imagine situation in which it would be useful, like using `dfm` on remote filesystem through NFS.

There is no flag for overriding current hostname, but you can do it by setting `HOST` enviroment variable.
