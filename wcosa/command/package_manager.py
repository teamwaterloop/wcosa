"""
A lightweight package management system for WCosa.
*_many functions operate on lists of data (common scenario, better efficiency).
"""

import json
import os
import re
import sys

import git

from wcosa.utils.output import write, writeln


class Package:
    def __init__(self, name, url, branch, version, path):
        self.unqualified_name = name
        self.url = url
        self.branch = branch
        self.version = version
        self.path = path
        self.name = (self.unqualified_name +
                     ('-' + self.branch if self.branch != 'master' else '') +
                     ('-' + self.version if self.version != 'master' else ''))

    def __repr__(self):
        return ('name: %s, url: %s, branch: %s, version: %s, path: %s' %
                (self.name, self.url, self.branch, self.version, self.path))


class PackageFormatError(Exception):
    def __init__(self, package_string):
        self.package_string = package_string

    def __str__(self):
        return 'Bad package format: ' + self.package_string


class GitFetchException(Exception):
    def __init__(self, package):
        self.url = (package.url +
                    (':' + package.branch if package.branch != 'master'
                     else '')
                    ('@' + package.version if package.version != 'master'
                     else ''))

    def __str__(self):
        return 'Could not fetch submodule from %s' + self.url


class AlreadyInstalledException(Exception):
    def __init__(self, link_updated):
        self.link_updated = link_updated


URL = r'(?P<url>https?://\S+/(?P<name>\S+))'
GITHUB = r'(?P<github>[\w\-]+/(?P<name>[\w\-]+))'
BRANCH = r'(:(?P<branch>[\w\-]+))?'
VERSION = r'(@(?P<version>\S+))?'
PATH = r'( as (?P<path>\S+))?'
VALID_SCHEMAS = [re.compile('^' + URL + BRANCH + VERSION + PATH + '$'),
                 re.compile('^' + GITHUB + BRANCH + VERSION + PATH + '$')]


def package_dir_path(path):
    """Return package path to package install directory"""
    return path + '/.pkg'


def package_string_parse_many(package_strings):
    """
    Convert package strings to package entities.
    Package strings must match (URL|GITHUB)[:BRANCH][@VERSION][ as PATH]
    where:
        URL is a valid URL pointing to a git repository
        GITHUB is of the form 'username/reponame'
        BRANCH [default master] is the branch to track
        VERSION [default master] is a tag on the given branch
        PATH [default 'lib/NAME'] is the relative path to install location
    """
    packages = []
    for package_string in package_strings:
        for schema in VALID_SCHEMAS:
            match = re.match(schema, package_string)
            if match:
                groups = match.groupdict()
                break
        if not match:
            raise PackageFormatError(package_string)
        if 'github' in groups:  # only a group if matched with github format
            url = 'https://github.com/' + groups['github']
        else:
            url = groups['url']
        name = groups['name']
        branch = 'master' if not groups['branch'] else groups['branch']
        version = 'master' if not groups['version'] else groups['version']
        path = 'lib/' + name if not groups['path'] else groups['path']
        packages.append(Package(name, url, branch, version, path))
    return packages


def package_list_read(pkgpath):
    """Read package list"""
    try:
        with open(pkgpath + '/pkglist', 'r') as pkglistfile:
            return json.loads(pkglistfile.read())
    except Exception:
        return []


def package_list_write_many(pkgpath, packages):
    """Update package list with the given list of packages"""
    if not packages:
        return  # Nothing to write
    repo = package_repo_open(pkgpath)
    newentries = []
    updentries = []
    with open(pkgpath + '/pkglist', 'r+') as pkglistfile:
        pkglist = json.loads(pkglistfile.read())
        pkgnames = list(map(lambda x: x['name'], pkglist))
        for package in packages:
            if package.name in pkgnames:
                index = pkgnames.index(package.name)
                if package.path in pkglist[index]['paths']:
                    continue
                pkglist[index]['paths'].append(package.path)
                updentries.append(package.name)
            else:
                pkglist.append(package.__dict__)
                pkglist[-1]['paths'] = [package.path]
                del pkglist[-1]['path']
                newentries.append(package.name)
        pkglistfile.seek(0)
        pkglistfile.write(json.dumps(pkglist))
    repo.index.add(['pkglist'])
    if repo.is_dirty():  # Something has changed
        repo.index.commit('Updated package list\n\n' +
                          ('New: %s\n' % ', '.join(newentries)
                           if newentries else '') +
                          ('Updated: %s\n' % ', '.join(updentries)
                           if updentries else ''))


def package_repo_open(pkgpath):
    """Try to open package repo; initalize upon failure"""
    try:
        return git.Repo(pkgpath)
    except Exception:
        return package_repo_init(pkgpath)


def package_repo_init(pkgpath):
    """Initialize package repo"""
    write('Initializing package repository... ')
    sys.stdout.flush()
    pkgrepo = git.Repo.init(pkgpath)

    with open(pkgpath + '/pkglist', 'w+') as pkglist:
        pkglist.write('[]')  # Start with empty package list

    pkgrepo.index.add(['pkglist'])
    pkgrepo.index.commit('Initialized repository')
    writeln('Done')

    return pkgrepo


def package_link(path, package):
    """Link package directory from pkgpath to package.path"""
    install_path = os.path.abspath(package_dir_path(path) + '/' + package.name)
    link_path = os.path.abspath(path + '/' + package.path)
    link_basedir = '/'.join(link_path.split('/')[:-1])
    try:
        os.mkdir(link_basedir)
    except Exception:
        pass  # Already exists or failed (then next try will fail)
    try:
        os.symlink(install_path, link_path)
    except Exception as e:
        try:  # Maybe the path is already linked
            current_path = os.readlink(link_path)
            if current_path == install_path:
                return  # Then we're done
        except Exception:
            pass
        raise (type(e))('Could not link package: ' + str(e))


def _package_install_unsafe(path, package, pkgrepo, pkglist, pkgnames):
    """
    NOT A PUBLIC INTERFACE: use package_install[_many] instead.

    Try to install a package and forward exceptions to the caller.
    Will leave package repository in dirty state.
    Returns
    """
    write('Installing %s... ' % package.name)
    sys.stdout.flush()
    if package.name in pkgnames:
        index = pkgnames.index(package.name)
        if package.path in pkglist[index]['paths']:
            writeln('Already installed.')
            raise AlreadyInstalledException(link_updated=False)
        else:
            write('Already installed, linking to %s... ' % package.path)
            sys.stdout.flush()
            package_link(path, package)
            writeln('Linked.')
            raise AlreadyInstalledException(link_updated=True)
    # If the above did not return, we need to actually install the package
    try:
        pkgrepo.create_submodule(package.name, package.name,
                                 url=package.url, branch=package.branch)
    except Exception:  # Default message is cryptic
        raise GitFetchException(package)
    package_link(path, package)
    writeln('Installed.')


def package_install(path, package, batch_mode=False, pkgrepo=None,
                    pkglist=None, pkgnames=None):
    """
    Install a package or roll back to last coherent state upon failure.
    If batch_mode is True, do not update package list (caller will update).
    Returns True on success, else (error or already installed) False.
    """
    pkgpath = package_dir_path(path)
    if pkgrepo is None:
        pkgrepo = package_repo_open(pkgpath)
    if pkglist is None:
        pkglist = package_list_read(pkgpath)
    if pkgnames is None:
        pkgnames = list(map(lambda x: x['name'], pkglist))
    try:
        _package_install_unsafe(path, package, pkgrepo, pkglist, pkgnames)
        pkgrepo.index.add(['.gitmodules', package.name])
        pkgrepo.index.commit('Installed ' + package.name)
        if not batch_mode:
            package_list_write_many(pkgpath, [package])
    except AlreadyInstalledException as e:
        return e.link_updated
    except Exception as e:  # Installation failed, roll back
        try:
            sm = pkgrepo.submodule(package.name)
            sm.remove()
        except Exception:
            pass
        pkgrepo.git.clean('-fdX')  # Remove all untracked files
        writeln('Install aborted.')
        writeln(str(e))
        return False
    return True


def package_install_many(path, packages):
    """Install a list of packages"""
    packages = package_string_parse_many(packages)
    installed_packages = []
    pkgpath = package_dir_path(path)
    pkglist = package_list_read(pkgpath)
    pkgnames = list(map(lambda x: x['name'], pkglist))
    pkgrepo = package_repo_open(pkgpath)

    for package in packages:
        if package_install(path, package, True, pkgrepo, pkglist, pkgnames):
            installed_packages.append(package)  # To be written to database
    if installed_packages:
        package_list_write_many(pkgpath, installed_packages)


def package_update_all(path):
    """Update all installed packages"""
    repo = package_repo_open(package_dir_path(path))
    for sm in repo.submodules:
        write('Updating %s... ' % sm.name)
        sm.update()
        writeln('Done.')
