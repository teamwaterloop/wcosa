"""@package parsers
Parses the boards.txt file and gathers information about the current board
"""

from core.scripts.others import helper


def get_board_properties(board, board_file_path):
    """parses the board file returns the properties of the board specified"""

    board_file = open(helper.linux_path(board_file_path))
    board_str = board_file.readlines()

    board_id_str = mcu_str = f_cpu_str = board_name_str = ""

    for line in board_str:
        if board in line:
            if "board=" in line:
                board_id_str = line[line.find('='):].strip("=").strip("\n").strip(" ")
            elif "mcu=" in line:
                mcu_str = line[line.find('='):].strip("=").strip("\n").strip(" ")
            elif "f_cpu=" in line:
                f_cpu_str = line[line.find('='):].strip("=").strip("\n").strip(" ")
            elif "name=" in line:
                board_name_str = line[line.find('='):].strip("=").strip("\n").strip(" ")

    return {"board-name": board_name_str, "board-id": board_id_str, "board-mcu": mcu_str, "board-f_cpu": f_cpu_str}


def get_all_board(board_file_path):
    """returns the common names of all the boards (these names will be used by the user)"""

    board_file = open(helper.linux_path(board_file_path))
    board_str = board_file.readlines()
    boards = []

    for line in board_str:
        if "board=" in line:
            boards.append(line[:line.find(".")])

    return boards
