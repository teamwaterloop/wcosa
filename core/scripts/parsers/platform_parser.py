"""@package parsers
Parses the platform.txt file and gathers information about the current platform
"""

import os
import json

from core.scripts.others import helper


def fill_template(flag, board_properties, version):
    if "{build.mcu}" in flag:
        flag = flag.replace("{build.mcu}", board_properties["board-mcu"])
    elif "{build.f_cpu}" in flag:
        flag = flag.replace("{build.f_cpu}", board_properties["board-f_cpu"])
    elif "{runtime.ide.version}" in flag:
        flag = flag.replace("{runtime.ide.version}", version)

    return flag


def get_raw_flags(lines, identifier, include_extra):
    raw_flags = ""

    for line in lines:
        if "compiler." + identifier + ".flags=" in line:
            raw_flags += line[line.find("=") + 1:].strip(" ").strip("\n")
        elif include_extra and "compiler." + identifier + ".extra_flags=" in line:
            raw_flags += " " + line[line.find("=") + 1:].strip(" ").strip("\n")

    return raw_flags


def get_c_compiler_flags(board_properties, platform_path, include_extra=True):
    platform_file = open(helper.linux_path(platform_path))
    raw_flags = get_raw_flags(platform_file.readlines(), "c", include_extra)

    settings_file = open(helper.linux_path(os.path.dirname(__file__) + "/../settings.json"))
    settings_data = json.load(settings_file)
    settings_file.close()

    processed_flags = ""

    for flag in raw_flags.split(" "):
        processed_flags += fill_template(flag, board_properties, settings_data["arduino-version"]) + " "

    return processed_flags.strip(" ")


def get_cxx_compiler_flags(board_properties, platform_path, include_extra=True):
    platform_file = open(helper.linux_path(platform_path))
    raw_flags = get_raw_flags(platform_file.readlines(), "cpp", include_extra)

    settings_file = open(helper.linux_path(os.path.dirname(__file__) + "/../settings.json"))
    settings_data = json.load(settings_file)
    settings_file.close()

    processed_flags = ""

    for flag in raw_flags.split(" "):
        processed_flags += fill_template(flag, board_properties, settings_data["arduino-version"]) + " "

    return processed_flags.strip(" ")
