package main

type configuration struct {
	Base        string
	MasterURI   string
	masterAddr  string
	masterSSL   bool
	SlaveURI    string
	slaveAddr   string
	slaveSSL    bool
	InsecureSSL bool
	CACert      string
}
