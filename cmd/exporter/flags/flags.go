package flags

import "flag"

func Parse(arguments []string) (namespace, binary, argName string, port int) {
	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	namespacePtr := flagSet.String("namespace", "mine", "Prometheus metric namespace.")
	binaryPtr := flagSet.String("binary", "", "Filter which processes to watch by binary name.")
	nameFlagPtr := flagSet.String("nameflag", "name", "Set Prometheus \"name\"-label value to value of this command line argument of the monitored processes.")
	portPtr := flagSet.Int("port", 80, "Port on which to listen to requests to /metrics.")
	err := flagSet.Parse(arguments)
	if err != nil {
		panic(err)
	}
	return *namespacePtr, *binaryPtr, *nameFlagPtr, *portPtr
}
