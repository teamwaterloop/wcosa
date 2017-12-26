import argparse
import handle
from colorama import Fore
from parsers import board_parser
from others import output, helper


def parse():
    """Adds command line arguments and parses them"""

    parser = argparse.ArgumentParser(description="WCosa create, build and upload Cosa AVR projects")

    parser.add_argument('action', help='action to perform (create, update, build, upload, serial and boards')
    parser.add_argument('--board', help='board to use for wcosa project',
                        type=str)
    parser.add_argument('--port', help='port to upload the AVR traget to (default: automatic)',
                        type=str)
    parser.add_argument('--programmer', help='port to upload the AVR traget to (default: usbtinyisp)',
                        type=str)
    parser.add_argument('--baud', help='buad rate for serial (default: 9600)',
                        type=int)
    parser.add_argument('--ide', help='create specific project structure for specific ide (default: none)',
                        type=str)
    parser.add_argument('--path', help='path to create the project at (default: curr dir)',
                        type=str)

    return parser.parse_args()


def verify_board(board):
    """Verify if the board provided is supported by wcosa"""

    boards = board_parser.get_all_board(helper.get_wcosa_path() + "/wcosa/boards.json")

    if board is not None and board not in boards:
        output.writeln("Board Invalid. Run wcosa script with boards option to see all the valid boards", Fore.RED)
        quit(2)


def print_boards():
    """Print all the available boards and their name"""

    boards = board_parser.get_all_board(helper.get_wcosa_path() + "/wcosa/boards.json")

    output.writeln("Boards compatible with this project are: ", Fore.CYAN)

    for board in boards:
        name = board_parser.get_board_properties(board, helper.get_wcosa_path() + "/wcosa/boards.json")["name"]
        output.writeln('{:15s} --->\t{}'.format(board, name))


if __name__ == "__main__":
    options = parse()

    # based on the action call scripts
    if options.action == "boards":
        print_boards()
    elif options.action == "create":
        if options.board is not None:
            # verify the board and create
            verify_board(options.board)
            handle.create_wcosa(options.path, options.board, options.ide)
        else:
            output.writeln("Board is needed for creating wcosa project", Fore.RED)
    elif options.action == "update":
        # verify the board and update
        verify_board(options.board)
        handle.update_wcosa(options.path, options.board)
