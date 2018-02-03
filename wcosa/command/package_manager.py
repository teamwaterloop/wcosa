"""
A lightweight package management system for WCosa.
*_list functions operate on lists of data (common scenario, better efficiency).
"""
import git
import json
import re
import os, sys
from wcosa.utils.output import write, writeln

class Package:
    def __init__(self, name, url, branch, version, path):
        self.name = name
        self.url = url
        self.branch = branch
        self.version = version
        self.path = path

    def __repr__(self):
        return ("name: %s, url: %s, branch: %s, version: %s, path: %s" %
                    (self.name, self.url, self.branch, self.version, self.path))

class PackageFormatError(Exception):
    def __init__(self, package_string):
        self.package_string = package_string

    def __str__(self):
        return "Bad package format: " + self.package_string

FULL_URL = r'(?P<url>https?://\S+/(?P<name>\S+))'
GITHUB = r'(?P<github>[\w\-]+/(?P<name>[\w\-]+))'
BRANCH = r'(:(?P<branch>[\w\-]+))?'
VERSION = r'(@(?P<version>\S+))?'
PATH = r'( as (?P<path>\S+))?'
VALID_SCHEMAS = [re.compile('^' + FULL_URL + BRANCH + VERSION + PATH + '$'),
                 re.compile('^' + GITHUB + BRANCH + VERSION + PATH + '$')]

def package_string_parse_list(package_strings):
    """
    Convert package strings to package entities.
    Package strings must match (BASE_URL|GITHUB)[:BRANCH][@VERSION][ as PATH]
    where:
        FULL_URL is a valid URL pointing to a git repository
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
        if 'github' in groups: # only a group if matched with github short form
            url = 'https://github.com/' + groups['github']
        else:
            url = groups['url']
        name = groups['name']
        branch = 'master' if not groups['branch'] else groups['branch']
        version = 'master' if not groups['version'] else groups['version']
        path = 'lib' + name if not groups['path'] else groups['path']
        packages.append(Package(name, url, branch, version, path))
    return packages

def package_info_write_list(pkgpath, packages):
    if not packages:
        return # Nothing to write
    repo = package_repo_open(pkgpath)
    newentries = []
    updentries = []
    with open(pkgpath + '/pkglist', 'r+') as pkglistfile:
        pkglist = json.loads(pkglistfile.read())
        pkgnames = list(map(lambda x: x['name'], pkglist))
        for package in packages:
            if package.name in pkgnames:
                index = pkgnames.index(package.name)
                if pkglist[index] == package.__dict__: # If same info present
                    continue
                pkglist[index] = package.__dict__
                updentries.append(package.name)
            else:
                pkglist.append(package.__dict__)
                newentries.append(package.name)
        pkglistfile.seek(0)
        pkglistfile.write(json.dumps(pkglist))
    repo.index.add(['pkglist'])
    repo.index.commit("Updated package list\n\n"
                      + ("New: %s\n" % ', '.join(newentries)
                          if newentries else "")
                      + ("Updated: %s\n" % ', '.join(updentries)
                          if updentries else ""))

def package_repo_open(pkgpath):
    try:
        return git.Repo(pkgpath)
    except: # If could not open, initialize
        return package_repo_init(pkgpath)

def package_repo_init(pkgpath):
    write("Initializing package repository... ")
    sys.stdout.flush()
    pkgrepo = git.Repo.init(pkgpath)

    with open(pkgpath + '/pkglist', 'w+') as pkglist:
        pkglist.write('[]') # Start with empty package list

    pkgrepo.index.add(['pkglist'])
    pkgrepo.index.commit("Initialized repository")
    writeln("Done")

    return pkgrepo

def package_install_list(path, packages):
    packages = package_name_parse_list(packages)
    pkgpath = path + '/.pkg'
    pkgrepo = package_repo_open(pkgpath)

    for package in packages:
        write("Installing %s... " % package.name)
        sys.stdout.flush()
        try:
            sm = pkgrepo.create_submodule(package.name, package.name,
                                          url=package.url, branch=package.branch)
            # FIXME: this is very slow. Do we need versioning functionality?
            # sm.repo.git.remote('add', 'origin', package.url)
            # sm.repo.git.fetch('--tags')
        except Exception as e:
            writeln("Package install failed: " + str(e))
        if package.version in sm.repo.tags:
            pass
            # FIXME: Checks out into pkgpath
            # sm.repo.git.checkout(package.version)
        else:
            writeln("No such version: %s. Checking out HEAD" % package.version)
        pkgrepo.index.commit("Installed " + package.name)
        try:
            installpath = os.path.abspath(pkgpath + '/' + package.name)
            linkpath = path + '/' + package.path
            linkbasedir = '/'.join(linkpath.split('/')[:-1])
            try:
                os.mkdir(linkbasedir)
            except:
                pass
            os.symlink(installpath, linkpath)
        except Exception as e:
            writeln("Could not link package to path: " + str(e))
        writeln("Installed.")
