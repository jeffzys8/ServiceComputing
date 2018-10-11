package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"math"
	"os/exec"
	flag "github.com/spf13/pflag"
)

type sp_args struct{
	start_page		int
	end_page		int
	in_filename		string
	page_len		int	/* default value, can be overriden by "-l number" on command line */
	page_type		bool /* false for lines-delimited, true for form-feed-delimited */
	print_dest		string
}

var progname string /* program name, for error msgs */
const INT_MAX = math.MaxInt64 /* MaxInt64 */
var err error /* for error output */

func main(){

	var sa sp_args

	/* save name by which program is invoked, for error messages */
	progname = os.Args[0]

	/* 原c程序中这里没有备注：此处对默认值进行设置 */
	sa.start_page = -1
	sa.end_page = -1
	sa.page_len = 72
	sa.page_type = false
	sa.print_dest = ""

	//using package 'pflag', no need to pass args from 'main' to following functions
	process_args(&sa)
	process_input(&sa)
	
}

func process_args(sa *sp_args) {
	

	/* check the command-line arguments for validity */
	if(len(os.Args) < 5){
		fmt.Fprintf(os.Stderr, "%s: %s\n", progname, "not enough arguments")
		usage()
		os.Exit(1)
	}

	/* handle mandatory args first */
	flag.IntVarP(&sa.start_page, "start-page", "s", 1, "start page number")
	flag.IntVarP(&sa.end_page, "end-page", "e", 1, "end page number")
	flag.IntVarP(&sa.page_len, "page-length", "l", 72, "lines per page")
	flag.BoolVarP(&sa.page_type, "page-type", "f", false, "page type")
	flag.StringVarP(&sa.print_dest, "print-dest", "d", "", "print dest")
	flag.Parse()

	/* handle 1st arg - start page */
	if sa.start_page < 0 || sa.start_page > (INT_MAX-1) {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %v\n", progname, sa.start_page)
		usage()
		os.Exit(2)
	}

	/* handle 2nd arg - end page */
	if sa.end_page < 0 || sa.end_page > (INT_MAX-1) || sa.end_page < sa.start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %v\n", progname, sa.end_page)
		usage()
		os.Exit(3)
	}

	/* now handle optional args */
	if sa.page_len < 1 || sa.page_len > (INT_MAX-1){
		fmt.Fprintf(os.Stderr, "%s: invalid page length %v\n", progname, sa.page_len)
		usage()
		os.Exit(6)
	}

	if sa.page_type == false {
		if sa.page_len < 1 {
			sa.page_len = 72
		}
	}

	/* check for the input file (if there is one) */
	if flag.NArg() == 1 {
		sa.in_filename = flag.Arg(0)
	} else {
		sa.in_filename = ""
	}

}


/*================================= process_input() ===============*/
func process_input(sa *sp_args) {

	var page_ctr int
	/* read the input file */
	reader := make_reader(sa)

	/* define the output dest */
	writer, sub_proc := make_writer(sa)

	if !sa.page_type {
		page_ctr = paging_by_line(reader, &writer, sa)
	} else {
		page_ctr = paging_by_page(reader, &writer, sa)
	}

	/* check for the page validity */
	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d),"+" no output written\n", progname, sa.start_page, page_ctr)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d),"+" less output than expected\n",progname, sa.end_page, page_ctr)
	}
	
	/* close the sub-process */
	if sub_proc != nil{
		writer.(io.WriteCloser).Close()
		sub_proc.Wait()
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}

func make_reader(sa *sp_args) *bufio.Reader{
	
	/* Stdin by default */
	in_fd := os.Stdin
	if len(sa.in_filename) > 0 {
		in_fd, err = os.Open(sa.in_filename)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	return bufio.NewReader(in_fd)
}

func make_writer(sa *sp_args) (io.Writer, *exec.Cmd) {
	var writer io.Writer
	var sub_proc *exec.Cmd

	/* Stdout by default */
	writer = os.Stdout
	if len(sa.print_dest) > 0 {

		/* define the input for the sub-process */
		sub_proc = exec.Command("cat", "-n")
		writer, err = sub_proc.StdinPipe()
		if err != nil && err != io.EOF {
			panic(err)
		}

		/* define the output for the sub-process */
		sub_proc.Stdout = os.Stdout
		sub_proc.Stderr = os.Stderr

		/* start the sub-process */
		sub_proc.Start()
	}
	return writer, sub_proc
}

func paging_by_line(reader *bufio.Reader, writer *io.Writer, sa *sp_args) int{
	var line string
	line_ctr := 0
	page_ctr := 1

	for {
		// get a line
		line, err = reader.ReadString('\n')

		// get to the end Ctrl^D
		if err == io.EOF {
			break
		}

		// error
		if err != nil && err != io.EOF {
			panic(err)
		}

		// page 'filpping' process
		line_ctr++
		if line_ctr > sa.page_len {
			page_ctr++
			line_ctr = 1
		}

		// output
		if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
			fmt.Fprintf(*writer, line)
		}
	}
	return page_ctr
}

func paging_by_page(reader *bufio.Reader, writer *io.Writer, sa *sp_args) int{
	var page string
	page_ctr := 1

	for {
		// fmt.Println("what")

		// get a page (till '\f')
		page, err = reader.ReadString('\f')

		// fmt.Println("the hell")

		// get to the end Ctrl^D
		if err == io.EOF {
			break
		}

		// error
		if err != nil && err != io.EOF {
			panic(err)
		}

		// page 'flipping' process
		page_ctr++
		if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
			fmt.Fprintf(*writer, page)
		}
	}
	return page_ctr
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s: %s\n", progname, "\n[USAGE] -s start_page -e end_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n")
}
