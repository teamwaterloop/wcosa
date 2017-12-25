"""@package templates
Parses and completes the config templates
"""

import json
import os

from wcosa.parsers import platform_parser, board_parser
from wcosa.others import helper


def fill_internal_config(path, curr_path, user_config_data):
    """fills the internal config file that will be used for internal build"""

    with open(helper.linux_path(path)) as f:
        internal_config_data = json.load(f)

    with open(helper.get_settings_path()) as f:
        settings_data = json.load(f)

    internal_config_data["project-name"] = os.path.basename(curr_path)
    internal_config_data["ide"] = user_config_data["ide"]
    internal_config_data["board"] = user_config_data["board"]
    internal_config_data["port"] = user_config_data["port"]
    internal_config_data["wcosa-path"] = helper.get_wcosa_path()
    internal_config_data["current-path"] = helper.linux_path(curr_path)
    internal_config_data["cmake-version"] = settings_data["cmake-version"]

    # get c and cxx flags
    board_properties = board_parser.get_board_properties(user_config_data["board"],
                                                         internal_config_data["wcosa-path"] + "/wcosa/boards.json")
    internal_config_data["cmake-c-flags"] = platform_parser.get_c_compiler_flags(board_properties,
                                                                                 internal_config_data[
                                                                                     "wcosa-path"] +
                                                                                 "/toolchain/cosa/platform.txt",
                                                                                 settings_data["include-extra-flags"])
    internal_config_data["cmake-cxx-flags"] = platform_parser.get_cxx_compiler_flags(board_properties,
                                                                                     internal_config_data[
                                                                                         "wcosa-path"] +
                                                                                     "/toolchain/cosa/platform.txt",
                                                                                     settings_data[
                                                                                         "include-extra-flags"])
    internal_config_data["cmake-cxx-standard"] = settings_data["cmake-cxx-standard"]
    internal_config_data["custom-definitions"] = user_config_data["build-flags"]
    internal_config_data["cosa-libraries"] = user_config_data["cosa-libraries"]

    with open(helper.linux_path(path), "w") as f:
        json.dump(internal_config_data, f, indent=settings_data["json-indent"])

    return internal_config_data


def fill_user_config(path, board, port, ide=""):
    """fills the user config file that will be used for internal build"""

    with open(helper.linux_path(path)) as f:
        user_config_data = json.load(f)

    with open(helper.get_settings_path()) as f:
        settings_data = json.load(f)

    user_config_data["board"] = board

    if ide != "":
        user_config_data["ide"] = ide

    user_config_data["framework"] = settings_data["framework"]
    user_config_data["port"] = port

    with open(helper.linux_path(path), "w") as f:
        json.dump(user_config_data, f, indent=settings_data["json-indent"])

    return user_config_data
