import sys
sys.path.append('gen-py')

import time
import argparse

from proxy import ProxyService

from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TCompactProtocol

class Args:
    def __init__(self):
        self.__parser = argparse.ArgumentParser(description='OSRM traffic updater against traffic proxy.', formatter_class=argparse.ArgumentDefaultsHelpFormatter)

        # params
        self.__parser.add_argument('-p', '--port', default=6666, help='traffic proxy listening port', type=int)
        self.__parser.add_argument('-c', '--ip', default='127.0.0.1', help='traffic proxy ip address', type=str)
        self.__parser.add_argument('-f', '--csv_file', default='traffic.csv', help='OSRM traffic csv file', type=str)

        # parse
        self.__args = self.__parser.parse_args()

    def get(self):
        return self.__args


def dump_flows_to_csv(csv_file, flows):
    with open(csv_file, "wb") as writer:
        i = 0
        lines_buff = []
        lines_count_per_write = 1000
        total_wrote_count = 0

        for flow in flows:
            osrm_csv_str_line = str(flow.fromId) + "," + str(flow.toId) + "," + str(flow.speed) + "\n"
            lines_buff.append(osrm_csv_str_line)
            if i < 10:  # print first 10 lines
                print "[ " + str(i) + " ] " + str(flow)
                print "[ " + str(i) + " ] " + osrm_csv_str_line        

            # append to csv file
            if len(lines_buff) >= lines_count_per_write:
                writer.writelines(lines_buff)
                total_wrote_count += len(lines_buff)
                lines_buff = []

            i += 1

        if len(lines_buff) > 0:
            writer.writelines(lines_buff)
            total_wrote_count += len(lines_buff)
            lines_buff = []
    print "total wrote to " + csv_file + " count: " + str(total_wrote_count)
    

def main():
    args = Args().get()

    # Make socket
    print "connect traffic_proxy: " + args.ip + ":" + str(args.port)
    transport = TSocket.TSocket(args.ip, args.port)

    # Buffering is critical. Raw sockets are very slow
    transport = TTransport.TFramedTransport(transport)

    # Wrap in a protocol
    protocol = TCompactProtocol.TCompactProtocol(transport)

    # Create a client to use the protocol encoder
    client = ProxyService.Client(protocol)

    # Connect!
    transport.open()

    # Get flows 
    start_time = time.clock()
    #flow = client.getFlowById(991239906)
    #print flow
    print "getting flows"
    flows = client.getAllFlows()
    print "got flows count: " + str(len(flows))
    after_get_flow_time = time.clock() 
    print "get flows time used: " + str(after_get_flow_time - start_time) + " seconds"

    # Dump to OSRM csv format
    print "dump flows to: " + args.csv_file
    dump_flows_to_csv(args.csv_file, flows)    

    end_time = time.clock()
    print "dump csv time used: " + str(end_time - after_get_flow_time) + " seconds"


    # Close!
    transport.close()


if __name__ == '__main__':
    main()
    
