import argparse
import glob


def main():
    parser = argparse.ArgumentParser(prog="flow.py",
                                     description="Welcome to GoFlow! Use this Python script to generate your code",
                                     usage="python %(prog)s [options]")
    parser.add_argument("--data", type=str, metavar="",
                        help="The name of _Data struct which stores all the data you need")
    parser.add_argument("--result", type=str, metavar="",
                        help="The name of _Result struct which stores all the output")
    parser.add_argument("--prepare", type=str, metavar="",
                        help="The name of _PrepareInput struct which stores the argument prepare function needs")
    parser.add_argument("-s", "--source", type=str, metavar="", nargs="?", default=".",
                        help="The directory of the template files of GoFlow")
    parser.add_argument("-o", "--output", type=str, metavar="", nargs="?", default=".",
                        help="The directory of the generated files to put in")
    parser.add_argument("-p","--package",type=str, metavar="",
                        help="The package name of the output Golang source file")

    args = parser.parse_args()

    # Check the arguments
    if args.data is None or args.result is None or args.prepare is None:
        print("please provided '--data', '--result' and '--prepare'")
        return

    # Find all the files
    for name in ['go_flow', 'structure']:
        file = glob.glob(args.source + f"/{name}.go")
        if file:
            with open(f'{args.output}/{name}.go', 'w') as output:
                with open(file[0], 'r') as source:
                    lines = []
                    for line in source.readlines():
                        line = line.replace("_Data", args.data).replace("_Result", args.result).replace("_PrepareInput",
                                                                                                     args.prepare)
                        if args.package is not None:
                            line = line.replace("package goflow",f"package {args.package}")
                        lines.append(line)
                    output.writelines(lines)

    # Print the message
    print("[SUCCESS]")


if __name__ == "__main__":
    main()
