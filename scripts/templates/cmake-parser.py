"""@package templates
Parses and completes the cmake templates
"""

import os
import scripts.module.helper as helper

lib_search_tag = "% lib-search\n"
firmware_gen_tag = "% firmware-gen\n"
end_tag = "% end\n"
fill_block_start = "{{"
fill_block_end = "}}"

src_file_exts = (".cpp", ".c", ".cc")
hdr_file_exts = (".hh", ".h")


def lib_search(content, project_data):
    """searches for library paths and then completes the templates to include search paths and build library"""

    str_to_return = ""
    for lib in helper.get_dirs(project_data["curr-path"] + "/lib"):
        src_files = []
        hdr_files = []

        # check if there is a src folder
        src_found = False

        # handle src folder
        for sub_dir in helper.get_dirs(lib):
            if os.path.basename(sub_dir) == 'src':
                # add all the src extensions in src folder
                src_files += helper.get_files_recursively(sub_dir, src_file_exts)

                # add all the header extensions in src folder
                hdr_files += helper.get_files_recursively(sub_dir, hdr_file_exts)
                src_found = True
                break

        if src_found is not True:
            # add all the src extensions
            src_files += helper.get_files_recursively(lib, src_file_exts)

            # add all the header extensions
            hdr_files += helper.get_files_recursively(lib, hdr_file_exts)

        # go through all files and generate cmake tags
        data = {'lib-path': ["\" \"".join(src_files + hdr_files)], 'name': os.path.basename(lib),
                'srcs': ["\" \"".join(src_files)],
                'hdrs': [" ".join(hdr_files)], 'board': project_data['board']}

        for line in content:
            line = line[2:len(line) - 3]
            str_to_return += helper.fill_template(line, data) + "\n"

    return str_to_return.strip(" ").strip("\n") + "\n"


def firmware_gen(content, project_data):
    """searches for src files and then generates the firmware code for linking and building the project"""

    curr_src_path = project_data["curr-path"] + "/src"
    curr_lib_path = project_data["curr-path"] + "/lib"
    str_to_return = ""

    src_files = "\" \"".join(helper.get_files_recursively(curr_src_path, src_file_exts + hdr_file_exts))
    lib_files = " ".join(helper.get_dirnames(curr_lib_path))

    data = {'name': project_data["cmake-project-name"], 'srcs': [src_files], 'libs': lib_files,
            'cosa-libs': project_data["cosa-libs"], 'port': project_data['port'], 'board': project_data['board']}

    for line in content:
        line = line[2:len(line) - 3]
        str_to_return += helper.fill_template(line, data) + "\n"

    return str_to_return.strip(" ").strip("\n") + "\n"


def get_elements(tpl_str, curr_index):
    """gather elements from the template block"""

    content = []

    # gather all the lines inside the loop block
    content_index = curr_index + 1
    while True:
        line = tpl_str[content_index]

        if line == end_tag:
            break
        else:
            content.append(line)

        content_index += 1

    return content, content_index


def parse(tpl_path, project_data):
    """reads the cmake template file and completes it using project data"""

    tpl_path = os.path.abspath(tpl_path)
    tpl_file = open(tpl_path)
    tpl_str = tpl_file.readlines()
    tpl_file.close()

    new_str = ""
    index = 0
    while index < len(tpl_str):
        curr_line = tpl_str[index]

        # handle loop statements
        if curr_line == lib_search_tag:
            result = get_elements(tpl_str, index)

            new_str += lib_search(result[0], project_data)
            index = result[1]
        elif curr_line == firmware_gen_tag:
            result = get_elements(tpl_str, index)

            new_str += firmware_gen(result[0], project_data)
            index = result[1]
        else:
            new_str += helper.fill_template(curr_line, project_data)
        index += 1
    print(new_str)


if __name__ == '__main__':
    curr_path = os.path.abspath(os.path.dirname(__file__) + "/../../core/Test")
    json_data = {'board': "uno", "curr-path": curr_path, "wcosa-path": curr_path, "cmake-version": "3.1.0",
                 "cmake-project-name": "deep", "port": "com4", "cosa-libs": ""}
    parse("../../core/templates/cmake/CMakeLists.txt.tpl", json_data)
