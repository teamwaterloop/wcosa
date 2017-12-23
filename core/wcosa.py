import argparse


def parse():
    parser = argparse.ArgumentParser(description="WCosa create, build and upload Cosa AVR projects")

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

    opts = parser.parse_args()


