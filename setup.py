#!/usr/bin/env python

import os

from setuptools import find_packages, setup

install_requires = [
    'colorama',
    'pyserial>=3,<4',
]


def package_files(directory):
    paths = []
    for (path, directories, filenames) in os.walk(directory):
        for filename in filenames:
            paths.append(os.path.join('..', path, filename))
    return paths


setup(
    name='WCosa',
    version='0.1.0',
    description='Create, Build, Upload and Monitor AVR Cosa Projects',
    author='Deep Dhillon, Jeff Niu, Ambareesh Balaji',
    author_email='deep.dhill6@gmail.com, jeffniu22@gmail.com, ambareeshbalaji@gmail.com',
    long_description=open('README.md').read(),
    license='MIT',
    packages=find_packages(),
    install_requires=install_requires,
    package_data={
        '': package_files('toolchain') + package_files('templates'),
        'wcosa': ['*.json'],
    },
    entry_points={
        'console_scripts': [
            'wcosa = wcosa.wcosa:main',
        ],
    },
    classifiers=[
        'Development Status :: 4 - Beta',
        'License :: MIT license',
    ],
    keywords=[
        'iot', 'embedded', 'arduino', 'avr', 'fpga', 'firmware',
        'hardware', 'microcontroller', 'debug', 'cosa', 'tool',
    ],
)
