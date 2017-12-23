"""
Handle handles creating and updating WCosa projects
"""

from shutil import copyfile
from colorama import Fore
from core.scripts.others.output import write, writeln
from core.scripts.others import helper
import core.scripts.templates.config as config
import core.scripts.templates.cmake as cmake


def create_wcosa(path, board, ide):
    """Creates WCosa project from scratch"""

    project_path = path

    if path is None:
        project_path = helper.get_working_directory()

    if ide is None:
        ide = ""
    else:
        ide = ide.strip(" ")

    templates_path = helper.get_wcosa_path() + "/core/templates"
    user_config_path = project_path + "/config.json"
    internal_config_path = project_path + "/wcosa/internal-config.json"
    general_cmake_path = project_path + "/wcosa/CMakeLists.txt"

    write("Creating work environment - ", color=Fore.CYAN)

    # create src, lib, and wcosa folders
    helper.create_folder(project_path + "/src", True)
    helper.create_folder(project_path + "/lib", True)
    helper.create_folder(project_path + "/wcosa", True)
    helper.create_folder(project_path + "/wcosa/bin", True)

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + "/cmake/CMakeLists.txt.tpl", general_cmake_path)
    copyfile(templates_path + "/config/internal-config.json.tpl", internal_config_path)
    copyfile(templates_path + "/config/config.json.tpl", user_config_path)

    if ide == "clion":
        copyfile(templates_path + "/ide/clion/CMakeLists.txt.tpl", project_path + "/CMakeLists.txt")

    writeln("done")
    write("Updating configurations based on the system - ", color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, "None", ide)  # give a dummy port right now
    project_data = config.fill_internal_config(internal_config_path, path, user_data)

    import core.scripts.templates.cmake as cmake

    cmake.parse_update(general_cmake_path, project_data)

    if ide != "":
        cmake.parse_update(project_path + "/CMakeLists.txt", project_data)

    writeln("done")
    writeln("Project Created and structure:", color=Fore.YELLOW)
    writeln("src    ->    All source files go here:", color=Fore.YELLOW)
    writeln("lib    ->    All custom libraries go here", color=Fore.YELLOW)
    writeln("wcosa  ->    All the build files are here (do no modify)", color=Fore.YELLOW)


def update_wcosa(path, board):
    """Updates existing WCosa project"""

    write("Updating work environment - ", color=Fore.CYAN)

    project_path = path

    if path is None:
        project_path = helper.get_working_directory()

    templates_path = helper.get_wcosa_path() + "/core/templates"
    user_config_path = project_path + "/config.json"
    internal_config_path = project_path + "/wcosa/internal-config.json"
    general_cmake_path = project_path + "/wcosa/CMakeLists.txt"

    # create src, lib, and wcosa folders
    helper.create_folder(project_path + "/src")
    helper.create_folder(project_path + "/lib")
    helper.create_folder(project_path + "/wcosa")
    helper.create_folder(project_path + "/wcosa/bin")

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + "/cmake/CMakeLists.txt.tpl", general_cmake_path)
    copyfile(templates_path + "/config/internal-config.json.tpl", internal_config_path)

    writeln("done")
    write("Updating configurations with new changes - ", color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, "None")  # give a dummy port right now
    ide = user_data["ide"]

    if ide == "clion":
        copyfile(templates_path + "/ide/clion/CMakeLists.txt.tpl", project_path + "/CMakeLists.txt")

    project_data = config.fill_internal_config(internal_config_path, path, user_data)

    cmake.parse_update(general_cmake_path, project_data)

    if ide != "":
        cmake.parse_update(project_path + "/CMakeLists.txt", project_data)

    writeln("done")
