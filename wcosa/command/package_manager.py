import git
import re
from wcosa.utils.output import write, writeln

class Package:
    def __init__(self, name, url, version):
        self.name = name
        self.url = url
        self.version = version

    def __str__(self):
        return "package %s version %s at %s" % (self.name, self.version, self.url)

    def __repr__(self):
        return "name: %s, version: %s, url: %s" % (self.name, self.version, self.url)

class PackageFormatError(Exception):
    def __init__(self, package_string):
        self.package_string = package_string

    def __str__(self):
        return "Bad package format: " + self.package_string

FULL_URL = r'(?P<url>https?://\S+/(?P<name>\S+))'
GITHUB = r'(?P<github>\w+/(?P<name>\w+))'
VERSION = r':(?P<version>\S+)'
EXPLICIT_NAME = r'( as (?P<explicit_name>\S+))?'
VALID_SCHEMAS = [re.compile('^' + FULL_URL + VERSION + EXPLICIT_NAME + '$'),
                 re.compile('^' + GITHUB + VERSION + EXPLICIT_NAME + '$')]

def parse_entries(package_strings):
    """
    Convert package strings to package entries.
    A package string is of the form '(BASE_URL|GITHUB_SHORTHAND):VERSION [as NAME]'
    where:
        BASE_URL is a valid URL pointing to a git repository
        GITHUB_SHORTHAND is of the form 'username/reponame'
        VERSION is some tag on the given repository
        NAME is the preferred package name
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
        if groups['explicit_name']: # always a group, possibly empty
            name = groups['explicit_name']
        else:
            name = groups['name']
        version = groups['version']
        packages.append(Package(name, url, version))
    return packages
