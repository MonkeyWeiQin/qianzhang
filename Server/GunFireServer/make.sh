#!/bin/bash
server_name=$(basename "$PWD");
os=`uname`
os_name=$(echo $os|tr 'A-Z' 'a-z')
tail=""

helpmsg=" -h   --help         show help info. \n-s   --server       Packaged file name \n-o   --os           Package running system [windows|linux|darwin]"

until [ $# -eq 0 ]
do

    case "$1" in
        -h|--help) 
			echo -e $helpmsg
        	exit 0
		;;

        -s|--server)
			shift
			server_name=$1
		;;

        -o|--os)
			shift
			os_name=$1

			case $os_name in "linux" | "darwin" )
					export GOOS=$os_name
					;;
				"" )
					;;					
				*)
					export GOOS="windows"
					tail=".exe"
					;;
			esac
		;;
        *) 
			echo " unknow prop $1"
			shift
		;;
    esac

    shift
done


out_file="${server_name}_package_${GOOS}${tail}"


go build -v -x  -o $out_file .

