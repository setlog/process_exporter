package flags

import "flag"

func Parse(arguments []string) (namespace, binary, argName string) {
	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	namespacePtr := flagSet.String("namespace", "mine", "Prometheus metric namespace.")
	binaryPtr := flagSet.String("binary", "", "Filter which processes to watch by binary name.")
	argNamePtr := flagSet.String("argname", "name", "Set Prometheus \"name\"-label value to value of this command line argument of the monitored processes.")
	err := flagSet.Parse(arguments)
	if err != nil {
		panic(err)
	}
	return *namespacePtr, *binaryPtr, *argNamePtr
}
